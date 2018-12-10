package mysql

import (
	"github.com/boxproject/boxwallet/bcconfig"
	"github.com/mitchellh/mapstructure"
)

const (
	keyStorageConfigKey = "mysql"
)

type Config struct {
	KeyStorage MySqlConf
}

func DecodeConfig(cfg bcconfig.Provider) (c Config, err error) {
	m := cfg.GetStringMap(keyStorageConfigKey)
	err = mapstructure.WeakDecode(m, &c)
	return
}

type MySqlConf struct {
	Link  string
	Limit int
}
