package clientseries

import (
	"github.com/boxproject/boxwallet/bcconfig"
	rpcclient "github.com/boxproject/lib-bitcore/selrpcclient"
	"github.com/ltcsuite/ltcd/chaincfg"
)

type LtcSeriesClient struct {
	C        *rpcclient.Client
	Env      *chaincfg.Params
	official *rpcclient.Client
}

func NewLtcSeriesClient(cfg bcconfig.Provider) *LtcSeriesClient {
	conf := struct {
		ip   string
		port string
		user string
		pass string
		net  string
	}{}

	conf.ip = cfg.GetString("ip")
	conf.port = cfg.GetString("port")
	conf.net = cfg.GetString("net")
	conf.user = cfg.GetString("user")
	conf.pass = cfg.GetString("passwd")

	btconf := &rpcclient.ConnConfig{
		Host:         conf.ip + ":" + conf.port,
		User:         conf.user,
		Pass:         conf.pass,
		HTTPPostMode: true,
		DisableTLS:   true,
	}
	intance := &LtcSeriesClient{}

	c, err := rpcclient.New(btconf, nil)
	if err != nil {
		panic("BtcSeriesClient initialization failed ")
	}
	intance.C = c
	switch conf.net {
	case "main":
		intance.Env = &chaincfg.MainNetParams
		break
	case "test":
		intance.Env = &chaincfg.TestNet4Params
		break
	case "regtest":
		intance.Env = &chaincfg.RegressionNetParams
		break
	default:
		intance.Env = &chaincfg.RegressionNetParams
		break
	}
	return intance

}
