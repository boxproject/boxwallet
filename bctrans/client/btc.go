package client

import (
	"encoding/hex"
	"strconv"

	"log"

	"github.com/boxproject/boxwallet/bccoin"
	"github.com/boxproject/boxwallet/bcconfig"
	"github.com/boxproject/boxwallet/bccore"
	"github.com/boxproject/boxwallet/bckey"
	"github.com/boxproject/boxwallet/bctrans/clientseries"
	"github.com/boxproject/boxwallet/errors"
	"github.com/boxproject/boxwallet/official"
	"github.com/boxproject/boxwallet/signature"
	"github.com/boxproject/lib-bitcore/sebtcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

var btcCliIntance *BtcClient

type BtcClient struct {
	*clientseries.OmniSeriesClient
	*bckey.KeyUtil
	nmap map[bool]*clientseries.OmniSeriesClient
	gap  int
}

func NewBtcClient(cfg bcconfig.Provider) *BtcClient {
	if btcCliIntance != nil {
		return btcCliIntance
	}
	btcCliIntance = &BtcClient{}
	bct := clientseries.NewOmniSeriesClient(cfg)
	btcCliIntance.OmniSeriesClient = bct
	btcCliIntance.KeyUtil = bckey.GetKeyUtilInstance()
	btcCliIntance.nmap = make(map[bool]*clientseries.OmniSeriesClient)
	btcCliIntance.nmap[true] = btcCliIntance.OmniSeriesClient
	btcCliIntance.nmap[false] = official.GetBtcNode()
	btcCliIntance.gap = cfg.GetInt("gap")
	return btcCliIntance
}
func GetBtcClientIntance() (*BtcClient, error) {
	if btcCliIntance == nil {
		return nil, errors.ERR_NIL_REFERENCE
	} else {
		return btcCliIntance, nil
	}
}
func (c *BtcClient) ImportAddress(address string, local bool) error {
	addr, err := btcutil.DecodeAddress(address, c.Env)
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
func (c *BtcClient) GetNewAddress() (bckey.GenericKey, error) {

	key, err := c.KeyUtil.GeneraterKey(nil, bccore.BC_BTC)
	if err != nil {
		return nil, err
	}
	address, err := btcutil.DecodeAddress(key.Address(), c.Env)
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
func (c *BtcClient) CreateTx(addrsFrom []string, token string, addrTo []*bccoin.AddressAmount, feeCeo float64) (txu signature.TxUtil, err error) {

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
	fee, _ := bccoin.NewCoinAmountFromFloat(bccore.BC_BTC, "", 0.0001*feeCeo)
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

	outsu, _ := bccoin.NewCoinAmount(bccore.BC_BTC, "", "0")     //unspent sum
	totalTran, _ := bccoin.NewCoinAmount(bccore.BC_BTC, "", "0") //tx cost sum
	for _, v := range addrTo {
		totalTran.Add(v.Amount) //总共花费
	}
	feesum, _ := bccoin.NewCoinAmount(bccore.BC_BTC, "", "0")
	tx := wire.NewMsgTx(wire.TxVersion) //构造tx

	totalTranDynamic, _ := bccoin.NewCoinAmount(bccore.BC_BTC, "", "0")
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
			//Dynamic calculation fee
			feesum.Set(fee)
			totalTranDynamic.Set(totalTran)

			i++
			size := float64(148*i + len(addrTo)*34 + 10 + 40)
			if size > 1000 {
				feesum, _ = bccoin.NewCoinAmountFromFloat(bccore.BC_BTC, "", 0.0001*feeCeo*size/1000)
			}
			totalTranDynamic.Add(feesum)
			cmpRe, err := outsu.Cmp(totalTranDynamic)
			if err != nil {
				return nil, err
			}
			if cmpRe == -1 {
				am, _ := bccoin.NewCoinAmount(bccore.BC_BTC, "", strconv.FormatFloat(vi.Amount, 'f', 8, 64))
				{
					//txin -------start-----------------
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
	// output 1, form----------------return-------------------
	addrf, err := btcutil.DecodeAddress(addrsFrom[0], c.Env)
	if err != nil {
		return
	}
	pkScriptf, err := txscript.PayToAddrScript(addrf)
	if err != nil {
		return
	}

	outsu.Sub(totalTranDynamic)
	tx.AddTxOut(wire.NewTxOut(outsu.Val().Int64(), pkScriptf))
	// output 2，to------------------pay-----------------
	for _, v := range addrTo {
		addrt, errInner := btcutil.DecodeAddress(v.Address, c.Env)
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
	return signature.NewBtcTx(tx, pkScripts, c.Env, local), nil
}
func (c *BtcClient) SendTx(txu signature.TxUtil) error {
	txv, err := txu.TxForSend()
	if err != nil {
		return nil
	}
	tx := txv.(*wire.MsgTx)
	_, err = c.nmap[txu.Local()].C.SendRawTransaction(tx, false)
	return err
}
func (c *BtcClient) GetBalance(address string, token string, local bool) (balance bccoin.CoinAmounter, err error) {
	unspents, err := c.getUnspentByAddress(address, local)
	if err != nil {
		return
	}
	balance, _ = bccoin.NewCoinAmount(bccore.BC_BTC, "", "0")
	for _, v := range unspents {

		f, errinner := bccoin.NewCoinAmount(bccore.BC_BTC, "", strconv.FormatFloat(v.Amount, 'f', 8, 64))
		if errinner != nil {
			err = errinner
			return
		}
		balance.Add(f)
	}
	return
}

///////////////////////////////////////////////////internal/////////////////////////////////////////////////////
func (c *BtcClient) getUnspentByAddress(address string, local bool) (unspents []sebtcjson.ListUnspentResult, err error) {
	btcAdd, err := btcutil.DecodeAddress(address, c.Env)
	if err != nil {
		return
	}
	log.Println(btcAdd.EncodeAddress())
	adds := [1]btcutil.Address{btcAdd}
	unspents, err = c.nmap[local].C.ListUnspentMinMaxAddresses(1, 999999, adds[:])
	return
}

func (c *BtcClient) addAddressToWallet(address btcutil.Address, local bool) error {
	err := c.checkAddressExists(address, local)
	if err != nil {
		return err
	}

	if err = c.C.ImportAddress(address.String()); err != nil {
		return err
	}
	return nil
}

func (c *BtcClient) checkAddressExists(address btcutil.Address, local bool) error {
	addrValid, err := c.nmap[local].C.ValidateAddress(address)
	if err != nil {
		return err
	}
	if addrValid.IsWatchOnly {
		return errors.ERR_DATA_EXISTS
	}
	return nil
}

func (c *BtcClient) ChooseClientNode() (local bool, err error) {
	bc, err := c.C.GetBlockChainInfo()
	if err != nil {
		return false, err
	}
	if bc.Headers-bc.Blocks > int32(c.gap) {
		return false, nil
	}
	return true, nil
}
