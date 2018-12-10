package daemon_test

import (
	"testing"

	"math/big"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"time"

	"strings"

	"github.com/boxproject/boxwallet/bccoin"
	"github.com/boxproject/boxwallet/bccore"
	"github.com/boxproject/boxwallet/daemon"
	_ "github.com/boxproject/boxwallet/mock"
)

var btcDaemon *daemon.BtcDaemon

func init() {
	btcDaemon = daemon.NewBtcDaemon(uint64(6), 64, time.Second)

}

func TestBtcDaemon_GetBlockHeight(t *testing.T) {
	height, h2, err := btcDaemon.GetBlockHeight()
	if err != nil {
		t.Fail()
		t.Error(err)
	} else {
		t.Log(height)
		t.Log(h2)
	}
	btcDaemon.C.GetBestBlock()
	a, _ := btcDaemon.C.GetBlockChainInfo()
	t.Log(a)
}

func TestBtcDaemon_GetBlockInfo(t *testing.T) {
	si, err := btcDaemon.GetBlockInfo(big.NewInt(283))
	if err != nil {
		t.Fail()
		t.Error(err)
		return
	}
	txIds := si.AnalyzeBlock([]string{"mhAfGecTPa9eZaaNkGJcV7fmUPFi3T2Ki8"})
	t.Log(txIds)
}

func TestBtcDaemon_GetTransaction(t *testing.T) {

	txi, err := btcDaemon.GetTransaction(&daemon.TxIdInfo{"6fae0307687ec2201768451a30d5ae308ea05169ff0a6514ed9b124a316f753b", nil}, true)
	if err != nil {
		t.Fail()
		t.Error(err)
		return
	}

	for _, v := range txi.In {
		t.Log("in:", v.Amt.String(), v.Addr)
	}
	for _, v := range txi.Out {
		t.Log("out:", v.Amt.String(), v.Addr)
	}
	t.Log(txi.Cnfm)
}

func TestBtcDaemon_spTTT(t *testing.T) {
	opReturnStr := "OP_RETURN"
	nulldata := "nulldata"
	txId, _ := chainhash.NewHashFromStr("e4c0eacc16dd29d7eccff66485c4851475681fad77f8fd954ced282baf6cebc0")
	txInfo, _ := btcDaemon.C.GetRawTransactionVerbose(txId)

	for _, v := range txInfo.Vout {
		if v.ScriptPubKey.Type == nulldata {
			arr := strings.Split(v.ScriptPubKey.Asm, " ")
			if len(arr) != 2 && arr[0] != opReturnStr {
				continue
			}
			omniRes, err := btcDaemon.C.OmniGetTransaction(txId)

			if err == nil && omniRes != nil && omniRes.Propertyid == int64(btcDaemon.UsdtId) {
				amt, _ := bccoin.NewCoinAmount(bccore.BC_BTC, bccore.Token("2147483656"), omniRes.Amount)
				t.Log(amt.String())
			} else {
				//TODO LOG

			}

		}
		/*for _, v := range txInfo.Vout {
			t.Log(string(v.ScriptPubKey.Hex), v.ScriptPubKey.Type, "0-0-0--0-0", v.ScriptPubKey.Asm)
			t.Log(v.ScriptPubKey.Type)
			if v.ScriptPubKey.Type == "OP_RETURN" {
				omniRes, err := mock.Usdt.C.OmniGetTransaction(txId)

				if err == nil && omniRes != nil && omniRes.Propertyid == int64(mock.Usdt.PropertyId) {
					//amt, _ := big.NewInt(0).SetString(omniRes.Amount, 10)

				} else {
					//TODO LOG

				}

			} else {
				t.Log(v.ScriptPubKey.Type, v.ScriptPubKey.Asm, v.ScriptPubKey.ReqSigs, v.ScriptPubKey.Addresses)
				t.Log(len(v.ScriptPubKey.Addresses))

			}
		}*/
	}
}
