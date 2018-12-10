package pipeline

import (
	"sync"

	"github.com/boxproject/boxwallet/bccore"
)

type Heights struct {
	PubHeight uint64
	CurHeight uint64
}

type HeightCache struct {
	m sync.Map
}

var heightCache *HeightCache

func GetHeightCacheInstance() *HeightCache {
	return heightCache
}
func init() {
	heightCache = &HeightCache{}
	heightCache.m.Store(bccore.BC_BTC, &Heights{})
	heightCache.m.Store(bccore.BC_ETH, &Heights{})
	heightCache.m.Store(bccore.BC_LTC, &Heights{})
}

func (h *HeightCache) Push(chainType bccore.BloclChainType, pubH, curH uint64) {
	h.m.Store(chainType, &Heights{pubH, curH})
}

func (h *HeightCache) LoadAll() map[bccore.BloclChainType]*Heights {
	m := make(map[bccore.BloclChainType]*Heights)
	m1, ok := h.m.Load(bccore.BC_BTC)
	if ok {
		m[bccore.BC_BTC] = m1.(*Heights)
	}
	m2, ok := h.m.Load(bccore.BC_ETH)
	if ok {
		m[bccore.BC_ETH] = m2.(*Heights)
	}
	m3, ok := h.m.Load(bccore.BC_LTC)
	if ok {
		m[bccore.BC_LTC] = m3.(*Heights)
	}
	return m
}
