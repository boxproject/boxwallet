package bckey

import (
	"github.com/boxproject/boxwallet/bccore"
	"github.com/boxproject/boxwallet/bckey/distribute"
)

type AddressFomatter interface {
	Address(pubkey string) (address string, err error)
}

//address factory
func AddressHandler(bc bccore.BloclChainType, net bccore.Net) AddressFomatter {
	switch bc {
	case bccore.BC_BTC:
		return &distribute.BtcAddress{Net: net}
	case bccore.BC_ETH:
		return &distribute.EthAddress{}
	case bccore.BC_LTC:
		return &distribute.LtcAddress{Net: net}
	}
	panic("this key not exists.")
}
