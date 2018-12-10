package signature

import (
	"math/big"
)

type TxUtil interface {
	Info(addrFrom string) (to []*AddressAmount, err error)
	FromAddresses() (froms []string)
	Sign(privKeys []string) (err error)
	TxForSend() (v interface{}, err error)
	IsSign() bool
	TxId() string
	Marshal() ([]byte, error)
	UnMarshal(data []byte) error
	Local() bool
}

type AddressAmount struct {
	Address string
	Amount  *big.Int
}
