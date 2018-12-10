/*
client Global entry
*/
package bctrans

import (
	"log"

	"github.com/boxproject/boxwallet/bccoin"
	"github.com/boxproject/boxwallet/bccore"
	"github.com/boxproject/boxwallet/bctrans/client"
	"github.com/boxproject/boxwallet/errors"
	"github.com/boxproject/boxwallet/official"
	"github.com/boxproject/boxwallet/pipeline"
	"github.com/boxproject/boxwallet/signature"
	"github.com/boxproject/boxwallet/util"
)

type Trans struct {
	Bct bccore.BloclChainType
	client.Walleter
	Pipe *pipeline.Pipeline
	lock bool //If you want to control the blocking of the address, pass false here
}

func NewTrans(bct bccore.BlockChainSign, lock bool) (ts *Trans, err error) {
	ts = &Trans{
		Pipe: &pipeline.Pipeline{},
		lock: lock,
	}
	switch bct {
	case bccore.STR_BTC:
		ts.Bct = bccore.BC_BTC
		ts.Walleter, err = client.GetBtcClientIntance()
		return
	case bccore.STR_ETH:
		ts.Bct = bccore.BC_ETH
		ts.Walleter, err = client.GetEthClientIntance()
		return
	case bccore.STR_ERC20:
		ts.Bct = bccore.BC_ERC20
		ts.Walleter, err = client.GetErc20ClientIntance()
		return
	case bccore.STR_USDT:
		ts.Bct = bccore.BC_USDT
		ts.Walleter, err = client.GetUsdtClientInstance()
		return
	case bccore.STR_LTC:
		ts.Bct = bccore.BC_LTC
		ts.Walleter, err = client.GetLtcClientIntance()
		return
	}
	return nil, errors.ERR_PARAM_NOT_VALID
}

func (t *Trans) GetBalance(address string, token string) (balance bccoin.CoinAmounter, err error) {
	return t.Walleter.GetBalance(address, token, true)
}

func (t *Trans) CreateTx(addrsFrom []string, token string, addrsTo []*bccoin.AddressAmount, feeCeo float64) (uuid string, txu signature.TxUtil, err error) {
	if feeCeo < 1 || feeCeo > 10 {
		return "", nil, errors.ERR_PARAM_NOT_VALID
	}
	if t.Pipe.AddressExist(t.Bct, addrsFrom) {
		return "", nil, errors.ERR_ADDRESS_QUEUE_BLOCKED
	}

	txu, err = t.Walleter.CreateTx(addrsFrom, token, addrsTo, feeCeo)
	if err != nil {
		return "", nil, err
	}
	if !txu.Local() {
		log.Println("official createtx successed!!! type:", t.Bct)
	}
	if t.lock {
		uuid = util.GetUUid()
		t.Pipe.Send(t.Bct, uuid, addrsFrom)
	}
	return
}

func (t *Trans) SendTx(txu signature.TxUtil, uuid string) error {
	err := t.Walleter.SendTx(txu)
	if err != nil {
		return err
	}
	if t.lock && txu.Local() {
		ch, err := t.Pipe.CheckSend(t.Bct, txu.TxId(), uuid)
		if err != nil {
			return err
		}
		result := false
		for c := range ch {
			result = c
		}
		if result {
			return nil
		} else {
			return errors.ERR_TX_END_NOT_NORMAL
		}
	} else if !txu.Local() {
		official.ListenTx(t.Bct, txu.TxId())
		log.Println("official sendtx successed!!! type:", t.Bct)
	}
	return nil
}
