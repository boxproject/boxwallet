package official

import (
	"github.com/boxproject/boxwallet/bcconfig"
	"github.com/boxproject/boxwallet/bctrans/clientseries"
)

var btcNodeInstance *clientseries.OmniSeriesClient
var ethNodeInstance *clientseries.EthSeriesClient
var ltcNodeInstance *clientseries.LtcSeriesClient
var usdtNodeInstance *clientseries.OmniSeriesClient

func InitBtcNode(cfg bcconfig.Provider) *clientseries.OmniSeriesClient {
	if btcNodeInstance != nil {
		return btcNodeInstance
	}
	btcNodeInstance = clientseries.NewOmniSeriesClient(cfg)
	return btcNodeInstance
}
func InitLtcNode(cfg bcconfig.Provider) *clientseries.LtcSeriesClient {
	if ltcNodeInstance != nil {
		return ltcNodeInstance
	}
	ltcNodeInstance = clientseries.NewLtcSeriesClient(cfg)
	return ltcNodeInstance
}

func InitEthNode(cfg bcconfig.Provider) *clientseries.EthSeriesClient {
	if ethNodeInstance != nil {
		return ethNodeInstance
	}
	ethNodeInstance = clientseries.NewEthSeriesClient(cfg)
	return ethNodeInstance
}

func InitUsdtNode(cfg bcconfig.Provider) *clientseries.OmniSeriesClient {
	if usdtNodeInstance != nil {
		return usdtNodeInstance
	}
	usdtNodeInstance = clientseries.NewOmniSeriesClient(cfg)
	return usdtNodeInstance
}

func GetBtcNode() *clientseries.OmniSeriesClient {
	if btcNodeInstance == nil {
		panic("btc official node connect failed")
	}
	return btcNodeInstance
}
func GetLtcNode() *clientseries.LtcSeriesClient {
	if ltcNodeInstance == nil {
		panic("ltc official node connect failed")
	}
	return ltcNodeInstance
}
func GetEthNode() *clientseries.EthSeriesClient {
	if ethNodeInstance == nil {
		panic("eth official node connect failed")
	}
	return ethNodeInstance
}
func GetUsdtNode() *clientseries.OmniSeriesClient {
	if usdtNodeInstance == nil {
		panic("usdt official node connect failed")
	}
	return usdtNodeInstance
}
