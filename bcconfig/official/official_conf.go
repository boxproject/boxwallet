package official

import (
	"github.com/boxproject/boxwallet/bcconfig"
	"github.com/mitchellh/mapstructure"
)

type Config struct {
	Btc DaemonCnf
	Ltc DaemonCnf
	Eth DaemonCnf
}

func DecodeConfig(cfg bcconfig.Provider) (c Config, err error) {
	m := cfg.GetStringMap("officDaemon")
	err = mapstructure.WeakDecode(m, &c)
	return
}

type DaemonCnf struct {
	Confirm uint64
	Gap     int
	Ticker  int64
	Unlock  uint64
}
