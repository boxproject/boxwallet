package bckey_test

import (
	"log"
	"testing"

	"crypto/ecdsa"

	"github.com/btcsuite/btcd/chaincfg"

	"github.com/boxproject/boxwallet/bckey"
	"github.com/boxproject/boxwallet/mock"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/crypto"
	lctchaincfg "github.com/ltcsuite/ltcd/chaincfg"
	ltchdkeychain "github.com/ltcsuite/ltcutil/hdkeychain"
	_ "github.com/stretchr/testify/mock"
)

func TestGenerateSimpleKey(t *testing.T) {
	key, err := bckey.GenerateSimpleKey([]byte("abc"))
	if err != nil {
		t.Fail()
		return
	}
	log.Println(key.String())
	if key.String() != mock.Prvkey {
		t.Fail()
	}

}

func TestGenerateMasterPubKey(t *testing.T) {
	prv, _ := hdkeychain.NewKeyFromString("tprv8ZgxMBicQKsPdys7eGCVDiQuFt1kzK5xR7TiL3KKBn9PFEpgHPKXdCA9EhE3jFDMCa5tSyJqFx1Ybsp23pZKnGiYK2jtyWqyy21ggLDJtbC")
	key, err := bckey.GenerateMasterPubKey(prv)
	if err != nil {
		t.Fail()
		return
	}
	log.Println(key.String())
}

func TestBase58(t *testing.T) {
	key1, _ := hdkeychain.NewKeyFromString("tpubD6NzVbkrYhZ4XStuXus5d851puXh9eGrzR4VcZMcc3wn5j5Sun97ogn1QntFfHHyQ3qzcbgMCbTfoByupV3Ve8wohJCCZJ27ntLqziDePkD")

	pub, _ := key1.ECPubKey()
	ethPub := &ecdsa.PublicKey{
		Curve: pub.Curve,
		X:     pub.X,
		Y:     pub.Y,
	}

	address := crypto.PubkeyToAddress(*ethPub).String()
	t.Log(address)
}

func TestBtcCmpLtc(t *testing.T) {
	prv, _ := ltchdkeychain.NewKeyFromString(mock.Prvkey)
	prv2, _ := hdkeychain.NewKeyFromString(mock.Prvkey)
	key, err := prv.Neuter()
	key2, _ := prv2.Neuter()
	if err != nil {
		t.Fail()
		return
	}
	log.Println("---------------------------")
	log.Println("ltc:", key.String())
	log.Println("btc:", key2.String())
	log.Println("---------------------------")
	addr1, _ := key.Address(&lctchaincfg.MainNetParams)
	addr2, _ := key2.Address(&chaincfg.MainNetParams)
	log.Println("ltc:", addr1)
	log.Println("btc:", addr2)
}
