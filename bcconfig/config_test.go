package bcconfig_test

import (
	"testing"

	"github.com/boxproject/boxwallet/bcconfig"
)

func TestFromConfigString(t *testing.T) {
	provide, err := bcconfig.FromConfigString("btc.yml", "yml")
	if err != nil {
		t.Fail()
		t.Error(err)
	}
	t.Log(provide)

	t.Log(provide.Get("aee"))
}
func TestGetStringSlicePreserveString(t *testing.T) {
	/*	provide,err:= bcconfig.FromConfigString("btc.yml","yml")
		if err!=nil{
			t.Fail()
			t.Error(err)
		}
		bcconfig.GetStringSlicePreserveString(provide,"passwd")*/
}
