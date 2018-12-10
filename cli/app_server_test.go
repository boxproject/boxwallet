package cli_test

import (
	"os/exec"
	"testing"

	"os"

	"path/filepath"
	"strings"

	"log"

	"github.com/boxproject/boxwallet/bccore"
	"github.com/boxproject/boxwallet/cli"
)

func TestNewAppServer(t *testing.T) {
	app := cli.NewAppServer()
	t.Log(app)

	//app.SaveMasterKey("tpubD6NzVbkrYhZ4XMziCvXd7MeobjcAkaCQZwErnt5qSeEkNfGBr7whYdXBcxGbrPnheKrcaTSpQuWkZh4VT2hd6fLj89Z8PUD7CiwtwXMRcmj")

	addr, _ := app.GetMasterAddress(bccore.BC_LTC)
	t.Log(addr.Address())
	/*
		keys, err := app.GeneraterKeys(10, bccore.STR_LTC)
		if err != nil {
			t.Fail()
			t.Error(err)
			return
		}
		for key, value := range keys {
			fmt.Printf("key:%v,[address=%v][CustomDeep=%v][CurrentNum=%v][IsMaster=%v][Key=%v]\n",
				key,
				value.Address(),
				value.CustomDeep(),
				value.CurrentNum(),
				value.IsMaster(),
				//value.Key(),
			)
		}
		t.Log(len(keys))
		t.Log(app.GetChildKeyCount())*/
	//select {}
}

func TestPath(t *testing.T) {
	t.Log(os.Args)
}
func GetAppPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	log.Println(file, path)
	return path[:index]
}
