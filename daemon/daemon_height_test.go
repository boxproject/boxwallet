package daemon_test

import (
	"testing"

	"os"

	"io/ioutil"

	"math/big"

	"time"

	"strconv"

	"strings"

	"github.com/boxproject/boxwallet/bccore"
	"github.com/boxproject/boxwallet/daemon"
)

var logs map[bccore.BloclChainType]string

func TestDir(t *testing.T) {
	logs = make(map[bccore.BloclChainType]string)
	fi, err := os.Open("/Users/rennbon/go/src/github.com/boxproject/boxwallet/bcconfig/blockheight/btc.log")
	if err != nil {
		panic("err")
	}
	t.Log(fi.Name())
}

func TestHeightWriteAsyn(t *testing.T) {
	daemon.HeightWriteAsyn(bccore.BC_ETH, big.NewInt(20200))
	time.Sleep(time.Second * 2)
}
func TestReadBlockHeight(t *testing.T) {
	name := "/Users/rennbon/go/src/github.com/boxproject/boxwallet/bcconfig/blockheight/ltc.log"
	data, err := ioutil.ReadFile(name)
	if err != nil {
		panic(err)
	}
	str := string(data)
	str = strings.TrimSpace(str)
	t.Log(str)
	num, err := strconv.ParseUint(str, 10, 64)
	t.Log(num)
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}
	height := big.NewInt(0)
	height = height.SetUint64(num)
	t.Log(height.Uint64())
}

func TestMol(t *testing.T) {
	a := big.NewInt(20)
	m := a.Int64() % 3
	t.Log(m)
}
