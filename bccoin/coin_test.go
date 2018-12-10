package bccoin_test

import (
	"testing"

	"github.com/boxproject/boxwallet/bccoin"
	"github.com/boxproject/boxwallet/bccore"
	"github.com/boxproject/boxwallet/errors"
	"github.com/boxproject/boxwallet/mock"
)

var (
	cache     = mock.Cache
	ca1, _    = bccoin.NewCoinAmount(bccore.BC_BTC, "", "123.1000")
	ca2, _    = bccoin.NewCoinAmount(bccore.BC_BTC, "", "123.1000")
	ca_sub, _ = bccoin.NewCoinAmount(bccore.BC_BTC, "", "0")
	ca_mul, _ = bccoin.NewCoinAmount(bccore.BC_BTC, "", "1515361000000")
	ca_add, _ = bccoin.NewCoinAmount(bccore.BC_BTC, "", "246.2")

	ca3, _ = bccoin.NewCoinAmount(bccore.BC_ETH, "", "2123.1000")
)

func reload() {
	ca1, _ = bccoin.NewCoinAmount(bccore.BC_BTC, "", "123.1000")
	ca2, _ = bccoin.NewCoinAmount(bccore.BC_BTC, "", "123.1000")
	ca_sub, _ = bccoin.NewCoinAmount(bccore.BC_BTC, "", "0")
	ca_mul, _ = bccoin.NewCoinAmount(bccore.BC_BTC, "", "1515361000000")
	ca_add, _ = bccoin.NewCoinAmount(bccore.BC_BTC, "", "246.2")

	ca3, _ = bccoin.NewCoinAmount(bccore.BC_ETH, "", "2123.1000")
}

func TestCoinCache_GetCoinInfo(t *testing.T) {
	cm, err := cache.GetCoinInfo(bccore.BC_BTC, "")
	if err != nil {
		t.Fail()
		t.Error(err)
	} else {
		t.Log(cm)
	}
}

func TestNewCoinAmount(t *testing.T) {
	ca, err := bccoin.NewCoinAmount(bccore.BC_BTC, "", "1000")
	if err != nil {
		t.Fail()
		t.Error(err)
	} else {

		t.Log(ca.String(), ca.Val(), ca.Sign())
	}
}

func TestEthCoin(t *testing.T) {
	//amount := big.NewInt(60000000000000000)
	ca, _ := bccoin.NewCoinAmountFromInt(bccore.BC_ETH, "", 60000000000000000)
	t.Log(ca.String())
}

func TestCoinAmount_Calculate(t *testing.T) {

	//1
	r, err := ca1.Cmp(ca2)
	if r != 0 || err != nil {
		t.Fail()
		t.Error("test 1 failed:", err)
	} else {
		t.Log("test 1=>success")
	}
	t.Log("cmp test 1：", r)

	//2
	r, err = ca1.Cmp(ca3)
	if err != errors.ERR_DIFF_UNIT {
		t.Fail()
		t.Error("test 2 failed:", err)
	} else {
		t.Log("test 2=>success")
	}
	t.Log("cmp test 2：", r)

	//3
	err = ca1.Sub(ca2)
	if ca1.String() != ca_sub.String() || err != nil {
		t.Fail()
		t.Error("test 3 failed:", err)
	} else {
		t.Log("test 3=>success")
	}
	t.Logf("sub test：The target：%s\r\nactual:%s\r\n", ca_sub.String(), ca1.String())
	reload()
	//4
	err = ca1.Add(ca2)
	if ca1.String() != ca_add.String() || err != nil {
		t.Fail()
		t.Error("test 4 failed:", err)
	} else {
		t.Log("test 4=>success")
	}
	t.Logf("add test：The target：%s\r\nactual:%s\r\n", ca_add.String(), ca1.String())
	reload()
	//5
	err = ca1.Mul(ca3)
	if err != errors.ERR_DIFF_UNIT {
		t.Fail()
		t.Error("test 5 failed:", err)
	} else {
		t.Log("test 5=>success")
	}
	t.Log("cmp test 5：", err)
	reload()
	//6
	err = ca1.Mul(ca2)
	if ca1.String() != ca_mul.String() || err != nil {
		t.Fail()
		t.Error("test 6 failed:", err)
	} else {
		t.Log("test 6=>success")
	}
	t.Logf("mul test：The target：%s\r\nactual:%s\r\n", ca_mul.String(), ca1.String())
	reload()
}

func TestCoinCache_GetAll(t *testing.T) {
	ch, err := cache.GetAll(bccore.BC_BTC, false)
	if err != nil {
		t.Fail()
		t.Error(err)
	}
	for c := range ch {
		t.Log(c.Token, c.CT)
	}
}
