package daemon_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/boxproject/boxwallet/daemon"
)

var ltcDaemon *daemon.LtcDaemon

func init() {

	ltcDaemon = daemon.NewLtcDaemon(uint64(6), 6, time.Second)

}
func TestLtcDaemon_GetTransaction(t *testing.T) {
	//0xe042795078c374a6c9ae929974757602d72564d8f360855b39560291ac49f876
	txId := &daemon.TxIdInfo{
		TxId:   "55e9cf10a5d67e48731316d01e8050df4b1ba3a2a9bbd5e23b257dfe35efa2bf",
		BlockH: big.NewInt(1080),
	}
	txi, err := ltcDaemon.GetTransaction(txId, false)
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
	t.Log(txi.Fee.String())
	t.Log(txi.Cnfm)
}
func TestLtcDaemon_GetBlockHeight(t *testing.T) {
	height, h2, err := ltcDaemon.GetBlockHeight()
	if err != nil {
		t.Fail()
		t.Error(err)
	} else {
		t.Log(height)
		t.Log(h2)
	}
}
func TestLtcDaemon_GetTransaction2(t *testing.T) {
	bk, _ := ltcDaemon.GetBlockInfo(big.NewInt(1529848))
	txInfos := bk.AnalyzeBlock([]string{"LXtwa3rx9JExquao2q23MXcSfvEmEfnrFv"})
	for _, v := range txInfos {
		txi, err := ltcDaemon.GetTransaction(v, false)
		if err != nil {
			t.Fail()
			return
		}
		t.Log(txi.TxId)
		for _, vi := range txi.In {
			t.Log(vi.Amt.Uint64(), vi.Addr)
		}
	}

}
