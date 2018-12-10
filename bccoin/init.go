package bccoin

import (
	"encoding/json"
	"io/ioutil"
)

func loadCoinInfo(path string) []*CoinInfo {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	cis := new([]*CoinInfo)
	err = json.Unmarshal(data, cis)
	if err != nil {
		panic(err)
	}
	return *cis
}
