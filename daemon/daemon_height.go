package daemon

import (
	"path"

	"io/ioutil"

	"math/big"

	"log"

	"strconv"

	"strings"

	"github.com/boxproject/boxwallet/bccore"
	"github.com/boxproject/boxwallet/util"
)

var logPath map[bccore.BloclChainType]string

func InitLogPath(dir string) {
	logPath = make(map[bccore.BloclChainType]string)
	logPath[bccore.BC_BTC] = path.Join(dir, "btc.log")
	logPath[bccore.BC_ETH] = path.Join(dir, "eth.log")
	logPath[bccore.BC_LTC] = path.Join(dir, "ltc.log")
}

func HeightWriteAsyn(chainType bccore.BloclChainType, height *big.Int) {

	go func(chainType bccore.BloclChainType, height *big.Int) {
		defer util.CatchPanic()
		if height.Uint64()%5 == 0 {
			data := []byte(height.String())
			err := ioutil.WriteFile(logPath[chainType], data, 0644)
			if err != nil {
				log.Println("🖊chainType:", chainType, ",block height write failed")
			} else {
				log.Println("🖊chainType:", chainType, ",block height write:", height.String())
			}
		}
	}(chainType, height)
}
func HeightRead(chainType bccore.BloclChainType) *big.Int {
	data, err := ioutil.ReadFile(logPath[chainType])
	height := big.NewInt(0)
	if err != nil || data == nil {
		return height
	}
	str := string(data)
	str = strings.TrimSpace(str)
	num, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return height
	}
	height = height.SetUint64(num)
	log.Println("type:", chainType, ",height:", height.String())
	if height == nil {
		height = big.NewInt(0)
	}
	return height
}
