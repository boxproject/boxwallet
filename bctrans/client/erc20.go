package client

import (
	"context"

	"strings"

	"math/big"

	"github.com/boxproject/boxwallet/bccoin"
	"github.com/boxproject/boxwallet/bcconfig"
	"github.com/boxproject/boxwallet/bccore"
	"github.com/boxproject/boxwallet/bckey"
	"github.com/boxproject/boxwallet/bctrans/clientseries"
	"github.com/boxproject/boxwallet/errors"
	"github.com/boxproject/boxwallet/official"
	"github.com/boxproject/boxwallet/signature"
	"github.com/boxproject/boxwallet/util"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var erc20CliIntance *Erc20Client

type Erc20Client struct {
	*clientseries.EthSeriesClient
	*bckey.KeyUtil
	abi.ABI
	nmap map[bool]*clientseries.EthSeriesClient
	gap  int
}

func NewErc20Client(cfg bcconfig.Provider) *Erc20Client {
	if erc20CliIntance != nil {
		return erc20CliIntance
	}
	erc20CliIntance = new(Erc20Client)
	erc20CliIntance.EthSeriesClient = clientseries.NewEthSeriesClient(cfg)
	erc20CliIntance.KeyUtil = bckey.GetKeyUtilInstance()
	abi, err := abi.JSON(strings.NewReader(util.TokenABI))
	if err != nil {
		panic(err)
	}
	erc20CliIntance.nmap = make(map[bool]*clientseries.EthSeriesClient)
	erc20CliIntance.nmap[true] = ethCliIntance.EthSeriesClient
	erc20CliIntance.nmap[false] = official.GetEthNode()
	erc20CliIntance.gap = cfg.GetInt("gap")
	erc20CliIntance.ABI = abi
	return erc20CliIntance
}
func GetErc20ClientIntance() (*Erc20Client, error) {
	if erc20CliIntance == nil {
		return nil, errors.ERR_NIL_REFERENCE
	} else {
		return erc20CliIntance, nil
	}
}
func (c *Erc20Client) ImportAddress(address string, local bool) error {
	return nil
}
func (c *Erc20Client) GetNewAddress() (bckey.GenericKey, error) {
	return c.KeyUtil.GeneraterKey(nil, bccore.BC_ERC20)
}
func (c *Erc20Client) GetBalance(address string, token string, local bool) (balance bccoin.CoinAmounter, err error) {

	tc, err := util.NewTokenCaller(common.HexToAddress(token), c.nmap[local].C)
	if err != nil {
		return nil, err
	}
	addr := common.HexToAddress(address)
	bal, err := tc.BalanceOf(nil, addr)
	if err != nil {
		return nil, err
	}
	balance, _ = bccoin.NewCoinAmountFromBigInt(bccore.BC_ERC20, bccore.Token(token), bal)
	return
}

func (c *Erc20Client) CreateTx(addrsFrom []string, token string, addrsTo []*bccoin.AddressAmount, feeCeo float64) (txu signature.TxUtil, err error) {
	if len(addrsTo) != 1 || len(addrsFrom) != 1 {
		return nil, errors.ERR_PARAM_NOT_VALID
	}
	local, err := c.ChooseClientNode()
	if err != nil {
		return nil, err
	}
	ctx := context.TODO()
	contract := common.HexToAddress(token)
	if err != nil {
		return nil, err
	}
	to := addrsTo[0]
	toAddress := common.HexToAddress(to.Address)
	fromAccDef := accounts.Account{
		Address: common.HexToAddress(addrsFrom[0]),
	}
	tbalance, err := c.GetBalance(addrsFrom[0], token, local)
	if err != nil {
		return nil, err
	}
	cmp, err := tbalance.Cmp(to.Amount)
	if err != nil {
		return nil, err
	}
	if cmp < 0 {
		return nil, errors.ERR_NOT_ENOUGH_COIN
	}
	gasprice, err := c.nmap[local].C.SuggestGasPrice(ctx)
	if err != nil {
		gasprice = c.DefGasPrice
	}
	//contract/abi
	nonce, _ := c.nmap[local].C.NonceAt(ctx, fromAccDef.Address, nil)
	input, err := c.ABI.Pack("transfer", toAddress, to.Amount.Val())
	if err != nil {
		return nil, err
	}
	msg := ethereum.CallMsg{From: fromAccDef.Address, To: &contract, Value: big.NewInt(0), Data: input}
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

	tx := types.NewTransaction(nonce, contract, big.NewInt(0), gasLimit, gasprice, input)
	cost := tx.Cost()
	balance, _ := c.nmap[local].C.BalanceAt(context.TODO(), fromAccDef.Address, nil)
	if balance.Cmp(cost) < 0 {
		return nil, errors.ERR_NOT_ENOUGH_COIN
	}
	return signature.NewEthTx(nonce, fromAccDef.Address, contract, big.NewInt(0), gasLimit, gasprice, input, local), nil
}

func (c *Erc20Client) SendTx(txu signature.TxUtil) error {
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

func (c *Erc20Client) ChooseClientNode() (local bool, err error) {
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
