package bcconfig

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

type Provider interface {
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringSlice(key string) []string
	Get(key string) interface{}
	Set(key string, value interface{})
	IsSet(key string) bool
	WatchConfig()
	OnConfigChange(run func(in fsnotify.Event))
	Unmarshal(rawVal interface{}, opts ...viper.DecoderConfigOption) error
}

// FromConfigString creates a config from the given YAML, JSON or TOML config. This is useful in tests.
func FromConfigString(path, configType string) (Provider, error) {
	v := viper.New()
	v.SetConfigType(configType)
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	return v, nil
}

// GetStringSlicePreserveString returns a string slice from the given config and key.
// It differs from the GetStringSlice method in that if the config value is a string,
// we do not attempt to split it into fields.
func GetStringSlicePreserveString(cfg Provider, key string) []string {
	sd := cfg.Get(key)
	if sds, ok := sd.(string); ok {
		return []string{sds}
	} else {
		return cast.ToStringSlice(sd)
	}
}
