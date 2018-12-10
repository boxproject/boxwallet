package client_test

import (
	"reflect"
	"testing"

	"log"

	"github.com/btcsuite/btcd/chaincfg"

	"time"

	"github.com/boxproject/boxwallet/bccoin"
	"github.com/boxproject/boxwallet/bccore"
	"github.com/boxproject/boxwallet/bctrans/client"
	"github.com/boxproject/boxwallet/mock"
	"github.com/boxproject/boxwallet/signature"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/core/types"
	ltcwire "github.com/ltcsuite/ltcd/wire"
)

type cliHandler struct {
	i        client.Walleter
	TypeName string
}

var (
	btc       = mock.Btc
	btcName   = "*client.BtcClient"
	usdt      = mock.Usdt
	usdtName  = "*client.UsdtClient"
	eth       = mock.Eth
	ethName   = "*client.EthClient"
	erc20     = mock.Erc20
	erc20Name = "*client.Erc20Client"
	ltc       = mock.Ltc
	ltcName   = "*client.LtcClient"

	handler cliHandler
	cache   = mock.Cache
)

func (h *cliHandler) LoadService(i client.Walleter) error {
	if i != nil {
		h.i = i
	}
	typ := reflect.TypeOf(i)
	h.TypeName = typ.String()
	return nil
}

func Test_ImportAddress(t *testing.T) {
	handler.LoadService(btc)
	var (
		address string
	)
	switch handler.TypeName {
	case btcName:
		address = "mzwAMJzjdReh4zuL4a1tk8T6RUS4TCffnz"
		break
	case usdtName:
		address = mock.Addr
		break
	case ethName:
		address = mock.AddrEth
		break
	case erc20Name:
		address = mock.AddrEth2
		break
	case ltcName:
		address = mock.Addr
		break
	}
	err := handler.i.ImportAddress(address, true)
	if err != nil {
		t.Fail()
		t.Error(err)
	}
}
func Test_GetBalance(t *testing.T) {
	handler.LoadService(ltc)
	var (
		address string
		token   string
	)
	switch handler.TypeName {
	case btcName:
		address = mock.Addr
		break
	case usdtName:
		address = mock.Addr
		break
	case ethName:
		address = mock.AddrEth
		break
	case erc20Name:
		address = mock.AddrEth2
		token = "0x6A671140c983EbA5636A32f23f1f0d1616c0fff6"
		break
	case ltcName:
		address = mock.Addr
		break
	}

	balance, err := handler.i.GetBalance(address, token, true)
	if err != nil {
		t.Fail()
		t.Error(err)
	} else {
		t.Log(balance)
	}
}

func TestGetBlance2(t *testing.T) {
	handler.LoadService(ltc)
	var (
		token string
	)
	addresses := []string{
		mock.Addr,
		mock.Addr2,
		mock.Addr3,
		mock.Addr4,
	}
	sum, _ := bccoin.NewCoinAmount(bccore.BC_LTC, "", "0")
	for _, v := range addresses {
		balance, err := handler.i.GetBalance(v, token, true)
		if err != nil {
			t.Fail()
			t.Error(err)
		} else {
			t.Log(balance)
		}
		sum.Add(balance)
	}
	t.Log("sum:", sum.String())
}

func Test_SendTx(t *testing.T) {
	handler.LoadService(ltc)

	var (
		bcs   []*bccoin.AddressAmount
		txu   signature.TxUtil
		err   error
		froms []string
		prvs  []string
	)
	time1 := time.Now()
	switch handler.TypeName {
	case btcName:
		froms = append(froms, "mwHZxRNXj9rQxjDAr2FCJnAjKJDT94h7P2",
			"mfiomhtor57G5dYZaiMpaDjQAQPMSBmayy",
			"mre55thoQCbpaUw1MhovdfBM6xo7CsyNFF",
			"myWzxpxWaKtD2MWdj8a4kr2hRAAGgUQPSw",
			"myGJZQ69Cpu8RctpHf4oh1UXw68gy7ZuVk")
		prvs = append(prvs, mock.Prvkey)
		bca, _ := bccoin.NewCoinAmount(bccore.BC_BTC, "", "100")

		bc1 := &bccoin.AddressAmount{
			Address: mock.Addr2,
			Amount:  bca,
		}
		//fee, _ := bccoin.NewCoinAmount(bccore.BC_BTC, "", "0.0001")
		//for i := 0; i < 30; i++ {
		bcs = append(bcs, bc1)
		//}
		txu, err = handler.i.CreateTx(froms, "", bcs, 1.1)

		break
	case usdtName:
		froms = append(froms, mock.Addr)
		prvs = append(prvs, mock.Prvkey)
		bca, _ := bccoin.NewCoinAmount(bccore.BC_USDT, mock.PropertyId, "15")
		bc1 := &bccoin.AddressAmount{
			Address: mock.Addr2,
			Amount:  bca,
		}
		//fee, _ := bccoin.NewCoinAmount(bccore.BC_BTC, "", "0.0001")
		bcs = append(bcs, bc1)
		txu, err = handler.i.CreateTx(froms, "", bcs, 1.1)
		break
	case ethName:
		froms = []string{mock.AddrEth}
		prvs = append(prvs, mock.Prvkey)
		bca, _ := bccoin.NewCoinAmount(bccore.BC_ETH, "", "1")
		bc1 := &bccoin.AddressAmount{
			Address: mock.AddrEth2,
			Amount:  bca,
		}
		//fee, _ := bccoin.NewCoinAmount(bccore.BC_ETH, "", "0.0001")
		bcs = append(bcs, bc1)
		txu, err = handler.i.CreateTx(froms, "", bcs, 1.1)

		break
	case erc20Name:
		froms = []string{mock.AddrEth}
		prvs = append(prvs, mock.Prvkey)
		bca, _ := bccoin.NewCoinAmount(bccore.BC_ERC20, "0x6A671140c983EbA5636A32f23f1f0d1616c0fff6", "1")
		bc1 := &bccoin.AddressAmount{
			Address: mock.AddrEth2,
			Amount:  bca,
		}
		bcs = append(bcs, bc1)
		txu, err = handler.i.CreateTx(froms, "0x6A671140c983EbA5636A32f23f1f0d1616c0fff6", bcs, 1.1)
		break
	case ltcName:
		froms = append(froms, mock.Addr, mock.Addr2, mock.Addr3, mock.Addr4)
		prvs = append(prvs, mock.Prvkey, mock.Prvkey2, mock.Prvkey4, mock.Prvkey3)
		bca, _ := bccoin.NewCoinAmount(bccore.BC_LTC, "", "155")

		bc1 := &bccoin.AddressAmount{
			Address: "mqK6VNApNrFgYAYoqTtDQjZSudBdrKVy8r",
			Amount:  bca,
		}
		//for i := 0; i < 3999; i++ {
		bcs = append(bcs, bc1)
		//}
		//fee, _ := bccoin.NewCoinAmount(bccore.BC_LTC, "", "0.0001")
		txu, err = handler.i.CreateTx(froms, "", bcs, 1.0)

		break
	}
	if err != nil {
		t.Fail()
		t.Error(err)
		return
	}
	t.Log("local:", txu.Local())
	txout, err := txu.Info("")
	if err != nil {
		t.Fail()
		t.Error(err)
		return
	}
	for _, v := range txout {
		t.Log(v.Address, v.Amount.String(), v.Amount.Sign())
	}

	//return //断点测试
	byteArr, err := txu.Marshal()
	if err != nil {
		t.Fail()
		t.Error(err)
		return
	}
	err = txu.UnMarshal(byteArr)
	if err != nil {
		t.Fail()
		t.Error(err)
		return
	}
	//return
	t.Log("first:", txu.TxId())
	err = txu.Sign(prvs)
	if err != nil {
		t.Fail()
		t.Error(err)
		return
	}
	t.Log("second:", txu.TxId())

	v, err := txu.TxForSend()
	if err != nil {
		t.Fail()
		t.Error(err)
	}
	t.Log(v)
	switch v.(type) {
	case *wire.MsgTx:
		t.Log("btctx")
		break
	case *types.Transaction:
		t.Log("ethtx")
		break
	case *ltcwire.MsgTx:
		t.Log("ltctx")
		break
	default:
		t.Fail()
		t.Log(reflect.TypeOf(v))
		return
	}
	time2 := time.Now()
	//return
	log.Println(txu.TxId())
	err = handler.i.SendTx(txu)
	if err != nil {
		t.Fail()
		t.Error(err)
		return
	}
	time3 := time.Now()
	t.Log(time3.Sub(time1).Nanoseconds())
	t.Log(time3.Sub(time2).Nanoseconds())
}

func Benchmark_SendTx(b *testing.B) {
	b.ReportAllocs()
	handler.LoadService(eth)
	for i := 0; i < b.N; i++ { //use b.N for looping
		var (
			bcs   []*bccoin.AddressAmount
			txu   signature.TxUtil
			err   error
			froms []string
			prvs  []string
		)

		switch handler.TypeName {
		case btcName:
			froms = append(froms, mock.Addr)
			prvs = append(prvs, mock.Prvkey)
			bca, _ := bccoin.NewCoinAmount(bccore.BC_BTC, "", "1")

			bc1 := &bccoin.AddressAmount{
				Address: mock.Addr2,
				Amount:  bca,
			}
			//fee, _ := bccoin.NewCoinAmount(bccore.BC_BTC, "", "0.0001")
			for i := 0; i < 30; i++ {
				bcs = append(bcs, bc1)
			}
			txu, err = handler.i.CreateTx(froms, "", bcs, 1.1)

			break
		case usdtName:
			froms = append(froms, mock.Addr)
			prvs = append(prvs, mock.Prvkey)
			bca, _ := bccoin.NewCoinAmount(bccore.BC_USDT, mock.PropertyId, "15")
			bc1 := &bccoin.AddressAmount{
				Address: mock.Addr2,
				Amount:  bca,
			}
			//fee, _ := bccoin.NewCoinAmount(bccore.BC_BTC, "", "0.0001")
			bcs = append(bcs, bc1)
			txu, err = handler.i.CreateTx(froms, "", bcs, 1.1)
			break
		case ethName:
			froms = []string{mock.AddrEth}
			prvs = append(prvs, mock.Prvkey)
			bca, _ := bccoin.NewCoinAmount(bccore.BC_ETH, "", "1")
			bc1 := &bccoin.AddressAmount{
				Address: mock.AddrEth2,
				Amount:  bca,
			}
			//fee, _ := bccoin.NewCoinAmount(bccore.BC_ETH, "", "0.0001")
			bcs = append(bcs, bc1)
			txu, err = handler.i.CreateTx(froms, "", bcs, 1.1)

			break
		case erc20Name:
			froms = []string{mock.AddrEth}
			prvs = append(prvs, mock.Prvkey)
			bca, _ := bccoin.NewCoinAmount(bccore.BC_ERC20, "0x6A671140c983EbA5636A32f23f1f0d1616c0fff6", "1")
			bc1 := &bccoin.AddressAmount{
				Address: mock.AddrEth2,
				Amount:  bca,
			}
			bcs = append(bcs, bc1)
			txu, err = handler.i.CreateTx(froms, "0x6A671140c983EbA5636A32f23f1f0d1616c0fff6", bcs, 1.1)
			break
		case ltcName:
			froms = append(froms, mock.Addr)
			prvs = append(prvs, mock.Prvkey)
			bca, _ := bccoin.NewCoinAmount(bccore.BC_LTC, "", "0.1")

			bc1 := &bccoin.AddressAmount{
				Address: mock.Addr3,
				Amount:  bca,
			}
			for i := 0; i < 30; i++ {
				bcs = append(bcs, bc1)
			}
			//fee, _ := bccoin.NewCoinAmount(bccore.BC_LTC, "", "0.0001")
			txu, err = handler.i.CreateTx(froms, "", bcs, 1.1)

			break
		}
		if err != nil {
			b.Log(err)
			return
		}
		_, err = txu.Info("")
		if err != nil {
			b.Log(err)
			return
		}
		//return //断点测试
		byteArr, err := txu.Marshal()
		if err != nil {
			b.Log(err)
			return
		}
		err = txu.UnMarshal(byteArr)
		if err != nil {
			b.Log(err)
			return
		}

		err = txu.Sign(prvs)
		if err != nil {
			b.Log(err)
			return
		}
		_, err = txu.TxForSend()
		if err != nil {
			b.Log(err)
		}
		/*err = handler.i.SendTx(txu)
		if err != nil {
			b.Log(err)
			return
		}*/
	}
}

func Test_(t *testing.T) {
	handler.LoadService(btc)
	switch handler.TypeName {
	case btcName:

	}
}

func TestTxSize(t *testing.T) {
	handler.LoadService(btc)
	var (
		bcs   []*bccoin.AddressAmount
		txu   signature.TxUtil
		froms []string
		prvs  []string
	)
	froms = append(froms, mock.Addr, mock.Addr2, mock.Addr3, mock.Addr4)
	prvs = append(prvs, mock.Prvkey, mock.Prvkey2)
	bca, _ := bccoin.NewCoinAmount(bccore.BC_BTC, "", "2")

	bc1 := &bccoin.AddressAmount{
		Address: mock.Addr3,
		Amount:  bca,
	}
	//fee, _ := bccoin.NewCoinAmount(bccore.BC_BTC, "", "0.0001")
	bcs = append(bcs, bc1)
	txu, _ = handler.i.CreateTx(froms, "", bcs, 1)
	txu.Sign([]string{mock.Prvkey})
	txv, err := txu.TxForSend()
	if err != nil {
		return
	}
	tx := txv.(*wire.MsgTx)

	size := float64(148*len(tx.TxIn) + len(tx.TxOut)*34 + 10)
	size2 := tx.SerializeSize()
	t.Log(size, size2)
}

func TestEncodeAddress(t *testing.T) {
	btcAdd, _ := btcutil.DecodeAddress("mtanjpKRbgMEJznptS27Z8vfw8fQb7dHCc", &chaincfg.RegressionNetParams)
	t.Log(btcAdd.EncodeAddress(), btcAdd.String())

}
