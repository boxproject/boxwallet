package client

import (
	"github.com/boxproject/boxwallet/bccoin"
	"github.com/boxproject/boxwallet/bccore"
	"github.com/boxproject/boxwallet/bckey"
	"github.com/boxproject/boxwallet/signature"
)

type TxInfo struct {
	T       bccore.BloclChainType
	TxBytes []byte
}

type Walleter interface {
	GetBalance(address string, token string, local bool) (balance bccoin.CoinAmounter, err error)
	CreateTx(addrsFrom []string, token string, addrsTo []*bccoin.AddressAmount, feeCeo float64) (txu signature.TxUtil, err error)
	SendTx(txu signature.TxUtil) error
	GetNewAddress() (bckey.GenericKey, error)
	ImportAddress(address string, local bool) error
	ChooseClientNode() (local bool, err error)
}
