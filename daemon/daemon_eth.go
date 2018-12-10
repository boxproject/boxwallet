package daemon

import (
	"math/big"
	"time"

	"context"

	"github.com/boxproject/boxwallet/bccore"
	"github.com/boxproject/boxwallet/bctrans/client"
	"github.com/boxproject/boxwallet/errors"
	"github.com/ethereum/go-ethereum/common"

	"encoding/hex"

	"strings"

	"github.com/ethereum/go-ethereum/core/types"
)

type EthDaemon struct {
	*client.EthClient
	TickSecond time.Duration
	ConfirmNum uint64
	UnlockNum  uint64
}

const tranferId = "a9059cbb"

func NewEthDaemon(confirmNum, unlocknNum uint64, tickSecond time.Duration) *EthDaemon {
	ethCli, err := client.GetEthClientIntance()
	if err != nil {
		panic(err)
	}
	bd := &EthDaemon{
		TickSecond: tickSecond,
		ConfirmNum: confirmNum,
		EthClient:  ethCli,
		UnlockNum:  unlocknNum,
	}
	return bd
}

func (d *EthDaemon) GetLoopDuration() time.Duration {
	return d.TickSecond
}

/*func (d *EthDaemon) GetPubBlockHeight() (uint64, error) {
	sync, err := d.C.SyncProgress(context.Background())
	if err != nil {
		return 0, err
	}
	if sync == nil {
		return 0, nil
	} else {
		return sync.HighestBlock, nil
	}
}*/
func (d *EthDaemon) GetBlockHeight() (*big.Int, uint64, error) {
	/*block, err := d.C.HeaderByNumber(context.Background(), nil)
	if err != nil {
		//TODO LOG
		return nil, err
	}
	return block.Number, nil*/
	sync, err := d.C.SyncProgress(context.Background())
	if err != nil {
		return nil, 0, err
	}
	if sync == nil {
		block, err := d.C.HeaderByNumber(context.Background(), nil)
		if err != nil {
			return nil, 0, err
		}
		return block.Number, block.Number.Uint64(), nil
	} else {
		return big.NewInt(int64(sync.CurrentBlock)), sync.HighestBlock, nil
	}
}

func (d *EthDaemon) GetBlockInfo(blockHeight *big.Int) (v BlockAnalyzer, err error) {
	block, err := d.C.BlockByNumber(context.Background(), blockHeight)
	if err != nil {
		return nil, err
	}
	return &EthBlock{Block: block, H: blockHeight}, nil
}

func (d *EthDaemon) CheckConfirmations(confirm uint64) bool {
	return confirm >= d.ConfirmNum
}
func (d *EthDaemon) CheckUnlock(confirm uint64) bool {
	return confirm >= d.UnlockNum
}

func (d *EthDaemon) GetTransaction(txId *TxIdInfo, easy bool) (*TxInfo, error) {
	txHash := common.HexToHash(txId.TxId)
	txInfo, ispending, err := d.C.TransactionByHash(context.Background(), txHash)
	if err != nil {
		return nil, err
	}
	if ispending {
		return nil, errors.ERR_TX_PENDING
	}

	block, err := d.C.BlockByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	confirm := block.Number().Sub(block.Number(), txId.BlockH)
	txi := &TxInfo{}
	txi.TxId = txId.TxId
	txi.H = txId.BlockH
	txi.Target = d.ConfirmNum
	txi.UnlockNum = d.UnlockNum
	txi.Cnfm = confirm.Uint64()
	txi.ExtValid = true

	if !easy {
		signer := types.NewEIP155Signer(txInfo.ChainId())
		msg, err := txInfo.AsMessage(signer)
		if err != nil {
			return nil, err
		}
		in := &AddrAmount{
			Addr: msg.From().String(),
		}
		txRe, errTxRe := d.C.TransactionReceipt(context.Background(), txHash)
		if errTxRe == nil && txRe != nil {
			fee := big.NewInt(0).SetUint64(txRe.GasUsed)
			fee = fee.Mul(fee, txInfo.GasPrice())
			txi.Fee = fee
		}
		if len(msg.Data()) == 0 {
			out := &AddrAmount{
				Addr: msg.To().String(),
				Amt:  msg.Value(),
			}

			txi.BCS = bccore.STR_ETH
			txi.BCT = bccore.BC_ETH
			txi.In = append(txi.In, in)
			txi.Out = append(txi.Out, out)

		} else {
			txi.ExtValid = false
			//MethodID: 0xa9059cbb
			//address: 64
			//amount: 64
			hexStr := hex.EncodeToString(msg.Data())
			if len(hexStr) >= 72 && hexStr[:8] == tranferId {
				to := strings.ToLower("0x" + hexStr[32:72])
				value := hexStr[72:]

				val := big.NewInt(0)
				val, bl := val.SetString(value, 16)
				if !bl {
					return nil, nil
				}
				out := &AddrAmount{
					Addr: to,
					Amt:  val,
				}
				txi.Token = msg.To().String()
				txi.BCT = bccore.BC_ERC20
				txi.BCS = bccore.STR_ERC20

				txi.InExt = append(txi.InExt, in)
				txi.OutExt = append(txi.OutExt, out)
				if err == nil && txRe != nil && len(txRe.Logs) == 1 && len(txRe.Logs[0].Topics) == 3 {
					rto := common.BytesToAddress(txRe.Logs[0].Topics[2].Bytes()[:32]).Hex()
					ramount := common.BytesToHash(txRe.Logs[0].Data[:32]).Big()
					raddress := txRe.Logs[0].Address.String()
					if strings.ToLower(rto) == to && ramount.Cmp(val) == 0 && txi.Token == raddress {
						//log.Println("success")
						txi.ExtValid = true
					}
				}
			} else {
				//other
				txi.BCT = bccore.BC_ERC20
				txi.In = append(txi.In, in)
			}

		}
	}
	return txi, nil
}

///////////////////////////////////////////////////////////////////////////

type EthBlock struct {
	*types.Block
	H *big.Int
}

func (b *EthBlock) AnalyzeBlock(addresses []string) (txIds []*TxIdInfo) {
	txIds = make([]*TxIdInfo, 0, 2)
	for _, v := range b.Block.Transactions() {
		signer := types.NewEIP155Signer(v.ChainId())
		msg, err := v.AsMessage(signer)
		if err != nil {
			continue
		}
		to := ""
		if len(msg.Data()) == 0 {
			if msg.To() != nil {
				to = strings.ToLower(msg.To().String())
			}
		} else {
			hexStr := hex.EncodeToString(msg.Data())
			if len(hexStr) >= 72 && hexStr[:8] == tranferId {
				to = strings.ToLower("0x" + hexStr[32:72])
			}
		}
		for _, curAddr := range addresses {
			if msg.From().String() == curAddr || to == strings.ToLower(curAddr) {
				txifo := &TxIdInfo{TxId: v.Hash().String(), BlockH: b.H}
				txIds = append(txIds, txifo)
				continue
			}
		}
	}
	return txIds
}
