package bccore

type BloclChainType uint32

//main chain type
const (
	BC_DEF BloclChainType = 0 //exception
	BC_BTC BloclChainType = 1
	BC_ETH BloclChainType = 2
	BC_LTC BloclChainType = 3
)

//Tokens，Use the main chain address, no additional generation of new addresses
const (
	BC_ERC20 = BC_ETH
	BC_USDT  = BC_BTC
)

type Token string

type BlockChainSign string

const (
	STR_BTC   BlockChainSign = "BTC"
	STR_ETH   BlockChainSign = "ETH"
	STR_ERC20 BlockChainSign = "ERC20"
	STR_USDT  BlockChainSign = "USDT"
	STR_LTC   BlockChainSign = "LTC"
)

var (
	BCMap = map[BlockChainSign]BloclChainType{
		STR_BTC:   BC_BTC,
		STR_ETH:   BC_ETH,
		STR_ERC20: BC_ERC20,
		STR_USDT:  BC_USDT,
		STR_LTC:   BC_LTC,
	}
)

type Net int8

const (
	MainNet Net = 1
	TestNet Net = 2
)
