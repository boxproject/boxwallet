package bckey

import (
	"bytes"
	"strconv"

	"github.com/boxproject/boxwallet/bccore"
	"github.com/boxproject/boxwallet/db"
	"github.com/boxproject/boxwallet/errors"
	"github.com/btcsuite/btcutil/hdkeychain"
)

/*
m/kt/customDeep[0]/customDeep[1]/..../customDeep[n]/curNum
curNum:Internal control increment
todo 可以加一层parentKey
*/
type key struct {
	key        *hdkeychain.ExtendedKey
	customDeep []uint32
	curNum     uint32
	kt         bccore.BloclChainType
	master     bool
}

type GenericKey interface {
	//Gets the database prefix  px:prefix  withCurNum: get current or get parent
	getStorageKey(px db.PrefixKey, withCurNum bool) (key []byte, err error)

	CustomDeep() []uint32

	KeyType() bccore.BloclChainType

	CurrentNum() uint32
	//key info
	Key() *hdkeychain.ExtendedKey

	SetCurrentNum(num uint32)
	SetKey(key *hdkeychain.ExtendedKey)

	Address() string
	IsMaster() bool
}

func newKey(customDeep []uint32, kt bccore.BloclChainType, curNum uint32) GenericKey {
	return &key{
		customDeep: customDeep,
		kt:         kt,
		curNum:     curNum,
	}
}
func emptyKey() GenericKey {
	return &key{}
}

///////////////GenericKey interface start//////////////////

func (k *key) Key() *hdkeychain.ExtendedKey {
	return k.key
}
func (k *key) CustomDeep() []uint32 {
	return k.customDeep
}
func (k *key) KeyType() bccore.BloclChainType {
	return k.kt
}
func (k *key) CurrentNum() uint32 {
	return k.curNum
}
func (k *key) SetCurrentNum(num uint32) {
	k.curNum = num
}
func (k *key) SetKey(key *hdkeychain.ExtendedKey) {
	k.key = key
}
func (k *key) getStorageKey(px db.PrefixKey, withCurNum bool) (key []byte, err error) {
	if k.kt == bccore.BC_DEF {
		err = errors.ERR_NIL_REFERENCE
		return
	}
	buff := &bytes.Buffer{}
	buff.Write(px)
	buff.WriteString(strconv.FormatUint(uint64(k.kt), 10))
	for _, v := range k.customDeep {
		buff.WriteString("_")
		buff.WriteString(strconv.FormatUint(uint64(v), 10))
	}
	if withCurNum {
		buff.WriteString(strconv.FormatUint(uint64(k.curNum), 10))
	}
	key = buff.Bytes()
	return
}
func (k *key) IsMaster() bool {
	return k.master
}
func (k *key) Address() string {
	handler := AddressHandler(k.kt, keyUtil.net)
	address, _ := handler.Address(k.key.String())
	return address
}

///////////////GenericKey interface end//////////////////
