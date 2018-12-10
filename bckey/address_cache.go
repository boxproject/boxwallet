package bckey

import (
	"sync"

	"github.com/boxproject/boxwallet/bccore"
)

// *[]string => address pointer
var globalAddressPool map[bccore.BloclChainType]*AddressMemCache

type AddressMemCache struct {
	AddrsPointer *[]string
	M            *sync.RWMutex
}

//Initialize existing addresses into memory
func init() {
	globalAddressPool = make(map[bccore.BloclChainType]*AddressMemCache)
	globalAddressPool[bccore.BC_BTC] = newAddressMemCache(bccore.BC_BTC)
	globalAddressPool[bccore.BC_ETH] = newAddressMemCache(bccore.BC_ETH)
	globalAddressPool[bccore.BC_LTC] = newAddressMemCache(bccore.BC_LTC)
	//log.Println("address cache init success")
	//*globalAddressPool[bccore.BC_LTC].AddrsPointer = append(*globalAddressPool[bccore.BC_LTC].AddrsPointer ,"LXtwa3rx9JExquao2q23MXcSfvEmEfnrFv")
}

func initAddressMemCache(net bccore.Net) {
	initAddresses(bccore.BC_BTC, net)
	initAddresses(bccore.BC_ETH, net)
	initAddresses(bccore.BC_LTC, net)
}
func newAddressMemCache(chainType bccore.BloclChainType) *AddressMemCache {
	amc := &AddressMemCache{
		AddrsPointer: new([]string),
		M:            new(sync.RWMutex),
	}
	return amc
}
func GetAddress(chainType bccore.BloclChainType) *[]string {
	if globalAddressPool[chainType] == nil {
		return &[]string{}
	}
	return globalAddressPool[chainType].AddrsPointer
}
func addAddress(chainType bccore.BloclChainType, address string) {
	globalAddressPool[chainType].M.Lock()
	defer globalAddressPool[chainType].M.Unlock()
	{
		*globalAddressPool[chainType].AddrsPointer = append(*globalAddressPool[chainType].AddrsPointer, address)
	}
}

func initAddresses(chainType bccore.BloclChainType, net bccore.Net) {
	count, err := keyUtil.GetDeepCount(chainType, nil)
	if err != nil {
		panic(err)
	}
	initMasterAddress(chainType, net)
	for i := uint32(1); i <= count; i++ {
		go func(num uint32) {
			key, err := keyUtil.GetKey(nil, chainType, num, true)
			if err == nil {
				//Persist in memory
				addAddress(chainType, key.Address())
			}
		}(i)
	}
}

func initMasterAddress(chainType bccore.BloclChainType, net bccore.Net) {
	mkey, _ := keyUtil.GetMasterKey()
	//If the main key is empty, there is no need to import it.
	if mkey == nil {
		return
	}
	addrFmt := AddressHandler(chainType, net)
	addr, _ := addrFmt.Address(mkey.String())
	addAddress(chainType, addr)
}
