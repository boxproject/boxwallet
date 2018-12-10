package signature

import (
	"encoding/json"
	"math/big"

	"github.com/boxproject/boxwallet/errors"
	"github.com/ltcsuite/ltcd/btcec"
	"github.com/ltcsuite/ltcd/chaincfg"
	"github.com/ltcsuite/ltcd/txscript"
	"github.com/ltcsuite/ltcd/wire"
	"github.com/ltcsuite/ltcutil"
	"github.com/ltcsuite/ltcutil/hdkeychain"
)

type LtcTx struct {
	T          *wire.MsgTx
	PkScripts  []*PKscripts
	Env        *chaincfg.Params
	Signed     bool
	PropertyId int64
	Loc        bool
}

func NewLtcTx(tx *wire.MsgTx, pkScripts []*PKscripts, env *chaincfg.Params, local bool) *LtcTx {
	return &LtcTx{T: tx, PkScripts: pkScripts, Env: env, Loc: local}
}
func (tx *LtcTx) FromAddresses() (froms []string) {
	froms = make([]string, 0, len(tx.PkScripts))
	for k, _ := range tx.PkScripts {
		froms = append(froms, tx.PkScripts[k].Address)
	}
	return
}

func (tx *LtcTx) Info(addrFrom string) (to []*AddressAmount, err error) {
	to = make([]*AddressAmount, len(tx.T.TxOut), len(tx.T.TxOut))
	i := 0
	for _, v := range tx.T.TxOut {

		if txscript.IsUnspendable(v.PkScript) {

		} else {
			scriptAddr := v.PkScript[3 : v.PkScript[2]+3]
			addr, err := ltcutil.NewAddressPubKeyHash(scriptAddr, tx.Env)
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
func (tx *LtcTx) IsSign() bool {
	return tx.Signed
}
func (tx *LtcTx) TxId() string {
	if !tx.IsSign() {
		return ""
	}
	return tx.T.TxHash().String()
}

//签名
func (tx *LtcTx) Sign(privKeys []string) error {

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
func (tx *LtcTx) TxForSend() (v interface{}, err error) {
	if tx.IsSign() {
		return tx.T, nil
	}
	return nil, errors.ERR_TX_WITHOUT_SGIN
}
func (tx *LtcTx) Local() bool {
	return tx.Loc

}
func (tx *LtcTx) Marshal() ([]byte, error) {
	return json.Marshal(tx)
}

func (tx *LtcTx) UnMarshal(data []byte) error {
	return json.Unmarshal(data, tx)
}
