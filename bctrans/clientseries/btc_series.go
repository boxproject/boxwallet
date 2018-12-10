package clientseries

import (
	"github.com/boxproject/boxwallet/bcconfig"
	rpcclient "github.com/boxproject/lib-bitcore/serpcclient"
	"github.com/btcsuite/btcd/chaincfg"
)

type OmniSeriesClient struct {
	C          *rpcclient.Client
	Env        *chaincfg.Params
	PropertyId int
}

func NewOmniSeriesClient(cfg bcconfig.Provider) *OmniSeriesClient {
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
	intance := &OmniSeriesClient{}

	c, err := rpcclient.New(btconf, nil)
	if err != nil {
		panic("BtcSeriesClient initialization failed ")
	}
	intance.PropertyId = cfg.GetInt("ppId")
	intance.C = c
	switch conf.net {
	case "main":
		intance.Env = &chaincfg.MainNetParams
		break
	case "test":
		intance.Env = &chaincfg.TestNet3Params
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
