package daemon

import (
	"math/big"
	"time"

	"github.com/boxproject/boxwallet/bccoin"
	"github.com/boxproject/boxwallet/bccore"
	"github.com/boxproject/boxwallet/bctrans/client"
	"github.com/boxproject/boxwallet/errors"
	"github.com/ltcsuite/ltcd/chaincfg"
	"github.com/ltcsuite/ltcd/chaincfg/chainhash"
	"github.com/ltcsuite/ltcd/txscript"
	"github.com/ltcsuite/ltcd/wire"
)

type LtcDaemon struct {
	*client.LtcClient
	TickSecond time.Duration
	ConfirmNum uint64
	UnlockNum  uint64
}

var (
	ltcBG *LtcDaemon
)

func NewLtcDaemon(confirmNum, unlocknNum uint64, tickSecond time.Duration) *LtcDaemon {
	lctCli, err := client.GetLtcClientIntance()

	if err != nil {
		panic(err)
	}
	bd := &LtcDaemon{
		TickSecond: tickSecond,
		ConfirmNum: confirmNum,
		LtcClient:  lctCli,
		UnlockNum:  unlocknNum,
	}
	return bd
}

func (d *LtcDaemon) GetLoopDuration() time.Duration {
	return d.TickSecond
}

func (d *LtcDaemon) GetBlockHeight() (*big.Int, uint64, error) {
	bc, err := d.C.GetBlockChainInfo()
	if err != nil {
		return nil, 0, err
	}
	return big.NewInt(int64(bc.Blocks)), uint64(bc.Headers), nil
}

//爬块
func (d *LtcDaemon) GetBlockInfo(blockHeight *big.Int) (v BlockAnalyzer, err error) {
	hash, err := d.C.GetBlockHash(blockHeight.Int64())
	if err != nil {
		return nil, err
	}
	blockInfo, err := d.C.GetBlock(hash)
	if err != nil {
		return nil, err
	}
	return &ltcBlock{blockInfo, d.Env, blockHeight}, nil
}
func (d *LtcDaemon) CheckConfirmations(confirm uint64) bool {
	return confirm >= d.ConfirmNum
}
func (d *LtcDaemon) CheckUnlock(confirm uint64) bool {
	return confirm >= d.UnlockNum
}

//分析 tx详情
func (d *LtcDaemon) GetTransaction(txId *TxIdInfo, easy bool) (*TxInfo, error) {
	txHash, err := chainhash.NewHashFromStr(txId.TxId)
	if err != nil {
		return nil, err
	}
	txInfo, err := d.C.GetRawTransactionVerbose(txHash)
	if err != nil {
		return nil, err
	}
	txi := &TxInfo{}
	txi.TxId = txId.TxId
	txi.Cnfm = txInfo.Confirmations
	txi.Target = d.ConfirmNum
	txi.UnlockNum = d.UnlockNum
	txi.ExtValid = true
	txi.BCT = bccore.BC_LTC
	txi.BCS = bccore.STR_LTC
	if !easy {

		txi.H = big.NewInt(0).Set(txId.BlockH)
		inSum := big.NewInt(0)
		for _, v := range txInfo.Vin {
			aaIn, err := d.GetFromAddr(v.Txid, v.Vout)
			if err != nil {
				return nil, err
			}
			inSum = inSum.Add(inSum, aaIn.Amt)
			txi.In = append(txi.In, aaIn)
		}
		outSum := big.NewInt(0)
		for _, v := range txInfo.Vout {
			if v.ScriptPubKey.Type == pubkeyhash {
				am, _ := bccoin.NewCoinAmountFromFloat(bccore.BC_LTC, "", v.Value)
				out := &AddrAmount{
					Addr: v.ScriptPubKey.Addresses[0],
					Amt:  am.Val(),
				}
				outSum = outSum.Add(outSum, am.Val())
				txi.Out = append(txi.Out, out)
			}
		}
		txi.Fee = inSum.Sub(inSum, outSum)
	}
	return txi, nil
}

func (d *LtcDaemon) GetFromAddr(txId string, vout uint32) (*AddrAmount, error) {
	txHash, err := chainhash.NewHashFromStr(txId)
	if err != nil {
		return nil, err
	}
	tx, err := d.C.GetRawTransactionVerbose(txHash)
	if err != nil {
		return nil, err
	}
	if len(tx.Vout) >= int(vout) {
		aa := &AddrAmount{}
		am1, _ := bccoin.NewCoinAmountFromFloat(bccore.BC_LTC, "", tx.Vout[vout].Value)
		aa.Amt = am1.Val()
		aa.Addr = tx.Vout[vout].ScriptPubKey.Addresses[0]
		return aa, nil
	} else {
		return nil, errors.ERR_TX_OUT_INDEX_OVERFLEW
	}
}

///////////////////////////////////////////////////////////////////////////

type ltcBlock struct {
	*wire.MsgBlock
	*chaincfg.Params
	H *big.Int
}

func (b *ltcBlock) AnalyzeBlock(addresses []string) (txIds []*TxIdInfo) {
	txIds = make([]*TxIdInfo, 0, 2)
	flag := false
	for _, v := range b.Transactions {
		flag = false
		for _, v1 := range v.TxOut {
			if txscript.IsUnspendable(v1.PkScript) {
				continue
			}
			_, addres, _, err := txscript.ExtractPkScriptAddrs(v1.PkScript, b.Params)
			if err != nil {
				return nil
			}
			if len(addres) == 0 {
				continue
			}
			addr := addres[0].EncodeAddress()
			for _, curAddr := range addresses {
				if curAddr == addr {
					txifo := &TxIdInfo{TxId: v.TxHash().String(), BlockH: b.H}
					txIds = append(txIds, txifo)
					flag = true
					break
				}
			}
			if flag {
				break
			}
		}
	}
	return txIds
}
