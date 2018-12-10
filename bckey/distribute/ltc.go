package distribute

import (
	"github.com/boxproject/boxwallet/bccore"
	"github.com/ltcsuite/ltcd/chaincfg"
	"github.com/ltcsuite/ltcutil/hdkeychain"
)

type LtcAddress struct {
	Net bccore.Net
}

func (a *LtcAddress) Address(pubkey string) (address string, err error) {
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
