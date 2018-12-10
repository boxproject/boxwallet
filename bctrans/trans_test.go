package bctrans_test

import (
	"testing"

	"github.com/boxproject/boxwallet/bccoin"
	"github.com/boxproject/boxwallet/bccore"
	"github.com/boxproject/boxwallet/bctrans"
	"github.com/boxproject/boxwallet/mock"
)

func TestTrans(t *testing.T) {
	//go daemon.Start()
	trans, err := bctrans.NewTrans(bccore.STR_BTC)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	var (
		bcs   []*bccoin.AddressAmount
		froms []string
		prvs  []string
	)
	froms = append(froms, mock.Addr)
	prvs = append(prvs, mock.Prvkey, mock.Prvkey2)
	bca, _ := bccoin.NewCoinAmount(bccore.BC_BTC, "", "1")
	bc1 := &bccoin.AddressAmount{
		Address: mock.Addr3,
		Amount:  bca,
	}
	//fee, _ := bccoin.NewCoinAmount(bccore.BC_BTC, "", "0.0001")
	bcs = append(bcs, bc1)
	uuid, txu, err := trans.CreateTx(froms, "", bcs, 1.1)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	by, err := txu.Marshal()
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	err = txu.UnMarshal(by)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	err = txu.Sign([]string{mock.Prvkey})
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	err = trans.SendTx(txu, uuid)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
}
