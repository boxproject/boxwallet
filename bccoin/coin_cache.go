package bccoin

import (
	"bytes"
	"encoding/json"
	"strconv"

	"github.com/boxproject/boxwallet/bccore"
	"github.com/boxproject/boxwallet/errors"
	"github.com/boxproject/boxwallet/util"
)

var defualtCoinCache *CoinCache

// prefix+coinType+token
type CoinCache struct {
	db     util.Database
	prefix []byte
}

func InitCoinCache(db util.Database, prefix []byte, path string) *CoinCache {
	if defualtCoinCache == nil {
		defualtCoinCache = &CoinCache{db: db, prefix: prefix}
		defualtCoinCache.BatchSaveCoinInfos(loadCoinInfo(path))
	}
	return defualtCoinCache
}

func GetCoinCacheIntance() (*CoinCache, error) {
	if defualtCoinCache == nil {
		return nil, errors.ERR_NIL_REFERENCE
	} else {
		return defualtCoinCache, nil
	}
}

func (c *CoinCache) BatchSaveCoinInfos(cis []*CoinInfo) error {

	kvp := make([]*util.Pair, len(cis), len(cis))
	for k, v := range cis {
		buff := &bytes.Buffer{}
		buff.Write(c.prefix)
		buff.Write(v.Sign())
		kvp[k] = &util.Pair{}
		kvp[k].Key = buff.Bytes()
		kvp[k].Val = v.Bytes()
	}
	return c.db.BatchPut(kvp)
}

func (c *CoinCache) SaveCoinInfo(ci *CoinInfo) error {
	buff := &bytes.Buffer{}
	buff.Write(c.prefix)
	buff.Write(ci.Sign())
	k := buff.Bytes()
	v := ci.Bytes()
	return c.db.Put(k, v)
}

func (c *CoinCache) GetCoinInfo(bc bccore.BloclChainType, t bccore.Token) (*CoinInfo, error) {
	buff := &bytes.Buffer{}
	buff.Write(c.prefix)
	buff.WriteString(strconv.FormatInt(int64(bc), 10))
	if t != "" {
		buff.WriteString("_")
		buff.WriteString(string(t))
	}
	v, err := c.db.Get(buff.Bytes())
	if err != nil {
		return nil, err
	}
	ci := &CoinInfo{}
	err = json.Unmarshal(v, ci)
	return ci, err
}

func (c *CoinCache) GetAll(chainType bccore.BloclChainType, hasToken bool) (<-chan *CoinInfo, error) {
	data := make(chan *CoinInfo, 127)
	ch, err := c.db.Iterator(c.prefix)
	if err != nil {
		return nil, err
	}
	for c := range ch {
		ci := &CoinInfo{}
		err = json.Unmarshal(c.Val, ci)
		if err == nil && ci.CT == chainType && (hasToken == (ci.Token != "")) {
			data <- ci
		}
	}
	close(data)
	return data, nil
}
