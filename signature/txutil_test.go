package signature_test

import (
	"fmt"
	"reflect"

	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil/hdkeychain"

	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/boxproject/boxwallet/mock"
	"github.com/boxproject/boxwallet/signature"
)

type txUHandler struct {
	i        signature.TxUtil
	TypeName string
}

var (
	btc     *signature.BtcTx
	btcName = "*signature.BtcTx"
	eth     *signature.EthTx
	ethName = "*signature.EthTx"
	handler txUHandler
)

var (
	bccache = mock.Cache
)

func (h *txUHandler) LoadService(i signature.TxUtil) error {
	if i != nil {
		h.i = i
	}
	typ := reflect.TypeOf(i)
	h.TypeName = typ.String()
	return nil
}

func (h *txUHandler) GetTx() signature.TxUtil {
	var t signature.TxUtil
	switch h.TypeName {
	case btcName:
		var pkscripts [][]byte
		pub, _ := hdkeychain.NewKeyFromString(mock.Pubkey)
		addr, _ := pub.Address(&chaincfg.RegressionNetParams)

		//prv, _ := hdkeychain.NewKeyFromString(prvkey)
		//addr, _ := prv.Address(&chaincfg.RegressionNetParams)

		originTx := wire.NewMsgTx(wire.TxVersion)
		prevOut := wire.NewOutPoint(&chainhash.Hash{}, ^uint32(0))
		txIn := wire.NewTxIn(prevOut, []byte{txscript.OP_0, txscript.OP_0}, nil)
		originTx.AddTxIn(txIn)
		pkScript, err := txscript.PayToAddrScript(addr)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		txOut := wire.NewTxOut(100000000, pkScript)
		originTx.AddTxOut(txOut)

		originTxHash := originTx.TxHash()

		prevOut = wire.NewOutPoint(&originTxHash, 0)
		txIn = wire.NewTxIn(prevOut, nil, nil)
		pkscripts = append(pkscripts, pkScript)
		txOut = wire.NewTxOut(0, nil)

		m := make(map[string][][]byte)
		m[addr.String()] = pkscripts
		t = signature.NewBtcTx(originTx, m, &chaincfg.RegressionNetParams, true)
		break
	case ethName:
		toAccDef := accounts.Account{
			Address: common.HexToAddress("0x0640287c23c3c3c59388f73cf904ec9277887820"),
		}
		fromAccDef := accounts.Account{
			Address: common.HexToAddress("0x430458dCfd8e88f37f1E60507EC6e0E84aa30118"),
		}
		t = signature.NewEthTx(uint64(0), fromAccDef.Address, toAccDef.Address, big.NewInt(1e18), 3231744, big.NewInt(18000000000), nil, true)
		break
	}
	return t
}

func Test_Info(t *testing.T) {
	handler.LoadService(btc)
	tx := handler.GetTx()
	byteArr, _ := tx.Marshal()
	err := tx.UnMarshal(byteArr)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	t.Log(string(byteArr))
	tos, err := tx.Info("")
	if err != nil {
		t.Fail()
		t.Error(err)
	}

	for _, v := range tos {
		t.Log(v)
	}
}
func Test_FromAddresses(t *testing.T) {
	handler.LoadService(eth)
	tx := handler.GetTx()
	froms := tx.FromAddresses()
	for _, v := range froms {
		t.Log(v)
	}
}
func Test_Sign(t *testing.T) {
	handler.LoadService(eth)
	tx := handler.GetTx()
	byteArr, _ := tx.Marshal()
	tx.UnMarshal(byteArr)
	err := tx.Sign([]string{mock.Prvkey})
	if err != nil {
		t.Fail()
		t.Error(err)
	}
}
func Test_TxForSend(t *testing.T) {
	handler.LoadService(btc)
	tx := handler.GetTx()
	byteArr, err := tx.Marshal()
	if err != nil {
		t.Fail()
		t.Error(err)
		return
	}
	err = tx.UnMarshal(byteArr)
	if err != nil {
		t.Fail()
		t.Error(err)
		return
	}
	err = tx.Sign([]string{mock.Prvkey})
	if err != nil {
		t.Fail()
		t.Error(err)
		return
	}
	by2, err := tx.Marshal()
	if err != nil {
		t.Fail()
		t.Error(err)
		return
	}
	err = tx.UnMarshal(by2)
	if err != nil {
		t.Fail()
		t.Error(err)
		return
	}
	v, err := tx.TxForSend()
	if err != nil {
		t.Fail()
		t.Error(err)
		return
	} else {
		switch v.(type) {
		case *wire.MsgTx:
			t.Log("btctx")
			break
		case *types.Transaction:
			t.Log("ethtx")
			break
		default:
			t.Fail()
			t.Log(reflect.TypeOf(v))
			break
		}
	}

}

func Test_(t *testing.T) {
	handler.LoadService(btc)
	//tx := handler.GetTx()
	switch handler.TypeName {
	case btcName:

	}
}
