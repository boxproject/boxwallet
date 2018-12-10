package bckey

import (
	"fmt"
	"strconv"
	"sync"

	"log"

	"github.com/boxproject/boxwallet/bccore"
	"github.com/boxproject/boxwallet/errors"
	"github.com/boxproject/boxwallet/util"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/dgraph-io/badger"
)

//Both public and private keys are available
type KeyUtil struct {
	db            util.Database
	Pfk_Key       []byte //秘钥存储前缀
	Pfk_key_Count []byte //秘钥组计数存储前缀
	net           bccore.Net
}

var (
	mu      = &sync.Mutex{}
	keyUtil *KeyUtil
)

// InitKeyUtil is used to initialize the key management tool
// pfk_key and pfk_key_count are the prefix of storage for the kvdb
func InitKeyUtil(db util.Database, pfk_key, pfk_key_count []byte, net bccore.Net) *KeyUtil {
	if keyUtil != nil {
		return keyUtil
	}
	keyUtil = &KeyUtil{db: db, Pfk_Key: pfk_key, Pfk_key_Count: pfk_key_count, net: net}
	initAddressMemCache(net)
	log.Println("KeyUtil init success")
	return keyUtil
}

func GetKeyUtilInstance() *KeyUtil {
	return keyUtil
}

//GeneraterKey generate self-increment key at the specified location
//Ps: now  => bc/custom/9
//    next => bc/custom/10
func (ku *KeyUtil) GeneraterKey(customDeep []uint32, bc bccore.BloclChainType) (key GenericKey, err error) {
	mu.Lock()
	defer mu.Unlock()
	{
		count, err := ku.GetDeepCount(bc, customDeep)
		if err != nil {
			return nil, err
		}
		key = newKey(customDeep, bc, 0)

		hdkey := &hdkeychain.ExtendedKey{}
		lastKey := emptyKey()

		if count > 0 {
			goto GENTLE
		} else {
			goto RUDE

		}
	GENTLE:
		lastKey, err = ku.GetKey(customDeep, bc, 0, false)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				goto RUDE
			} else {
				return nil, err
			}
		}
		count++
		key.SetCurrentNum(count)
		hdkey, err = lastKey.Key().Child(count)
		if err != nil {
			goto RUDE
		} else {
			fmt.Println(hdkey.String())
			err = ku.SaveChildKey(bc, customDeep, count, hdkey.String())
			if err != nil {
				return nil, err
			}
			key.SetKey(hdkey)
			addAddress(bc, key.Address())
			return key, nil
		}
	RUDE:
		masterKey, err := ku.GetMasterKey()
		if err != nil {
			return nil, err
		}
		fmt.Println(masterKey.Address(&chaincfg.RegressionNetParams))
		hdkey, err = masterKey.Child(uint32(bc))
		if err != nil {
			return nil, err
		}
		fmt.Println(hdkey.Address(&chaincfg.RegressionNetParams))
		for _, v := range customDeep {
			hdkey, err = hdkey.Child(v)
			if err != nil {
				return nil, err
			}
		}
		count++
		key.SetCurrentNum(count)
		hdkey, err = hdkey.Child(count)
		if err != nil {
			return nil, err
		}
		fmt.Println(hdkey.Address(&chaincfg.RegressionNetParams))
		err = ku.SaveChildKey(bc, customDeep, count, hdkey.String())
		if err != nil {
			return nil, err
		}

		key.SetKey(hdkey)
		addAddress(bc, key.Address())
		return key, nil
	}
}

//SaveMasterKey must save the key first,because all the public keys need him to spawn them
func (ku *KeyUtil) SaveMasterKey(key string) error {
	err := ku.db.Put([]byte(ku.Pfk_Key), []byte(key))
	if err != nil {
		return err
	}
	//Special handling
	initMasterAddress(bccore.BC_BTC, keyUtil.net)
	initMasterAddress(bccore.BC_ETH, keyUtil.net)
	initMasterAddress(bccore.BC_LTC, keyUtil.net)
	return nil
}

func (ku *KeyUtil) GetMasterKey() (key *hdkeychain.ExtendedKey, err error) {
	v, err := ku.db.Get([]byte(ku.Pfk_Key))
	if err != nil {
		return
	}
	key, err = hdkeychain.NewKeyFromString(string(v))
	return
}

func (ku *KeyUtil) SaveChildKey(bc bccore.BloclChainType, parentDeep []uint32, curNum uint32, key string) (err error) {
	ky := newKey(parentDeep, bc, curNum)
	sk, err := ky.getStorageKey(ku.Pfk_Key, true)
	if err != nil {
		return
	}
	totalCount, err := ku.GetChildKeyCount()
	if err != nil {
		return err
	}
	if totalCount >= 10000 {
		return errors.ERR_KEY_OVERFLOW
	}
	err = ku.db.Put([]byte(sk), []byte(key))
	if err != nil {
		return
	}

	count, err := ku.GetDeepCount(bc, parentDeep)
	if err != nil {
		return
	}
	count++
	err = ku.SaveDeepCount(bc, parentDeep, count)
	if err != nil {
		return
	}
	return
}

//GetMasterGenericKey get masterkey and fomat to specifies the type of address
func (ku *KeyUtil) GetMasterGenericKey(bc bccore.BloclChainType) (k GenericKey, err error) {
	mk, err := ku.GetMasterKey()
	if err != nil {
		return nil, err
	}
	mkey := &key{
		key:    mk,
		master: true,
		kt:     bc,
	}
	return mkey, nil
}

// GetKey get deep key by params，but cannot be used to get the master key
// withCurNum:  false=>get parent, true=>get current
func (ku *KeyUtil) GetKey(customDeep []uint32, bc bccore.BloclChainType, curNum uint32, withCurNum bool) (k GenericKey, err error) {
	if customDeep == nil && !withCurNum {
		mk, err := ku.GetMasterKey()
		if err != nil {
			return nil, err
		}
		father, err := mk.Child(uint32(bc))
		if err != nil {
			return nil, err
		}
		mkey := &key{
			key:    father,
			master: true,
			kt:     bc,
		}
		return mkey, nil
	}
	ky := newKey(customDeep, bc, curNum)
	sk, err := ky.getStorageKey(ku.Pfk_Key, withCurNum)
	if err != nil {
		return
	}
	v, err := ku.db.Get(sk)
	if err != nil {
		return
	}
	keyStr := string(v)
	hdkey, err := hdkeychain.NewKeyFromString(keyStr)
	if err != nil {
		return
	}
	if withCurNum {
		ky.SetKey(hdkey)
		ky.SetCurrentNum(curNum)
		k = ky
	} else {
		l := len(customDeep)
		k = newKey(customDeep[:l-2], bc, 0)
		k.SetCurrentNum(customDeep[l-2:][0])
		k.SetKey(hdkey)
	}
	return
}

//Gets the total number of public keys at the specified depth
func (ku *KeyUtil) GetDeepCount(bc bccore.BloclChainType, parentDeep []uint32) (count uint32, err error) {
	k := newKey(parentDeep, bc, 0)
	sk, err := k.getStorageKey(ku.Pfk_key_Count, false)
	if err != nil {
		return
	}
	v_count, err := ku.db.Get(sk)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return 0, nil
		} else {
			return
		}
	}
	count64, err := strconv.ParseUint(string(v_count), 10, 32)
	count = uint32(count64)
	return
}

func (ku *KeyUtil) SaveDeepCount(bc bccore.BloclChainType, parentDeep []uint32, count uint32) (err error) {
	k := newKey(parentDeep, bc, 0)
	sk, err := k.getStorageKey(ku.Pfk_key_Count, false)
	if err != nil {
		return
	}
	err = ku.db.Put(sk, []byte(strconv.FormatInt(int64(count), 10)))
	return
}

func (ku *KeyUtil) GetChildKeyCount() (int, error) {
	ch, err := ku.db.Iterator(ku.Pfk_key_Count)
	if err != nil {
		return 0, err
	}
	sum := 0
	for c := range ch {
		count, _ := strconv.ParseUint(string(c.Val), 10, 32)
		sum += int(count)
	}
	return sum, nil
}
