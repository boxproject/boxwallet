package client

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/boxproject/boxwallet/bccoin"
	"github.com/boxproject/boxwallet/bcconfig"
	"github.com/boxproject/boxwallet/bccore"
	"github.com/boxproject/boxwallet/bckey"
	"github.com/boxproject/boxwallet/bctrans/clientseries"
	"github.com/boxproject/boxwallet/errors"
	"github.com/boxproject/boxwallet/official"
	"github.com/boxproject/boxwallet/signature"
	"github.com/boxproject/lib-bitcore/sebtcjson"
	"github.com/ltcsuite/ltcd/chaincfg/chainhash"
	"github.com/ltcsuite/ltcd/txscript"
	"github.com/ltcsuite/ltcd/wire"
	"github.com/ltcsuite/ltcutil"
)

var ltcCliIntance *LtcClient

type LtcClient struct {
	*clientseries.LtcSeriesClient
	*bckey.KeyUtil
	nmap map[bool]*clientseries.LtcSeriesClient
	gap  int
}

func NewLtcClient(cfg bcconfig.Provider) *LtcClient {
	if ltcCliIntance != nil {
		return ltcCliIntance
	}
	ltcCliIntance = &LtcClient{}
	bct := clientseries.NewLtcSeriesClient(cfg)
	ltcCliIntance.LtcSeriesClient = bct
	ltcCliIntance.KeyUtil = bckey.GetKeyUtilInstance()
	ltcCliIntance.nmap = make(map[bool]*clientseries.LtcSeriesClient)
	ltcCliIntance.nmap[true] = ltcCliIntance.LtcSeriesClient
	ltcCliIntance.nmap[false] = official.GetLtcNode()
	ltcCliIntance.gap = cfg.GetInt("gap")
	return ltcCliIntance
}
func GetLtcClientIntance() (*LtcClient, error) {
	if ltcCliIntance == nil {
		return nil, errors.ERR_NIL_REFERENCE
	} else {
		return ltcCliIntance, nil
	}
}
func (c *LtcClient) ImportAddress(address string, local bool) error {
	addr, err := ltcutil.DecodeAddress(address, c.Env)
	if err != nil {
		return err
	}
	err = c.checkAddressExists(addr, local)
	if err != nil {
		if err == errors.ERR_DATA_EXISTS {
			return nil
		}
		return err
	}
	err = c.nmap[local].C.ImportAddressRescan(address, "", false)
	return err
}
func (c *LtcClient) GetNewAddress() (bckey.GenericKey, error) {

	key, err := c.KeyUtil.GeneraterKey(nil, bccore.BC_LTC)
	if err != nil {
		return nil, err
	}
	address, err := ltcutil.DecodeAddress(key.Address(), c.Env)
	if err != nil {
		return nil, err
	}
	err = c.checkAddressExists(address, true)
	if err != nil {
		if err == errors.ERR_DATA_EXISTS {
			return key, nil
		}
		return nil, err
	}
	err = c.C.ImportAddressRescan(key.Address(), "", false)
	return key, err
}
func (c *LtcClient) CreateTx(addrsFrom []string, token string, addrTo []*bccoin.AddressAmount, feeCeo float64) (txu signature.TxUtil, err error) {
	if len(addrsFrom) == 0 || len(addrTo) == 0 {
		return nil, errors.ERR_PARAM_NOT_VALID
	}
	local, err := c.ChooseClientNode()
	if err != nil {
		return nil, err
	}
	if !local {
		for _, v := range addrsFrom {
			err = c.ImportAddress(v, local)
			if err != nil {
				return nil, err
			}
		}
	}
	fee, _ := bccoin.NewCoinAmountFromFloat(bccore.BC_BTC, "", 0.001*feeCeo)
	unspentlist := make([][]sebtcjson.ListUnspentResult, 0, len(addrsFrom))
	for _, v := range addrsFrom {
		unspents, err := c.getUnspentByAddress(v, local)
		if err != nil {
			continue
		} else {
			unspentlist = append(unspentlist, unspents)
		}
	}
	pkScripts := make([]*signature.PKscripts, 0)
	// fee：size = inputsNum * 148 + outputsNum * 34 + 10
	outsu, _ := bccoin.NewCoinAmount(bccore.BC_LTC, "", "0")
	totalTran, _ := bccoin.NewCoinAmount(bccore.BC_LTC, "", "0")
	for _, v := range addrTo {
		totalTran.Add(v.Amount)
	}
	feesum, _ := bccoin.NewCoinAmount(bccore.BC_LTC, "", "0")
	tx := wire.NewMsgTx(wire.TxVersion)

	totalTranDynamic, _ := bccoin.NewCoinAmount(bccore.BC_LTC, "", "0")
	i := 0
	for k, v := range unspentlist {
		scripts := &signature.PKscripts{
			Address: addrsFrom[k],
			Scripts: make([][]byte, 0),
		}
		for _, vi := range v {
			if vi.Amount == 0 {
				continue
			}

			feesum.Set(fee)
			totalTranDynamic.Set(totalTran)

			i++
			size := float64(148*i + len(addrTo)*34 + 10 + 40)
			if size > 1000 {
				feesum, _ = bccoin.NewCoinAmountFromFloat(bccore.BC_LTC, "", 0.0001*feeCeo*size/1000)
			}
			totalTranDynamic.Add(feesum)
			cmpRe, err := outsu.Cmp(totalTranDynamic)
			if err != nil {
				return nil, err
			}
			if cmpRe == -1 {
				am, _ := bccoin.NewCoinAmount(bccore.BC_LTC, "", strconv.FormatFloat(vi.Amount, 'f', 8, 64))
				{
					//txin-------start-----------------
					hash, _ := chainhash.NewHashFromStr(vi.TxID)
					outPoint := wire.NewOutPoint(hash, vi.Vout)
					txIn := wire.NewTxIn(outPoint, nil, nil)

					outsu.Add(am)
					tx.AddTxIn(txIn)
					txinPkScript, errInner := hex.DecodeString(vi.ScriptPubKey)
					if errInner != nil {
						return nil, errInner
					}
					scripts.Scripts = append(scripts.Scripts, txinPkScript)
				}
			} else {
				break
			}
		}
		if len(scripts.Scripts) > 0 {
			pkScripts = append(pkScripts, scripts)
		}
	}
	cmpRe, err := outsu.Cmp(totalTranDynamic)
	if err != nil {
		return nil, err
	}
	if cmpRe == -1 {
		err = errors.ERR_NOT_ENOUGH_COIN
		return
	} else if cmpRe == 0 && i == 0 {
		err = errors.ERR_NOT_ENOUGH_COIN
		return
	}
	// form----------------return-------------------
	addrf, err := ltcutil.DecodeAddress(addrsFrom[0], c.Env)
	if err != nil {
		return
	}
	pkScriptf, err := txscript.PayToAddrScript(addrf)
	if err != nil {
		return
	}
	outsu.Sub(totalTranDynamic)
	tx.AddTxOut(wire.NewTxOut(outsu.Val().Int64(), pkScriptf))
	//to------------------pay-----------------
	for _, v := range addrTo {
		addrt, errInner := ltcutil.DecodeAddress(v.Address, c.Env)
		if errInner != nil {
			err = errInner
			return
		}
		pkScriptt, errInner := txscript.PayToAddrScript(addrt)
		if errInner != nil {
			err = errInner
			return
		}
		tx.AddTxOut(wire.NewTxOut(v.Amount.Val().Int64(), pkScriptt))

	}
	return signature.NewLtcTx(tx, pkScripts, c.Env, true), nil
}

func (c *LtcClient) SendTx(txu signature.TxUtil) error {
	txv, err := txu.TxForSend()
	if err != nil {
		return nil
	}
	tx := txv.(*wire.MsgTx)
	txHash, err := c.nmap[txu.Local()].C.SendRawTransaction(tx, false)
	fmt.Println(txHash)
	return err
}
func (c *LtcClient) GetBalance(address string, token string, local bool) (balance bccoin.CoinAmounter, err error) {
	unspents, err := c.getUnspentByAddress(address, local)
	if err != nil {
		return
	}
	balance, _ = bccoin.NewCoinAmount(bccore.BC_LTC, "", "0")
	for _, v := range unspents {

		f, errinner := bccoin.NewCoinAmount(bccore.BC_LTC, "", strconv.FormatFloat(v.Amount, 'f', 8, 64))
		if errinner != nil {
			err = errinner
			return
		}
		balance.Add(f)
	}
	return
}

///////////////////////////////////////////////////internal/////////////////////////////////////////////////////
func (c *LtcClient) getUnspentByAddress(address string, local bool) (unspents []sebtcjson.ListUnspentResult, err error) {
	btcAdd, err := ltcutil.DecodeAddress(address, c.Env)
	if err != nil {
		return
	}
	adds := [1]ltcutil.Address{btcAdd}
	unspents, err = c.nmap[local].C.ListUnspentMinMaxAddresses(1, 999999, adds[:])
	return
}

func (c *LtcClient) addAddressToWallet(address ltcutil.Address) error {
	err := c.checkAddressExists(address, true)
	if err != nil {
		return err
	}

	if err = c.C.ImportAddress(address.String()); err != nil {
		return err
	}
	return nil
}

func (c *LtcClient) checkAddressExists(address ltcutil.Address, local bool) error {
	addrValid, err := c.nmap[local].C.ValidateAddress(address)
	if err != nil {
		return err
	}
	if addrValid.IsWatchOnly {
		return errors.ERR_DATA_EXISTS
	}
	return nil
}

func (c *LtcClient) ChooseClientNode() (local bool, err error) {

	bc, err := c.C.GetBlockChainInfo()
	if err != nil {
		return false, err
	}
	if bc.Headers-bc.Blocks > int32(c.gap) {
		return false, nil
	}
	return true, nil
}
