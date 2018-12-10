/*
As a communication channel between client and daemon
Only two modules can be referenced one way each
To control the address，Facilitate monitoring

## concept
1. To block the address
2. A single address can only be serialized
3. Multiple addresses can be parallel


## Design ideas
1. when createTx ,uuid is a unique identifier，cache addresses first,
2. when sendTx  uuid,txid,addresses(get from cache)，And add a txid primary key data
3. after scan block ,txid pk =>uuid pk =>address pk
*/
package pipeline

import (
	"time"

	"sync"

	"log"

	"github.com/boxproject/boxwallet/bccore"
	"github.com/boxproject/boxwallet/errors"
)

func init() {
	go poolGC()
}

type addrLockKey string
type addrLockBody struct {
	createAt  time.Time
	sendCheck bool      //set true after send sccuess
	deadline  time.Time //clear after deadline
}

type txLockKey string //uuid
type txLockBody struct {
	ch        chan bool
	used      bool
	addresses []string
}

type txScanKey string //txId
type txScanBody struct {
	uuid string
	used bool
}

func getAddrLockKey(chainType bccore.BloclChainType, address string) addrLockKey {
	return addrLockKey(string(chainType) + address)
}
func getTxLockKey(chainType bccore.BloclChainType, uuid string) txLockKey {
	return txLockKey(string(chainType) + uuid)
}
func getScanKey(chainType bccore.BloclChainType, txId string) txScanKey {
	return txScanKey(string(chainType) + txId)
}

type lockPool struct {
	m   *sync.Mutex
	kv  map[addrLockKey]*addrLockBody
	uti map[txLockKey]*txLockBody
	ti  map[txScanKey]*txScanBody
}

var p = &lockPool{
	m:   &sync.Mutex{},
	kv:  make(map[addrLockKey]*addrLockBody, 1000),
	uti: make(map[txLockKey]*txLockBody, 1000),
	ti:  make(map[txScanKey]*txScanBody, 1000),
}

type Pipeline struct{}

func (*Pipeline) AddressExist(chainType bccore.BloclChainType, addresses []string) bool {
	p.m.Lock()
	defer p.m.Unlock()
	{
		for _, v := range addresses {
			k := getAddrLockKey(chainType, v)
			if p.kv[k] != nil {
				return true
			}
		}
		return false
	}
}

func (*Pipeline) Send(chainType bccore.BloclChainType, uuid string, addresses []string) bool {
	p.m.Lock()
	defer p.m.Unlock()
	{
		keys := make([]addrLockKey, 0, len(addresses))
		for _, v := range addresses {
			key := getAddrLockKey(chainType, v)
			if p.kv[key] != nil {
				return false
			}
			keys = append(keys, key)
		}
		tb := &addrLockBody{
			createAt:  time.Now(),
			sendCheck: false,
			deadline:  time.Now().AddDate(0, 0, 5),
		}
		for _, k := range keys {
			ktmp := k
			p.kv[ktmp] = tb
		}
		tik := getTxLockKey(chainType, uuid)
		p.uti[tik] = &txLockBody{
			addresses: addresses,
		}
		return true
	}
}

func (*Pipeline) CheckSend(chainType bccore.BloclChainType, txId, uuid string) (ch <-chan bool, err error) {
	p.m.Lock()
	defer p.m.Unlock()
	{
		utik := getTxLockKey(chainType, uuid)
		if p.uti[utik] == nil {
			return nil, errors.ERR_PIPELINE_DATA_ILLEGAL
		}
		ch := make(chan bool)
		p.uti[utik].ch = ch
		for _, v := range p.uti[utik].addresses {
			key := getAddrLockKey(chainType, v)
			if p.kv[key] != nil {
				p.kv[key].sendCheck = true
			}
		}
		tik := getScanKey(chainType, txId)
		p.ti[tik] = &txScanBody{
			uuid: uuid,
		}
		return ch, nil
	}
}
func (*Pipeline) TxOver(chainType bccore.BloclChainType, txId string, result bool) {
	p.m.Lock()
	defer p.m.Unlock()
	{
		//txid find uuid
		tik := getScanKey(chainType, txId)
		if p.ti[tik] == nil || p.ti[tik].used {
			return
		} else {
			p.ti[tik].used = true
		}
		utik := getTxLockKey(chainType, p.ti[tik].uuid)
		if p.uti[utik] != nil {
			if !p.uti[utik].used {
				p.uti[utik].ch <- result
				close(p.uti[utik].ch)
				p.uti[utik].used = true
			}
			//uuid find addresses
			for _, v := range p.uti[utik].addresses {
				addr := v
				k := getAddrLockKey(chainType, addr)
				delete(p.kv, k)
			}
		}
	}
}

func poolGC() {
	tick := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-tick.C:
			p.m.Lock()
			{
				for k, v := range p.ti {
					if v.used {
						log.Println(k)
						delete(p.ti, k)
					}
				}
				for k, v := range p.uti {
					if v.used {
						log.Println(k)
						delete(p.uti, k)
					}
				}
				for k, v := range p.kv {
					if !v.sendCheck && v.deadline.After(time.Now()) {
						log.Println(k, time.Now().Format("Mon Jan _2 15:04:05 2006"))
						delete(p.kv, k)
					}
				}
			}
			p.m.Unlock()
		}
	}
}
