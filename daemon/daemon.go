package daemon

import (
	"math/big"

	"reflect"

	"time"

	"log"

	"github.com/boxproject/boxwallet/bcconfig/official"
	"github.com/boxproject/boxwallet/bccore"
	"github.com/boxproject/boxwallet/bckey"
	"github.com/boxproject/boxwallet/db/mysql"
	"github.com/boxproject/boxwallet/pipeline"
	"github.com/boxproject/boxwallet/util"
)

var (
	btcAddressSilce map[bccore.BloclChainType][]string
	pipe            = &pipeline.Pipeline{}
	heightCache     = pipeline.GetHeightCacheInstance()
)

type Daemoner interface {
	//Scan block cycle
	GetLoopDuration() time.Duration
	GetBlockHeight() (*big.Int, uint64, error)
	GetBlockInfo(blockHeight *big.Int) (v BlockAnalyzer, err error)
	CheckConfirmations(confirm uint64) bool
	CheckUnlock(confirm uint64) bool
	GetTransaction(txId *TxIdInfo, easy bool) (*TxInfo, error)
}

type TxIdInfo struct {
	TxId   string
	BlockH *big.Int
}

type BlockAnalyzer interface {
	AnalyzeBlock(addresses []string) (txIds []*TxIdInfo)
}

var BlockTypeMap map[bccore.BloclChainType]reflect.Type

type AddrAmount struct {
	Addr string
	Amt  *big.Int
}

type TxInfo struct {
	TxId      string
	BCS       bccore.BlockChainSign
	BCT       bccore.BloclChainType
	Target    uint64        //The target counts as many blocks as successful
	T         bccore.Token  //contract address
	In        []*AddrAmount //for main chain
	Out       []*AddrAmount //fro main chain
	Cnfm      uint64        //Confirmation number
	UnlockNum uint64        //the tx count fro unlock pipeline
	H         *big.Int      //block height
	Time      time.Time
	Fee       *big.Int
	Token     string

	InExt  []*AddrAmount //for tokens
	OutExt []*AddrAmount //fro tokens
	// for tokens check
	ExtValid bool
}

var (
	daemons map[bccore.BloclChainType]Daemoner
	//btc 10min ;address count 10004;
	//600 seconds
	w         = util.NewTimingWheel(100*time.Millisecond, 6000)
	txStorage *mysql.TxStorage
)

//////////////////////////////////////////////Main process/////////////////////////////////////////////

func Start(config official.Config) {
	daemons := initDaemons(config)
	txStorage = mysql.GetTxStorageInstance()
	for k, v := range daemons {
		go func(key bccore.BloclChainType, value Daemoner) {
			daemonGo(value, key)
		}(k, v)
	}
	log.Println("daemon Start-up success")
	select {}
}

func initDaemons(config official.Config) map[bccore.BloclChainType]Daemoner {
	m := make(map[bccore.BloclChainType]Daemoner)
	m[bccore.BC_BTC] = NewBtcDaemon(config.Btc.Confirm, config.Btc.Unlock, time.Duration(config.Btc.Ticker)*time.Second)
	m[bccore.BC_ETH] = NewEthDaemon(config.Eth.Confirm, config.Eth.Unlock, time.Duration(config.Eth.Ticker)*time.Second)
	m[bccore.BC_LTC] = NewLtcDaemon(config.Ltc.Confirm, config.Ltc.Unlock, time.Duration(config.Ltc.Ticker)*time.Second)
	return m
}

//layer-1
func daemonGo(d Daemoner, chainType bccore.BloclChainType) {
	tick := d.GetLoopDuration()
	lastBlockHeight := HeightRead(chainType)
	log.Println("chainType:", chainType, ",init first height:", lastBlockHeight.Uint64())
	for {
		select {
		case <-w.After(tick):
			curIndex, pubHeight, err := d.GetBlockHeight()
			if err != nil {
				continue
			}
			heightCache.Push(chainType, pubHeight, curIndex.Uint64())
			if curIndex.Cmp(lastBlockHeight) == 0 {
				continue
			}
			if lastBlockHeight.Cmp(big.NewInt(0)) == 0 {
				lastBlockHeight.Set(curIndex)
				continue
			}
			start := big.NewInt(0).Set(lastBlockHeight)
			end := big.NewInt(0).Set(curIndex)
			log.Println("📏chainType:", chainType, ",block height:", start.String(), "-", end.String())

			/*		///DDD
					start = big.NewInt(1529848)
					end = big.NewInt(1529849)
			*/
			go blockListen(d, start, end, chainType)
			log.Println("📏chainType:", chainType, "The last time the block height is updated：", lastBlockHeight.String())
			HeightWriteAsyn(chainType, lastBlockHeight)
			lastBlockHeight.Set(curIndex)
			log.Println("📏chainType:", chainType, "After the latest block height update：", lastBlockHeight.String())
		}
	}
}

//layer-2
func blockListen(d Daemoner, start, end *big.Int, chainType bccore.BloclChainType) {
	log.Println("🔍📒chainType:", chainType, ",blocklisten:start ", start.Uint64(), "--end ", end.Uint64())
	for i := big.NewInt(0).Set(start); i.Cmp(end) < 0; i.Add(i, big.NewInt(1)) {
		bi, err := d.GetBlockInfo(big.NewInt(0).Set(i))
		if err != nil {
			continue
		}
		log.Println("🔍📒chainType:", chainType, ",blocklisten:", i.Uint64())
		txInfos := bi.AnalyzeBlock(*bckey.GetAddress(chainType))
		go txListen(txInfos, d, chainType)

	}
}

//layer-3
func txListen(t []*TxIdInfo, d Daemoner, chainType bccore.BloclChainType) {
	for _, txId := range t {
		log.Println("🔍📋chainType:", chainType, ",block heigth：", txId.BlockH.String(), ";txId:", txId.TxId)
		go func(t *TxIdInfo) {
			txPolling(t, d, chainType, false)
		}(txId)
	}
}

//layer-4
func txPolling(t *TxIdInfo, d Daemoner, chainType bccore.BloclChainType, withoutPipeline bool) {
	tick := d.GetLoopDuration()
	for {
		select {
		case <-w.After(tick):
			tx := txStorage.GetTx(t.TxId)
			easy := false
			if tx != nil {
				easy = tx.ExtValid
			}
			txInfo, err := d.GetTransaction(t, easy)
			if err != nil {
				log.Println("💿chainType:", chainType, err.Error())
				continue
			}
			if tx == nil {
				//save
				if !txStorage.AddTx(ConvertTxInfo(txInfo)) {
					log.Println("💿chainType:", chainType, "save failed")
				} else {
					log.Println("💿chainType:", chainType, "save:", txInfo.TxId, ";confirm:", txInfo.Cnfm, ":height:", txInfo.H.String())
				}
			} else if txInfo.Cnfm > tx.Confirmations {
				//update
				txStorage.UpdateTx(tx.TxId, txInfo.Cnfm, txInfo.ExtValid)
				log.Println("💿chainType:", chainType, "update:", txInfo.TxId, ";confirm:", txInfo.Cnfm, ";extValid:", txInfo.ExtValid)
			} else {
				log.Println("💾chainType:", chainType, "update:", txInfo.TxId, "already exists")
			}
			if d.CheckUnlock(txInfo.Cnfm) {
				if !withoutPipeline {
					pipe.TxOver(chainType, t.TxId, true)
				}
			}
			if d.CheckConfirmations(txInfo.Cnfm) {
				return
			}
		}
	}
}
