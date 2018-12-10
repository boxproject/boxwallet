package mock

import (
	"os"
	"path"

	"github.com/boxproject/boxwallet/bccoin"
	"github.com/boxproject/boxwallet/bcconfig"
	"github.com/boxproject/boxwallet/bcconfig/keynet"
	mysqlConf "github.com/boxproject/boxwallet/bcconfig/mysql"
	offcnf "github.com/boxproject/boxwallet/bcconfig/official"
	"github.com/boxproject/boxwallet/bccore"
	"github.com/boxproject/boxwallet/bckey"
	"github.com/boxproject/boxwallet/bctrans/client"
	"github.com/boxproject/boxwallet/bctrans/token"
	"github.com/boxproject/boxwallet/daemon"
	"github.com/boxproject/boxwallet/db"
	mysqlDB "github.com/boxproject/boxwallet/db/mysql"
	"github.com/boxproject/boxwallet/official"
)

func init() {
	initBadger()
	initOfficial()
	daemon.InitLogPath(path.Join(Gopath, ProjectDir, ConfigDir, "blockheight"))
	Cache = initCoinCache()
	KeyUtil = initKeyUtil()

	//////////BTC
	path1 := path.Join(Gopath, ProjectDir, ConfigDir, "btc.yml")
	provide1, err := bcconfig.FromConfigString(path1, YmlExt)
	if err != nil {
		panic(err)
	}
	Btc = client.NewBtcClient(provide1)
	//////////usdt
	path3 := path.Join(Gopath, ProjectDir, ConfigDir, "usdt.yml")
	provide3, err := bcconfig.FromConfigString(path3, YmlExt)
	Usdt = client.NewUsdtClient(provide3)
	//////////ETH
	path2 := path.Join(Gopath, ProjectDir, ConfigDir, "eth.yml")
	provide2, err := bcconfig.FromConfigString(path2, YmlExt)
	if err != nil {
		panic(err)
	}
	//////////ERC20
	Eth = client.NewEthClient(provide2)

	Erc20 = client.NewErc20Client(provide2)
	Erc20Token, _ = token.GetErc20TokenInstance()
	/////////LTC
	path4 := path.Join(Gopath, ProjectDir, ConfigDir, "ltc.yml")
	provide4, err := bcconfig.FromConfigString(path4, YmlExt)
	Ltc = client.NewLtcClient(provide4)

	go daemon.Start(DaemonCnf)
}

var (
	Cache         *bccoin.CoinCache
	KeyUtil       *bckey.KeyUtil
	Gopath        = os.Getenv("GOPATH")
	ProjectDir    = "src/github.com/boxproject/boxwallet"
	MySqlConfPath = "bcconfig/config.yml"
	YmlExt        = "yml"

	ConfigDir = "bcconfig"

	///client
	Btc   *client.BtcClient
	Usdt  *client.UsdtClient
	Ltc   *client.LtcClient
	Eth   *client.EthClient
	Erc20 *client.Erc20Client

	Erc20Token *token.Erc20Token
	//DB
	DaemonCnf  = initDaemonCnf()
	TxStorage  = initTxStorage()
	PropertyId = bccore.Token("2147483651")
)

var (
	Seed    = "我是谁"
	Prvkey  = "tprv8ZgxMBicQKsPdcnN6QTU491ex39He27ZCz6MXWKTuukfgZDDkwWa7KMuDs9KCy32R7im9xRvG7WAXy9E4cHvnBkCq1kFHQqGmHFnEYT7X8s"
	Pubkey  = "tpubD6NzVbkrYhZ4X5p9z484TYfmX4fDoMJTnHh8p2MmLBZ4X3TzPLLAHoymPy3UvpHLgqDeNeUBdefVN8C1MuoWDYnPGDE8tsvja92sioSG4na"
	Addr    = "n1kNSSzaxc6kdDkmKdBWmL7TKHQDKmm1R2"
	AddrEth = "0x95C8F9D46E50F19A1eFCA20C9bDe264eaf68E6E0"

	Seed2    = "太白菜"
	Prvkey2  = "tprv8ZgxMBicQKsPehCkD1GFhCkTTUyLNQmUg48DUD8z3XTECm7haQmNhLUurghBT4xatBLhzkeRVFWGBrDqq4fHo8kknxB7BCnqb7o6wwMfpbR"
	Pubkey2  = "tpubD6NzVbkrYhZ4YAEY6evr6cQa2WVGXjxPFMizkjBHToFd3FNUCoaxsq6n2oa6gHBpTCv9B6NygMdDZrQkTmAMg7Z94vAKKX2f11HtLoUqh2b"
	Addr2    = "mopBiZmA4iXvEEzw7feScBSHdVnWakyFC5"
	AddrEth2 = "0x430458dCfd8e88f37f1E60507EC6e0E84aa30118"

	Seed3    = "呵呵呵"
	Prvkey3  = "tprv8ZgxMBicQKsPd1CNcn35ZWEAJxMo9h9osH3WTykNv7pxNg6DDuRX3iK5Q2qpfPxwMHv2gVq2wmBQxfYrtQPZnLU9qHGppc64ckdRG6KHTT1"
	Pubkey3  = "tpubD6NzVbkrYhZ4WUEAWRhfxutGsysjK2LiSaeHkVngLPdMDALyrJF7ECvwa8u52kmkQfnP9CsyjF5famnU6MDGZeYHK7UL8jpziPERkTtfCbE"
	Addr3    = "miDdEcGWejYnweerVAvbDYPPWQ9rSFD4aD"
	AddrEth3 = "0xb8A88a13172d84d647c8460817Bfbb049ca8C24c"

	Seed4   = "abc"
	Prvkey4 = "tprv8ZgxMBicQKsPdys7eGCVDiQuFt1kzK5xR7TiL3KKBn9PFEpgHPKXdCA9EhE3jFDMCa5tSyJqFx1Ybsp23pZKnGiYK2jtyWqyy21ggLDJtbC"
	Pubkey4 = "tpubD6NzVbkrYhZ4XStuXus5d851puXh9eGrzR4VcZMcc3wn5j5Sun97ogn1QntFfHHyQ3qzcbgMCbTfoByupV3Ve8wohJCCZJ27ntLqziDePkD"
	Addr4   = "n1fNQ3XRrJqdsMFYe9xkYx7cT6mQWEH6c2"
)

func initBadger() {
	path := path.Join(Gopath, ProjectDir, "/db/data")
	if db.GetStore() != nil {
		return
	}
	err := db.Open(path)
	if err != nil {
		panic(err)
	}
}

func initCoinCache() *bccoin.CoinCache {
	path1 := path.Join(Gopath, ProjectDir, "bccoin/distribute/coin_info.json")
	return bccoin.InitCoinCache(db.GetStore(), db.Coin_Info, path1)
}
func initKeyUtil() *bckey.KeyUtil {
	path1 := path.Join(Gopath, ProjectDir, MySqlConfPath)
	provide, err := bcconfig.FromConfigString(path1, YmlExt)
	if err != nil {
		panic(err)
	}
	conf, err := keynet.DecodeConfig(provide)
	if conf.Net != 1 && conf.Net != 2 {
		conf.Net = 2
	}
	return bckey.InitKeyUtil(db.GetStore(), db.Pfk_PubKey, db.Pfk_Pubkey_Count, conf.Net)
}
func initTxStorage() *mysqlDB.TxStorage {
	path1 := path.Join(Gopath, ProjectDir, MySqlConfPath)
	provide, err := bcconfig.FromConfigString(path1, YmlExt)
	if err != nil {
		panic(err)
	}
	conf, err := mysqlConf.DecodeConfig(provide)
	if err != nil {
		panic(err)
	}
	return mysqlDB.NewTxStorage(conf.KeyStorage)
}

func initOfficial() {
	path1 := path.Join(Gopath, ProjectDir, ConfigDir, "official/btc.yml")
	btcp, err := bcconfig.FromConfigString(path1, YmlExt)
	if err != nil {
		panic(err)
	}
	official.InitBtcNode(btcp)
	path2 := path.Join(Gopath, ProjectDir, ConfigDir, "official/ltc.yml")
	ltcp, err := bcconfig.FromConfigString(path2, YmlExt)
	if err != nil {
		panic(err)
	}
	official.InitLtcNode(ltcp)
	path3 := path.Join(Gopath, ProjectDir, ConfigDir, "official/eth.yml")
	Ethp, err := bcconfig.FromConfigString(path3, YmlExt)
	if err != nil {
		panic(err)
	}
	official.InitEthNode(Ethp)
	path4 := path.Join(Gopath, ProjectDir, ConfigDir, "official/usdt.yml")
	usdtp, err := bcconfig.FromConfigString(path4, YmlExt)
	if err != nil {
		panic(err)
	}
	official.InitUsdtNode(usdtp)
	official.InitOfficDaemons(DaemonCnf)
}
func initDaemonCnf() offcnf.Config {
	path1 := path.Join(Gopath, ProjectDir, ConfigDir, "daemon_cnf.yml")
	provide, err := bcconfig.FromConfigString(path1, "yml")
	if err != nil {
		panic(err)
	}
	conf, err := offcnf.DecodeConfig(provide)
	if err != nil {
		panic(err)
	}
	return conf

}
