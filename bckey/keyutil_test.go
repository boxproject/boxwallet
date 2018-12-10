package bckey_test

import (
	"testing"

	"github.com/boxproject/boxwallet/bccore"
	"github.com/boxproject/boxwallet/mock"
	"github.com/btcsuite/btcutil/hdkeychain"
)

var (
	keyUtil = mock.KeyUtil
)

func TestKeyUtil_SaveMasterKey(t *testing.T) {
	err := keyUtil.SaveMasterKey(mock.Pubkey)
	if err != nil {
		t.Fail()
		t.Error(err)
	}
}
func TestKeyUtil_GetMasterKey(t *testing.T) {
	key, err := keyUtil.GetMasterKey()
	if err != nil {
		t.Fail()
		t.Error(err)
	} else {
		t.Log(key)
	}
}
func TestKeyUtil_GeneraterKey(t *testing.T) {
	key, err := keyUtil.GeneraterKey(nil, bccore.BC_LTC)
	if err != nil {
		t.Fail()
		t.Error(err)
	}
	t.Logf("success, chain type：%d\r\n customdeep:%d\r\n curnum:%d\r\n", key.KeyType(), key.CustomDeep(), key.CurrentNum())

	pub, _ := hdkeychain.NewKeyFromString(mock.Pubkey)
	keyErr, _ := pub.Child(key.CurrentNum())
	t.Log(keyErr.String())
	key1, _ := pub.Child(uint32(bccore.BC_LTC))
	t.Log(key1.String())
	key2, _ := key1.Child(key.CurrentNum())
	t.Log(key2.String())
	/*for i := 0; i < 10; i++ {
		key, err = keyUtil.GeneraterKey([]uint32{1}, bccore.BC_BTC)
		if err != nil {
			t.Fail()
			t.Error(err)
		}
		t.Logf("success, chain type：%d\r\n customdeep:%d\r\n curnum:%d\r\n", key.KeyType(), key.CustomDeep(), key.CurrentNum())
	}
	for i := 0; i < 10; i++ {
		key, err = keyUtil.GeneraterKey([]uint32{1, 1}, bccore.BC_BTC)
		if err != nil {
			t.Fail()
			t.Error(err)
		}
		t.Logf("success, chain type：%d\r\n customdeep:%d\r\n curnum:%d\r\n", key.KeyType(), key.CustomDeep(), key.CurrentNum())
	}*/
}

func TestKeyUtil_GetDeepCount(t *testing.T) {
	inner := []struct {
		kt   bccore.BloclChainType
		deep []uint32
	}{
		{
			kt:   bccore.BC_BTC,
			deep: nil,
		},
		{
			kt:   bccore.BC_BTC,
			deep: []uint32{1},
		},
		{
			kt:   bccore.BC_BTC,
			deep: []uint32{1, 1},
		},
	}
	for _, v := range inner {
		count, err := keyUtil.GetDeepCount(v.kt, v.deep)
		if err != nil {
			t.Fail()
			t.Error(err)
		}
		t.Log(v.kt, v.deep, count)
	}
}
func TestKeyUtil_GetKey(t *testing.T) {
	key, err := keyUtil.GetKey(nil, bccore.BC_BTC, 32, true)
	if err != nil {
		t.Fail()
		t.Error(err)
	}
	t.Log(key.Key().String())
}

func TestKeyUtil_GetChildKeyCount(t *testing.T) {
	t.Log(keyUtil.GetChildKeyCount())
}
