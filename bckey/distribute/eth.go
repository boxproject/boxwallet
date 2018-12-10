package distribute

import (
	"crypto/ecdsa"

	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/crypto"
)

type EthAddress struct {
}

func (a *EthAddress) Address(pubkey string) (address string, err error) {
	key, err := hdkeychain.NewKeyFromString(pubkey)
	if err != nil {
		return
	}
	pub, err := key.ECPubKey()
	if err != nil {
		return
	}
	ethPub := &ecdsa.PublicKey{
		Curve: pub.Curve,
		X:     pub.X,
		Y:     pub.Y,
	}
	address = crypto.PubkeyToAddress(*ethPub).String()
	return
}
