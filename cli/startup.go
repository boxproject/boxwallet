package cli

import (
	"path"

	"os"

	"log"

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
	init_1_Path()
	initBadger()
	initOfficial()
	daemon.InitLogPath(path.Join(ROOT, WalletPath.HEIGHT))
	Cache = initCoinCache()

	KeyUtil = initKeyUtil()

	//////////BTC
	path1 := path.Join(ROOT, WalletPath.LOCAL, "btc.yml")
	provide1, err := bcconfig.FromConfigString(path1, YmlExt)
	if err != nil {
		panic(err)
	}
	Btc = client.NewBtcClient(provide1)
	//////////usdt
	path3 := path.Join(ROOT, WalletPath.LOCAL, "usdt.yml")
	provide3, err := bcconfig.FromConfigString(path3, YmlExt)
	Usdt = client.NewUsdtClient(provide3)
	//////////ETH
	path2 := path.Join(ROOT, WalletPath.LOCAL, "eth.yml")
	provide2, err := bcconfig.FromConfigString(path2, YmlExt)

	if err != nil {
		panic(err)
	}
	//////////ERC20
	Eth = client.NewEthClient(provide2)

	Erc20 = client.NewErc20Client(provide2)
	Erc20Token, _ = token.GetErc20TokenInstance()
	/////////LTC
	path4 := path.Join(ROOT, WalletPath.LOCAL, "ltc.yml")
	provide4, err := bcconfig.FromConfigString(path4, YmlExt)
	Ltc = client.NewLtcClient(provide4)

	go daemon.Start(DaemonCnf)
	log.Println("Everything is ready！！！")
}

var (
	Cache      *bccoin.CoinCache
	keyNet     bccore.Net
	KeyUtil    *bckey.KeyUtil
	ROOT       = os.Getenv("PWD")
	WalletPath = init_1_Path()

	YmlExt = "yml"
	///client
	Btc   *client.BtcClient
	Usdt  *client.UsdtClient
	Ltc   *client.LtcClient
	Eth   *client.EthClient
	Erc20 *client.Erc20Client

	Erc20Token *token.Erc20Token
	//DB
	DaemonCnf = initDaemonCnf()
	TxStorage = initTxStorage()
)

func init_1_Path() *bcconfig.WalletPath {
	path := path.Join(ROOT, "path.yml")
	return bcconfig.InitPath(path, YmlExt)
}

//init kvdb
func initBadger() {
	path := path.Join(ROOT, WalletPath.KVDB)
	err := db.Open(path)
	if err != nil {
		panic(err)
	}
}

func initCoinCache() *bccoin.CoinCache {
	path1 := path.Join(ROOT, WalletPath.COINJSON, "coin_info.json")
	return bccoin.InitCoinCache(db.GetStore(), db.Coin_Info, path1)
}
func initKeyUtil() *bckey.KeyUtil {
	path1 := path.Join(ROOT, WalletPath.COMMON, "config.yml")
	provide, err := bcconfig.FromConfigString(path1, YmlExt)
	if err != nil {
		panic(err)
	}
	conf, err := keynet.DecodeConfig(provide)

	if err != nil || (conf.Net != 1 && conf.Net != 2) {
		conf.Net = 2
	}
	return bckey.InitKeyUtil(db.GetStore(), db.Pfk_PubKey, db.Pfk_Pubkey_Count, conf.Net)
}
func initTxStorage() *mysqlDB.TxStorage {
	path1 := path.Join(ROOT, WalletPath.COMMON, "config.yml")
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
	path1 := path.Join(ROOT, WalletPath.OFFIC, "btc.yml")
	btcp, err := bcconfig.FromConfigString(path1, YmlExt)
	if err != nil {
		panic(err)
	}
	official.InitBtcNode(btcp)
	path2 := path.Join(ROOT, WalletPath.OFFIC, "ltc.yml")
	ltcp, err := bcconfig.FromConfigString(path2, YmlExt)
	if err != nil {
		panic(err)
	}
	official.InitLtcNode(ltcp)
	path3 := path.Join(ROOT, WalletPath.OFFIC, "eth.yml")
	Ethp, err := bcconfig.FromConfigString(path3, YmlExt)
	if err != nil {
		panic(err)
	}
	official.InitEthNode(Ethp)
	path4 := path.Join(ROOT, WalletPath.OFFIC, "usdt.yml")
	usdtp, err := bcconfig.FromConfigString(path4, YmlExt)
	if err != nil {
		panic(err)
	}
	official.InitUsdtNode(usdtp)
	official.InitOfficDaemons(DaemonCnf)
}
func initDaemonCnf() offcnf.Config {
	path1 := path.Join(ROOT, WalletPath.COMMON, "daemon_cnf.yml")
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
