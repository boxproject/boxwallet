package bckey_test

import (
	"testing"

	"reflect"

	"crypto/ecdsa"

	"github.com/boxproject/boxwallet/bccore"
	"github.com/boxproject/boxwallet/bckey"
	"github.com/boxproject/boxwallet/bckey/distribute"
	"github.com/boxproject/boxwallet/mock"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/crypto"
)

type AddressHandler struct {
	bckey.AddressFomatter
	TypeName string
}

func (ch *AddressHandler) LoadService(g bckey.AddressFomatter) error {
	if g != nil {
		ch.AddressFomatter = g
	}
	typ := reflect.TypeOf(g)
	ch.TypeName = typ.String()
	return nil
}

var (
	btc     = &distribute.BtcAddress{bccore.TestNet}
	btcs    = "*distribute.BtcAddress"
	eth     = &distribute.EthAddress{}
	eths    = "*distribute.EthAddress"
	ltc     = &distribute.LtcAddress{bccore.MainNet}
	ltcs    = "*distribute.LtcAddress"
	handler AddressHandler
)

func TestKey_Address(t *testing.T) {
	var (
		pub     = mock.Pubkey4
		address string
		err     error
	)
	handler.LoadService(btc)
	switch handler.TypeName {
	case btcs:
		address, err = btc.Address(pub)
		if err == nil {
			if address != mock.Addr {
				t.Fail()
				t.Error("transfer err")
			}
		}
	case eths:
		address, err = eth.Address(pub)
		if err == nil {
			if address != "0x6f3b6e51477DBccCF88EcC64c78F9C92eA4039c7" {
				t.Fail()
				t.Error("transfer err")
			}
		}
	case ltcs:
		address, err = ltc.Address(pub)
		if err == nil {
			if address != "0x6f3b6e51477DBccCF88EcC64c78F9C92eA4039c7" {
				t.Fail()
				t.Error("transfer err")
			}
		}
	}
	if err != nil {
		t.Fail()
		t.Error(err)
	} else {
		t.Log(address)
	}
}
func BenchmarkKey_Address(b *testing.B) {
	//key, _ := hdkeychain.NewKeyFromString(mock.Pubkey2)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		//keyConvert(key)
		eth.Address(mock.Pubkey2)
	}
}
func keyConvert(key *hdkeychain.ExtendedKey) (string, error) {
	pub, err := key.ECPubKey()
	if err != nil {
		return "", err
	}
	ethPub := &ecdsa.PublicKey{
		Curve: pub.Curve,
		X:     pub.X,
		Y:     pub.Y,
	}
	return crypto.PubkeyToAddress(*ethPub).String(), nil
}
func BenchmarkChildkey(b *testing.B) {
	key, _ := hdkeychain.NewKeyFromString(mock.Pubkey2)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		childkey(key)
	}
}

func childkey(key *hdkeychain.ExtendedKey) (*hdkeychain.ExtendedKey, error) {
	return key.Child(20)
}
func TestGetAddress(t *testing.T) {
	arr := bckey.GetAddress(bccore.BC_ETH)
	for _, v := range *arr {
		t.Log(v)
	}
}
