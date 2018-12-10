package clientseries

import (
	"math/big"

	"strconv"

	"github.com/boxproject/boxwallet/bcconfig"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthSeriesClient struct {
	C           *ethclient.Client
	DefGasPrice *big.Int
	DefGasLimit uint64
}

func NewEthSeriesClient(cfg bcconfig.Provider) *EthSeriesClient {
	url := cfg.GetString("url")
	gas := cfg.GetString("gasPrice")
	limit := cfg.GetString("gasLimit")
	dg := big.NewInt(0)
	dg, flag := dg.SetString(gas, 10)
	if !flag {
		panic("set gasPrice error")
	}
	gasLimit, err := strconv.ParseUint(limit, 10, 64)
	if err != nil {
		panic("ser gasLimit error")
	}
	c, err := ethclient.Dial(url)
	if err != nil {
		panic(err)
	}
	intance := &EthSeriesClient{
		C:           c,
		DefGasPrice: dg,
		DefGasLimit: gasLimit,
	}
	return intance
}
