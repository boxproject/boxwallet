package official

import (
	"context"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"time"

	"log"

	daemonCnf "github.com/boxproject/boxwallet/bcconfig/official"
	"github.com/boxproject/boxwallet/bccore"
	"github.com/boxproject/boxwallet/bctrans/clientseries"
	"github.com/boxproject/boxwallet/errors"
	"github.com/boxproject/boxwallet/pipeline"
	"github.com/ethereum/go-ethereum/common"
	ltchainhash "github.com/ltcsuite/ltcd/chaincfg/chainhash"
)

var (
	officDaemons map[bccore.BloclChainType]OfficDaemon
	pipe         = &pipeline.Pipeline{}
)

func InitOfficDaemons(config daemonCnf.Config) {
	officDaemons = make(map[bccore.BloclChainType]OfficDaemon)
	officDaemons[bccore.BC_BTC] = &BtcOfficDaemon{
		C:       btcNodeInstance,
		Confirm: uint64(config.Btc.Confirm),
		Ticker:  time.Duration(config.Btc.Ticker) * time.Second,
	}
	officDaemons[bccore.BC_ETH] = &EthOfficDaemon{
		C:       ethNodeInstance,
		Confirm: int(config.Eth.Confirm),
		Ticker:  time.Duration(config.Eth.Ticker) * time.Second,
	}
	officDaemons[bccore.BC_LTC] = &BtcOfficDaemon{
		C:       btcNodeInstance,
		Confirm: uint64(config.Ltc.Confirm),
		Ticker:  time.Duration(config.Ltc.Ticker) * time.Second,
	}
	log.Println("official daemon init success")
}

type BtcOfficDaemon struct {
	C       *clientseries.OmniSeriesClient
	Confirm uint64
	Ticker  time.Duration
}
type EthOfficDaemon struct {
	C       *clientseries.EthSeriesClient
	Confirm int
	Ticker  time.Duration
}

type LtcOfficDaemon struct {
	C       *clientseries.LtcSeriesClient
	Confirm int
	Ticker  time.Duration
}

type OfficDaemon interface {
	GetTxResult(txId string) error
}

func ListenTx(chainType bccore.BloclChainType, txId string) {
	go func(txId string) {
		log.Println("offic tx listening ...")
		err := officDaemons[chainType].GetTxResult(txId)
		if err != nil {
			log.Println("official err:", err)
		}
		log.Println("offic tx Listen to complete！")
		pipe.TxOver(chainType, txId, err == nil)
	}(txId)
}

func (d *BtcOfficDaemon) GetTxResult(txId string) error {
	timeDeadline := time.Now().Add(d.Ticker * 100)
	ticker := time.NewTicker(d.Ticker)
	for {
		select {
		case <-ticker.C:
			txHash, err := chainhash.NewHashFromStr(txId)
			if err != nil {
				return err
			}
			txInfo, err := btcNodeInstance.C.GetRawTransactionVerbose(txHash)
			if err != nil {
				if time.Now().After(timeDeadline) {
					return errors.ERR_TIME_OUT
				}
				continue
			}
			if txInfo.Confirmations > d.Confirm {
				return nil
			} else {
				continue
			}
		}
	}
}
func (d *LtcOfficDaemon) GetTxResult(txId string) error {
	timeDeadline := time.Now().AddDate(0, 0, 2)
	ticker := time.NewTicker(d.Ticker)
	for {
		select {
		case <-ticker.C:
			txHash, err := ltchainhash.NewHashFromStr(txId)
			if err != nil {
				return err
			}
			txInfo, err := ltcNodeInstance.C.GetRawTransactionVerbose(txHash)
			if err != nil {
				if time.Now().After(timeDeadline) {
					return errors.ERR_TIME_OUT
				}
				continue
			}
			if txInfo.Confirmations > uint64(d.Confirm) {
				return nil
			} else {
				continue
			}
		}
	}
}

func (d *EthOfficDaemon) GetTxResult(txId string) error {
	timeDeadline := time.Now().AddDate(0, 0, 2)
	ticker := time.NewTicker(d.Ticker)
	for {
		select {
		case <-ticker.C:
			txHash := common.HexToHash(txId)
			_, ispending, err := ethNodeInstance.C.TransactionByHash(context.Background(), txHash)
			if err != nil {
				continue
			}
			if time.Now().After(timeDeadline) {
				return errors.ERR_TIME_OUT
			}
			if ispending {
				continue
			} else {
				return nil
			}
		}
	}
}
