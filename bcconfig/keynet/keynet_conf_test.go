package keynet

import (
	"testing"

	"path"

	"github.com/boxproject/boxwallet/bcconfig"
)

func TestDecodeConfig(t *testing.T) {
	path1 := path.Join("../", "config.yml")
	provide, err := bcconfig.FromConfigString(path1, "yml")
	if err != nil {
		t.Fail()
		t.Error(err)
		return
	}
	conf, err := DecodeConfig(provide)
	if err != nil {
		t.Fail()
		t.Error(err)
		return
	}
	t.Log(conf)
}
