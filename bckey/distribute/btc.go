package distribute

import (
	"github.com/btcsuite/btcd/chaincfg"

	"github.com/boxproject/boxwallet/bccore"
	"github.com/btcsuite/btcutil/hdkeychain"
)

type BtcAddress struct {
	Net bccore.Net
}

func (a *BtcAddress) Address(pubkey string) (address string, err error) {
	key, err := hdkeychain.NewKeyFromString(pubkey)
	if err != nil {
		return
	}
	curNet := &chaincfg.RegressionNetParams
	if a.Net == bccore.MainNet {
		curNet = &chaincfg.MainNetParams
	}

	pub, err := key.Address(curNet)
	if err != nil {
		return
	}
	address = pub.EncodeAddress()
	return
}
