package client

import (
	"context"

	"math/big"

	"github.com/boxproject/boxwallet/bccoin"
	"github.com/boxproject/boxwallet/bcconfig"
	"github.com/boxproject/boxwallet/bccore"
	"github.com/boxproject/boxwallet/bckey"
	"github.com/boxproject/boxwallet/bctrans/clientseries"
	"github.com/boxproject/boxwallet/errors"
	"github.com/boxproject/boxwallet/official"
	"github.com/boxproject/boxwallet/signature"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var ethCliIntance *EthClient

type EthClient struct {
	*clientseries.EthSeriesClient
	*bckey.KeyUtil
	nmap map[bool]*clientseries.EthSeriesClient
	gap  int
}

func NewEthClient(cfg bcconfig.Provider) *EthClient {
	if ethCliIntance != nil {
		return ethCliIntance
	}
	ethCliIntance = new(EthClient)
	ethCliIntance.EthSeriesClient = clientseries.NewEthSeriesClient(cfg)
	ethCliIntance.KeyUtil = bckey.GetKeyUtilInstance()
	ethCliIntance.nmap = make(map[bool]*clientseries.EthSeriesClient)
	ethCliIntance.nmap[true] = ethCliIntance.EthSeriesClient
	ethCliIntance.nmap[false] = official.GetEthNode()
	ethCliIntance.gap = cfg.GetInt("gap")
	return ethCliIntance
}
func GetEthClientIntance() (*EthClient, error) {
	if ethCliIntance == nil {
		return nil, errors.ERR_NIL_REFERENCE
	} else {
		return ethCliIntance, nil
	}
}

func (c *EthClient) GetBalance(address string, token string, local bool) (balance bccoin.CoinAmounter, err error) {
	amount, err := c.nmap[local].C.BalanceAt(context.TODO(), common.HexToAddress(address), nil)
	if err != nil {
		return nil, err
	}
	balance, _ = bccoin.NewCoinAmountFromBigInt(bccore.BC_ETH, "", amount)
	return
}
func (c *EthClient) ImportAddress(address string, local bool) error {
	return nil
}
func (c *EthClient) GetNewAddress() (bckey.GenericKey, error) {
	key, err := c.KeyUtil.GeneraterKey(nil, bccore.BC_ETH)
	return key, err
}

func (c *EthClient) CreateTx(addrsFrom []string, token string, addrsTo []*bccoin.AddressAmount, feeCeo float64) (txu signature.TxUtil, err error) {
	if len(addrsTo) != 1 || len(addrsFrom) != 1 {
		return nil, errors.ERR_PARAM_NOT_VALID
	}
	local, err := c.ChooseClientNode()
	if err != nil {
		return nil, err
	}
	ctx := context.TODO()
	fromAccDef := accounts.Account{
		Address: common.HexToAddress(addrsFrom[0]),
	}
	to := addrsTo[0]
	toAccDef := accounts.Account{
		Address: common.HexToAddress(to.Address),
	}
	nonce, _ := c.nmap[local].C.NonceAt(ctx, fromAccDef.Address, nil)
	gasprice := big.NewInt(0)

	gasprice, err = c.nmap[local].C.SuggestGasPrice(ctx)
	if err != nil {
		gasprice = c.DefGasPrice
	}

	msg := ethereum.CallMsg{From: fromAccDef.Address, To: nil, Value: big.NewInt(0), Data: nil}
	gasLimit, err := c.nmap[local].C.EstimateGas(ctx, msg)
	if err != nil {
		gasLimit = c.DefGasLimit
	}

	excessCeo := (feeCeo - 1) * 100
	if excessCeo != 0 {
		prec, _ := big.NewFloat(excessCeo).Int64()
		excess := big.NewInt(0).Quo(gasprice, big.NewInt(prec))
		gasprice.Add(gasprice, excess)
	}

	// nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte
	tx := types.NewTransaction(
		nonce,
		toAccDef.Address,
		to.Amount.Val(),
		gasLimit,
		gasprice,
		nil,
	)
	balance, _ := c.nmap[local].C.BalanceAt(context.TODO(), fromAccDef.Address, nil)

	cost := tx.Cost()
	if balance.Cmp(cost) < 0 {
		return nil, errors.ERR_NOT_ENOUGH_COIN
	}
	return signature.NewEthTx(nonce, fromAccDef.Address, toAccDef.Address, to.Amount.Val(), gasLimit, gasprice, nil, local), nil

}

func (c *EthClient) SendTx(txu signature.TxUtil) error {
	txv, err := txu.TxForSend()
	if err != nil {
		return nil
	}
	tx := txv.(*types.Transaction)
	err = c.nmap[txu.Local()].C.SendTransaction(context.TODO(), tx)
	if err != nil {
		return err
	}
	return nil
}

func (c *EthClient) ChooseClientNode() (local bool, err error) {
	sync, err := c.C.SyncProgress(context.Background())
	if err != nil {
		return true, err
	}
	if sync == nil {
		return true, nil
	} else {
		if sync.HighestBlock-sync.CurrentBlock > uint64(c.gap) {
			return false, nil
		} else {
			return true, nil
		}
	}
}
