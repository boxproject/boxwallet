package daemon_test

import (
	"context"
	"testing"

	"math/big"

	"strings"

	"encoding/hex"

	"fmt"

	"time"

	"github.com/boxproject/boxwallet/daemon"
	"github.com/boxproject/boxwallet/util"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var ethDaemon *daemon.EthDaemon

func init() {
	ethDaemon = daemon.NewEthDaemon(uint64(6), 6, time.Second)

}

func TestEthDaemon_GetBlockHeight(t *testing.T) {
	height, h2, err := ethDaemon.GetBlockHeight()
	if err != nil {
		t.Fail()
		t.Error(err)
	} else {
		t.Log(height, h2)
	}
	sync, err := ethDaemon.C.SyncProgress(context.Background())
	if sync != nil {
		//delayBlocks := sync.HighestBlock - sync.CurrentBlock
		t.Log(sync.HighestBlock, sync.CurrentBlock)
	}
}

func TestEthDaemon_GetBlockInfo(t *testing.T) {
	si, err := ethDaemon.GetBlockInfo(big.NewInt(1080))
	if err != nil {
		t.Fail()
		t.Error(err)
		return
	}
	txIds := si.AnalyzeBlock([]string{"mhAfGecTPa9eZaaNkGJcV7fmUPFi3T2Ki8"})
	t.Log(txIds)
}

func TestEthDaemon_GetTransaction(t *testing.T) {
	//0xe042795078c374a6c9ae929974757602d72564d8f360855b39560291ac49f876
	txId := &daemon.TxIdInfo{
		TxId:   "0xe8a44553f6fa84b63d9a82631f93b48f7345a7589fc587c33be402a3e5cd3fb9",
		BlockH: big.NewInt(1080),
	}
	/*txId = &daemon.TxIdInfo{
		"0xf758f559aa78703b8fa012f6547da4a088e4f784bb805be0e2e61aadf47b418e",
		big.NewInt(945),
	}*/
	/////
	txHash := common.HexToHash(txId.TxId)
	txRe, err := ethDaemon.C.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		t.Error(err)
		//return
	} else {

		from := common.BytesToAddress(txRe.Logs[0].Topics[1].Bytes()[:32]).Hex()
		to := common.BytesToAddress(txRe.Logs[0].Topics[2].Bytes()[:32]).Hex()

		// amount
		amount := common.BytesToHash(txRe.Logs[0].Data[:32]).Big()
		t.Log(from, to, amount, txRe.Logs[0].Address.String())
	}
	t.Log(txRe)
	/*for _, v := range txRe.Logs {
		t.Log(v)
	}
	*/
	txi, err := ethDaemon.GetTransaction(txId, false)
	if err != nil {
		t.Fail()
		t.Error(err)
		return
	}

	for _, v := range txi.In {
		t.Log("in:", v.Amt.String(), v.Addr)
	}
	for _, v := range txi.Out {
		t.Log("out:", v.Amt.String(), v.Addr)
	}
	t.Log(txi.Cnfm)

}
func TestUnPack(t *testing.T) {

	tokenAbi, err := abi.JSON(strings.NewReader(util.TokenABI))
	if err != nil {
		t.Error(err)
		return
	}
	var transferEvent struct {
		From  common.Address
		To    common.Address
		Value *big.Int
	}
	//"0xa9059cbb
	encodedData := "a9059cbb00000000000000000000000095c8f9d46e50f19a1efca20c9bde264eaf68e6e000000000000000000000000000000000000000000000000140980e4a1898e000"

	fmt.Println(encodedData[:8], encodedData[8:8+64], encodedData[8+64:])
	str := string("0x" + encodedData[32:72])
	if strings.ToLower("0x95C8F9D46E50F19A1eFCA20C9bDe264eaf68E6E0") == strings.ToLower(str) {
		t.Log("OK")
	}
	//   a9059cbb 00000000000000000000000067fa2c06c9c6d4332f330e14a66bdf1873ef3d2b 0000000000000000000000000000000000000000000000000de0b6b3a7640000
	//0x a9059cbb 000000000000000000000000ee7b38d5e730e31f9376416be9d0efa88f426667 00000000000000000000000000000000000000000000000000000004a817c800
	decodeData, err := hex.DecodeString(encodedData)
	if err != nil {
		t.Error(err)
		return
	}
	err = tokenAbi.Unpack(&transferEvent, "Transfer", decodeData)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(transferEvent.To.Hex(), transferEvent.Value.String())
}

func TestGetTxByHash(t *testing.T) {
	txHash := common.HexToHash("0xe8a44553f6fa84b63d9a82631f93b48f7345a7589fc587c33be402a3e5cd3fb9")
	txInfo, _, err := ethDaemon.C.TransactionByHash(context.Background(), txHash)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	gasused := big.NewInt(0).SetUint64(txInfo.Gas())
	t.Log(txInfo.Value())
	t.Log(txInfo.Cost().String())
	t.Log(txInfo.GasPrice())
	t.Log(txInfo.Gas())
	gasused.Mul(txInfo.GasPrice(), gasused)
	t.Log("1aaaaaa:", gasused.String())
	signer := types.NewEIP155Signer(txInfo.ChainId())
	msg, err := txInfo.AsMessage(signer)
	fee := big.NewInt(0)
	fee = fee.Sub(txInfo.Cost(), msg.Value())
	t.Log(msg.Value())
	t.Log(fee.String())
	txRe, err := ethDaemon.C.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		t.Error(err)
		//return
	}

	gu2 := big.NewInt(0).SetUint64(txRe.CumulativeGasUsed)
	gu2.Mul(gu2, txInfo.GasPrice())
	t.Log(gu2)
	t.Log(txRe.GasUsed)
}
