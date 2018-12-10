/*
1. 初始化部分已订好的token
2. 添加部分未初始化的token


*/
package token

import (
	"github.com/boxproject/boxwallet/bccoin"
	"github.com/boxproject/boxwallet/bccore"
	"github.com/boxproject/boxwallet/bctrans/client"
	"github.com/boxproject/boxwallet/util"

	"github.com/dgraph-io/badger"
	"github.com/ethereum/go-ethereum/common"
)

type Erc20Token struct {
	cache  *bccoin.CoinCache
	client *client.Erc20Client
}

var defErc20Token *Erc20Token

func GetErc20TokenInstance() (*Erc20Token, error) {
	if defErc20Token != nil {
		return defErc20Token, nil
	}
	cc, err := bccoin.GetCoinCacheIntance()
	if err != nil {
		return nil, err
	}
	eci, err := client.GetErc20ClientIntance()
	if err != nil {
		return nil, err
	}
	defErc20Token = &Erc20Token{
		cache:  cc,
		client: eci,
	}
	//先放这里，目前只有erc20
	initTokenMemCache()
	return defErc20Token, nil
}

func (t *Erc20Token) GetTokenInfo(contract bccore.Token) (*bccoin.CoinInfo, error) {
	ci, err := t.cache.GetCoinInfo(bccore.BC_ERC20, contract)
	if err != nil && err != badger.ErrKeyNotFound {
		return nil, err
	}
	if ci != nil {
		return ci, nil
	}
	tc, err := util.NewTokenCaller(common.HexToAddress(string(contract)), t.client.C)
	if err != nil {
		return nil, err
	}
	symbol, err := tc.Symbol(nil)
	if err != nil {
		return nil, err
	}
	decimals, err := tc.Decimals(nil)
	if err != nil {
		return nil, err
	}
	name, err := tc.Name(nil)
	if err != nil {
		return nil, err
	}
	ci = &bccoin.CoinInfo{
		CT:       bccore.BC_ERC20,
		Token:    contract,
		Symbol:   symbol,
		Decimals: int(decimals.Int64()),
		Name:     name,
	}
	t.cache.SaveCoinInfo(ci)
	addToken(bccore.BC_ERC20, string(ci.Token))
	return ci, nil
}
