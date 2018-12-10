package bccoin

import "testing"

func TestLoadCoinInfo(t *testing.T) {
	cis := loadCoinInfo("./distribute/coin_info.json")
	for _, v := range cis {
		t.Log(v.Name, v.CT, v.Decimals, v.Symbol, v.Token)
	}
}
