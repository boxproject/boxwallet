package signature

import (
	"strconv"

	"github.com/btcsuite/btcd/btcec"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"

	"encoding/json"

	"strings"

	"math/big"

	"github.com/boxproject/boxwallet/errors"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil/hdkeychain"
)

type BtcTx struct {
	T          *wire.MsgTx
	PkScripts  []*PKscripts
	Env        *chaincfg.Params
	Signed     bool
	PropertyId int64 //only for omni
	Loc        bool  //is local
}
type PKscripts struct {
	Address string
	Scripts [][]byte
}

var (
	usdtType = "6f6d6e69"
	opreturn = "OP_RETURN"
)

func NewBtcTx(tx *wire.MsgTx, pkScripts []*PKscripts, env *chaincfg.Params, local bool) *BtcTx {
	return &BtcTx{T: tx, PkScripts: pkScripts, Env: env, Loc: local}
}
func NewBtcOmniTx(tx *wire.MsgTx, pkScripts []*PKscripts, env *chaincfg.Params, propertyId int, local bool) *BtcTx {
	return &BtcTx{T: tx, PkScripts: pkScripts, Env: env, PropertyId: int64(propertyId), Loc: local}
}
func (tx *BtcTx) FromAddresses() (froms []string) {
	froms = make([]string, 0, len(tx.PkScripts))
	for k, _ := range tx.PkScripts {
		froms = append(froms, tx.PkScripts[k].Address)
	}
	return
}

func (tx *BtcTx) Info(addrFrom string) (to []*AddressAmount, err error) {
	to = make([]*AddressAmount, len(tx.T.TxOut), len(tx.T.TxOut))
	i := 0
	for _, v := range tx.T.TxOut {

		if txscript.IsUnspendable(v.PkScript) {
			opstr, err := txscript.DisasmString(v.PkScript)
			if err != nil {
				return nil, err
			}
			//6f6d6e69 0000 0000 0000001f 0000001ac68eac4b
			//
			//6f6d6e69 => omni
			//
			//0000 => version：0
			//
			//0000 => txtype：Simple Send（0）
			//
			//0000001f = 31 => token：USDT
			//
			//0000001ac68eac4b = 115000388683 = 1150.00388683 => transfer amount
			//
			//OP_RETURN 6f6d6e69 0000 0000 80000003 0000000008f0d180
			arr := strings.Split(opstr, " ")
			if len(arr) != 2 {
				return nil, errors.ERR_PARAM_NOT_VALID
			}
			omniNum, err := strconv.ParseInt(arr[1][16:24], 16, 64)
			if err != nil {
				return nil, err
			}
			if arr[0] == opreturn && arr[1][:8] == usdtType && omniNum == tx.PropertyId {
				amount, err := strconv.ParseInt(arr[1][24:], 16, 64)
				if err != nil {
					return nil, err
				}
				/*s := strconv.FormatInt(omniNum, 10)
				amountOmni, err := bccoin.NewCoinAmountFromInt(bccore.BC_USDT, bccore.Token(s), amount)*/
				if err != nil {
					return nil, err
				}
				aa := &AddressAmount{
					Address: to[i-1].Address,
					Amount:  big.NewInt(amount),
				}
				to[i] = aa
			} else {
				return nil, errors.ERR_PARAM_NOT_VALID
			}
		} else {
			scriptAddr := v.PkScript[3 : v.PkScript[2]+3]
			addr, err := btcutil.NewAddressPubKeyHash(scriptAddr, tx.Env)
			if err != nil {
				return nil, err
			}
			if addr.String() == addrFrom {
				continue
			}

			aa := &AddressAmount{
				Address: addr.String(),
				Amount:  big.NewInt(0).SetInt64(v.Value),
			}
			to[i] = aa
		}
		i++
	}
	return
}
func (tx *BtcTx) IsSign() bool {
	return tx.Signed
}
func (tx *BtcTx) TxId() string {
	if !tx.IsSign() {
		return ""
	}
	return tx.T.TxHash().String()
}

func (tx *BtcTx) Sign(privKeys []string) error {

	keyMap := make(map[string]*btcec.PrivateKey)
	for _, v := range privKeys {
		hdprv, err := hdkeychain.NewKeyFromString(v)
		if err != nil {
			return err
		}
		if !hdprv.IsPrivate() {
			return errors.ERR_PARAM_NOT_VALID
		}
		prv, err := hdprv.ECPrivKey()
		if err != nil {
			return err
		}

		addr, err := hdprv.Address(tx.Env)
		if err != nil {
			return err
		}
		keyMap[addr.String()] = prv
	}
	sli := []struct {
		script []byte
		prv    *btcec.PrivateKey
	}{}
	for _, v := range tx.PkScripts {
		for _, vi := range v.Scripts {
			prv := keyMap[v.Address]
			if prv == nil {
				return errors.ERR_PARAM_NOT_VALID
			}
			tmpV := vi
			sli = append(sli, struct {
				script []byte
				prv    *btcec.PrivateKey
			}{script: tmpV, prv: prv})
		}
	}
	for i, _ := range tx.T.TxIn {
		script, err := txscript.SignatureScript(tx.T, i, sli[i].script, txscript.SigHashAll, sli[i].prv, true)
		//script, err := txscript.SignTxOutput(runenv, tx, i, pkScripts[i], txscript.SigHashAll, txscript.KeyClosure(lookupKey), nil, nil)
		if err != nil {
			return err
		}
		tx.T.TxIn[i].SignatureScript = script
		vm, err := txscript.NewEngine(sli[i].script, tx.T, i,
			txscript.StandardVerifyFlags, nil, nil, -1)
		if err != nil {
			return err
		}
		err = vm.Execute()
		if err != nil {
			return err
		}
	}
	tx.Signed = true
	return nil
}
func (tx *BtcTx) TxForSend() (v interface{}, err error) {
	if tx.IsSign() {
		return tx.T, nil
	}
	return nil, errors.ERR_TX_WITHOUT_SGIN
}
func (tx *BtcTx) Local() bool {
	return tx.Loc

}
func (tx *BtcTx) Marshal() ([]byte, error) {
	return json.Marshal(tx)
}

func (tx *BtcTx) UnMarshal(data []byte) error {
	return json.Unmarshal(data, tx)
}
