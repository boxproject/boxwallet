package clientseries_test

import (
	"path"
	"testing"

	"github.com/boxproject/boxwallet/bcconfig"
	"github.com/boxproject/boxwallet/bctrans/clientseries"
	"github.com/boxproject/boxwallet/mock"
)

func TestNewEthSeriesClient(t *testing.T) {
	path := path.Join(mock.Gopath, mock.ProjectDir, mock.ConfigDir, "eth.yml")
	provide, err := bcconfig.FromConfigString(path, mock.YmlExt)
	if err != nil {
		t.Fail()
		t.Error(err)
	}
	client := clientseries.NewEthSeriesClient(provide)
	t.Log(client)
}

func TestNewUsdtSeriesClient(t *testing.T) {
	path := path.Join(mock.Gopath, mock.ProjectDir, mock.ConfigDir, "usdt.yml")
	provide, err := bcconfig.FromConfigString(path, mock.YmlExt)
	if err != nil {
		t.Fail()
		t.Error(err)
	}
	client := clientseries.NewOmniSeriesClient(provide)
	t.Log(client)
}
func TestNewLtcClient(t *testing.T) {
	path := path.Join(mock.Gopath, mock.ProjectDir, mock.ConfigDir, "ltc.yml")
	provide, err := bcconfig.FromConfigString(path, mock.YmlExt)
	if err != nil {
		t.Fail()
		t.Error(err)
	}
	client := clientseries.NewLtcSeriesClient(provide)
	/*txhash, _ := chainhash.NewHashFromStr("29462f45c42e1e0097b87c7a0aebc13184f1192d84b2723b9b747ff43ed17472")
	tx, err := client.C.GetTransaction(txhash)*/
	t.Log(client.C.GetBlockCount())
}
