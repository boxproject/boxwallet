package token

import (
	"sync"

	"github.com/boxproject/boxwallet/bccoin"
	"github.com/boxproject/boxwallet/bccore"
)

//全局address Map
// *[]string => address pointer
//对于引用
var globalTokenPool map[bccore.BloclChainType]*TokenMemCache

type TokenMemCache struct {
	TokenPointer *[]string
	M            *sync.RWMutex
}

//todo 初始化所有地址，先这样用着，先从之前的key中直接同步到内存中
func init() {
	//go不能直接使用循环遍历"枚举"，因为这里的枚举其实是const常量
	globalTokenPool = make(map[bccore.BloclChainType]*TokenMemCache)
	globalTokenPool[bccore.BC_ERC20] = newTokenMemCache(bccore.BC_ERC20)
}
func initTokenMemCache() {
	initToken(bccore.BC_ERC20)
}
func newTokenMemCache(chainType bccore.BloclChainType) *TokenMemCache {
	amc := &TokenMemCache{
		TokenPointer: new([]string),
		M:            new(sync.RWMutex), //其实直接用mutex就够了
	}
	return amc
}
func GetToken(chainType bccore.BloclChainType) *[]string {
	if globalTokenPool[chainType] == nil {
		return &[]string{}
	}
	return globalTokenPool[chainType].TokenPointer
}
func addToken(chainType bccore.BloclChainType, address string) {
	globalTokenPool[chainType].M.Lock()
	defer globalTokenPool[chainType].M.Unlock()
	{
		*globalTokenPool[chainType].TokenPointer = append(*globalTokenPool[chainType].TokenPointer, address)
	}
}

func initToken(chainType bccore.BloclChainType) {
	cache, _ := bccoin.GetCoinCacheIntance()
	ch, err := cache.GetAll(chainType, true)
	if err != nil {
		panic(err)
	}
	for c := range ch {
		addToken(chainType, string(c.Token))
	}
}
