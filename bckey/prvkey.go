package bckey

import (
	"crypto/rand"

	"github.com/btcsuite/btcd/chaincfg"

	"github.com/btcsuite/btcutil/hdkeychain"
)

func GenerateSimpleKey(existSeed []byte) (ext *hdkeychain.ExtendedKey, err error) {
	var (
		seed [32]byte
	)
	if existSeed != nil {
		copy(seed[:], existSeed)
	} else {
		if _, err = rand.Read(seed[:]); err != nil {
			return
		}
	}
	ext, err = hdkeychain.NewMaster(seed[:], &chaincfg.RegressionNetParams)
	return
}

func GenerateMasterPubKey(prvkey *hdkeychain.ExtendedKey) (ext *hdkeychain.ExtendedKey, err error) {
	ext, err = prvkey.Neuter()
	return
}
