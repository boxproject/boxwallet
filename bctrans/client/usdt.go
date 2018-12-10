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
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

var usdtCliInstance *UsdtClient

type UsdtClient struct {
	*clientseries.OmniSeriesClient
	*bckey.KeyUtil
	gap  int
	nmap map[bool]*clientseries.OmniSeriesClient
}

func NewUsdtClient(cfg bcconfig.Provider) *UsdtClient {
	if usdtCliInstance != nil {
		return usdtCliInstance
	}
	usdtCliInstance = &UsdtClient{}
	usdt := clientseries.NewOmniSeriesClient(cfg)
	usdtCliInstance.OmniSeriesClient = usdt
	usdtCliInstance.KeyUtil = bckey.GetKeyUtilInstance()
	usdtCliInstance.nmap = make(map[bool]*clientseries.OmniSeriesClient)
	usdtCliInstance.nmap[true] = usdtCliInstance.OmniSeriesClient
	usdtCliInstance.nmap[false] = official.GetUsdtNode()
	usdtCliInstance.gap = cfg.GetInt("gap")
	return usdtCliInstance
}
func GetUsdtClientInstance() (*UsdtClient, error) {
	if usdtCliInstance == nil {
		return nil, errors.ERR_NIL_REFERENCE
	} else {
		return usdtCliInstance, nil
	}
}
func (c *UsdtClient) ImportAddress(address string, local bool) error {
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
	err = c.C.ImportAddressRescan(address, "", false)
	return err
}
func (c *UsdtClient) GetNewAddress() (bckey.GenericKey, error) {
	key, err := c.KeyUtil.GeneraterKey(nil, bccore.BC_USDT)
	if err != nil {
		return nil, err
	}
	address, err := btcutil.DecodeAddress(key.Address(), c.Env)
	if err != nil {
		return nil, err
	}
	err = c.checkAddressExists(address, true)
	if err != nil {
		return nil, err
	}
	err = c.C.ImportAddressRescan(address.String(), "", false)
	return key, err
}

//1 to 1 only
func (c *UsdtClient) CreateTx(addrsFrom []string, token string, addrTo []*bccoin.AddressAmount, feeCeo float64) (txu signature.TxUtil, err error) {
	if len(addrsFrom) != 1 || len(addrTo) != 1 {
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
	addrF := addrsFrom[0]
	addrT := addrTo[0]
	//part-1------cmp usdt balance/////cmp usdt balance start///////

	omniBalance, err := c.GetBalance(addrF, "", local)
	if err != nil {
		return nil, err
	}
	cmp, err := omniBalance.Cmp(addrT.Amount)
	if err != nil {
		return nil, err
	}
	if cmp == -1 {
		return nil, errors.ERR_NOT_ENOUGH_COIN
	}
	///////////////////cmp usdt balance end/////////////

	//part-2 --------cmp btc balance----omniPrice and fee
	//omni transfer btc 0.00000546
	omniPrice, _ := bccoin.NewCoinAmount(bccore.BC_BTC, "", "0.00000546")
	unspents, err := c.getUnspentByAddress(addrF, local)
	if err != nil {
		return
	}
	// size = inputsNum * 148 + outputsNum * 34 + 10
	outsu, _ := bccoin.NewCoinAmount(bccore.BC_BTC, "", "0")
	feesum, _ := bccoin.NewCoinAmount(bccore.BC_BTC, "", "0")

	totalTran, _ := bccoin.NewCoinAmount(bccore.BC_BTC, "", "0")
	totalTran.Add(omniPrice)

	pkScripts := make([]*signature.PKscripts, 0)
	totalTranDynamic, _ := bccoin.NewCoinAmount(bccore.BC_BTC, "", "0")
	tx := wire.NewMsgTx(wire.TxVersion)
	scripts := &signature.PKscripts{
		Address: addrF,
		Scripts: make([][]byte, 0),
	}
	for k, v := range unspents {
		if v.Amount == 0 {
			continue
		}
		feesum.Set(fee)
		totalTranDynamic.Set(totalTran)
		////////////Buoyancy 40
		size := float64(148*(k+1) + 2*34 + 10 + 40)
		if size > 1000 {
			feesum, _ = bccoin.NewCoinAmountFromFloat(bccore.BC_BTC, "", 0.0001*feeCeo*size/1000)
		}
		totalTranDynamic.Add(feesum)
		cmpRe, err := outsu.Cmp(totalTranDynamic)
		if err != nil {
			return nil, err
		}
		if cmpRe == -1 {
			am, _ := bccoin.NewCoinAmount(bccore.BC_BTC, "", strconv.FormatFloat(v.Amount, 'f', 8, 64))
			{
				//txin-------start-----------------
				hash, _ := chainhash.NewHashFromStr(v.TxID)
				outPoint := wire.NewOutPoint(hash, v.Vout)
				txIn := wire.NewTxIn(outPoint, nil, nil)

				outsu.Add(am)
				tx.AddTxIn(txIn)

				txinPkScript, errInner := hex.DecodeString(v.ScriptPubKey)
				if errInner != nil {
					return nil, errInner
				}
				scripts.Scripts = append(scripts.Scripts, txinPkScript)
			}
		} else {
			break
		}
	}
	pkScripts = append(pkScripts, scripts)
	cmpRe, err := outsu.Cmp(totalTranDynamic)
	if err != nil {
		return nil, err
	}
	if cmpRe == -1 {
		err = errors.ERR_NOT_ENOUGH_COIN
		return
	}
	////////btc cmp end/////////////////////////

	// form----------------return-------------------
	addrf, err := btcutil.DecodeAddress(addrF, c.Env)
	if err != nil {
		return
	}
	pkScriptf, err := txscript.PayToAddrScript(addrf)
	if err != nil {
		return
	}
	outsu.Sub(totalTranDynamic)
	tx.AddTxOut(wire.NewTxOut(outsu.Val().Int64(), pkScriptf))
	//to------------------pay----Must be put last------------
	addrt, errInner := btcutil.DecodeAddress(addrT.Address, c.Env)
	if errInner != nil {
		err = errInner
		return
	}
	pkScriptt, errInner := txscript.PayToAddrScript(addrt)
	if errInner != nil {
		err = errInner
		return
	}
	tx.AddTxOut(wire.NewTxOut(omniPrice.Val().Int64(), pkScriptt))
	//opreturn Set omni information
	c.addOpReturn(tx, addrT.Amount.Val().Uint64())
	return signature.NewBtcOmniTx(tx, pkScripts, c.Env, c.PropertyId, local), nil
}

func (c *UsdtClient) addOpReturn(rawtx *wire.MsgTx, amount uint64) {
	var omni = "omni"
	sOmni := fmt.Sprintf("%08x", omni)
	sPropertyId := fmt.Sprintf("%016X", c.PropertyId)
	sAmount := fmt.Sprintf("%016X", amount)
	msg := sOmni + sPropertyId + sAmount
	bMsg, _ := hex.DecodeString(msg)
	pkScript, _ := txscript.NewScriptBuilder().AddOp(txscript.OP_RETURN).AddData(bMsg).Script()
	rawtx.AddTxOut(wire.NewTxOut(0, pkScript))
}

func (c *UsdtClient) SendTx(txu signature.TxUtil) error {
	txv, err := txu.TxForSend()
	if err != nil {
		return nil
	}
	tx := txv.(*wire.MsgTx)
	txHash, err := c.nmap[txu.Local()].C.SendRawTransaction(tx, false)
	fmt.Println(txHash)
	return err
}
func (c *UsdtClient) GetBalance(address string, token string, local bool) (balance bccoin.CoinAmounter, err error) {
	omniBalance, err := c.nmap[local].C.OmniGetbalance(address, c.PropertyId)
	if err != nil {
		return
	}
	balance, err = bccoin.NewCoinAmount(bccore.BC_USDT, bccore.Token(strconv.FormatInt(int64(c.PropertyId), 10)), omniBalance.Balance)
	return
}
func getBtcFee(rawtx *wire.MsgTx, fee bccoin.CoinAmounter) bccoin.CoinAmounter {
	size := rawtx.SerializeSize() + 200
	realFee, _ := bccoin.NewCoinAmount(bccore.BC_BTC, "", "0")
	realFee.Set(fee)
	if size > 1000 {
		factor, _ := bccoin.NewCoinAmount(bccore.BC_BTC, "", strconv.FormatFloat(float64(size)/1000, 'f', 8, 64))
		realFee.Mul(factor)
	}
	return realFee
}

///////////////////////////////////////////////////internal/////////////////////////////////////////////////////
func (c *UsdtClient) getUnspentByAddress(address string, local bool) (unspents []sebtcjson.ListUnspentResult, err error) {
	btcAdd, err := btcutil.DecodeAddress(address, c.Env)
	if err != nil {
		return
	}
	adds := [1]btcutil.Address{btcAdd}
	unspents, err = c.nmap[local].C.ListUnspentMinMaxAddresses(1, 999999, adds[:])

	return
}

func (c *UsdtClient) addAddressToWallet(address btcutil.Address, local bool) error {
	err := c.checkAddressExists(address, local)
	if err != nil {
		return err
	}

	if err = c.nmap[local].C.ImportAddress(address.String()); err != nil {
		return err
	}
	return nil
}

func (c *UsdtClient) checkAddressExists(address btcutil.Address, local bool) error {
	addrValid, err := c.nmap[local].C.ValidateAddress(address)
	if err != nil {
		return err
	}
	if addrValid.IsWatchOnly {
		return errors.ERR_DATA_EXISTS
	}
	return nil
}

func (c *UsdtClient) ChooseClientNode() (local bool, err error) {

	bc, err := c.C.GetBlockChainInfo()
	if err != nil {
		return false, err
	}
	if bc.Headers-bc.Blocks > int32(c.gap) {
		return false, nil
	}
	return true, nil
}
