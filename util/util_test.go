package util_test

import (
	"testing"

	"math/big"

	"github.com/boxproject/boxwallet/util"
)

var su util.StrUtil

/////////////////////测试用实体///////////////////////////
type testModle struct {
	from, to, target string
	gap, index       int
}

var slc = []testModle{
	{index: 0, target: "<< 2 and not overflow", from: "12345.6789", to: "123.456789", gap: 2},
	{index: 1, target: ">> 2 and not overflow", from: "12345.6789", to: "1234567.89", gap: -2},
	{index: 2, target: "<< 10 and overflow append 0", from: "12345.6789", to: "0.00000123456789", gap: 10},
	{index: 3, target: ">> 10 and overflow append 0", from: "12345.6789", to: "123456789000000", gap: -10},
	{index: 4, target: "just <<", from: "12345.6789", to: "0.123456789", gap: 5},
	{index: 5, target: "just >>", from: "12345.6789", to: "123456789", gap: -4},
}

////////////////////////////////////////////////////////

func TestStrUtil_SplitStrToNum(t *testing.T) {
	resmap := map[string]bool{
		"12312312332":                           true,
		"1111111.33333":                         true,
		"812391231237123123.111123123812399922": true,
		"000000.0000000":                        true,
		"a11231312313":                          false,
		"123123hello.1231231":                   false,
		"box test error":                        false,
		"1231^&.123123":                         false,
	}
	failure := make([]string, 0, 0)
	for k, v := range resmap {
		l, f, err := su.SplitStrToNum(k, true)
		if (err == nil) == v {
			t.Logf("string:'%s'\r\nsccuess\r\nleft:%s\r\nright:%s\r\n", k, l, f)
		} else {
			t.Errorf("string:'%s'\r\nfailed，error:%s\r\n", k, err)
			failure = append(failure, k)
		}
	}
	if len(failure) > 0 {
		t.Fail()
	}
}
func BenchmarkStrUtil_SplitStrToNum(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ { //use b.N for looping
		su.SplitStrToNum("12345.67890", false)
	}
}

func TestStrUtil_MoveDecimalPosition(t *testing.T) {
	for k, v := range slc {
		str, err := su.MoveDecimalPosition(v.from, v.gap, false)
		if err != nil {
			t.Errorf("sccuess，index:%d\r\n")
			t.Fail()
		} else if str != v.to {
			t.Errorf("failed,index：%d\r\n param:%s\r\n target:%s\r\n Gap:%d\r\n actul:%s\r\n target:%s\r\n", k, v.from, v.to, v.gap, str, v.target)
			t.Fail()
		} else {
			t.Logf("sccuess,index：%d\r\n param:%s\r\n target:%s\r\n Gap:%d\r\n actul:%s\r\n target:%s\r\n", k, v.from, v.to, v.gap, str, v.target)
		}
	}

}

func BenchmarkStrUtil_MoveDecimalPosition(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ { //use b.N for looping
		for _, v := range slc {
			su.MoveDecimalPosition(v.from, v.gap, false)

		}
	}
}
func TestBig(t *testing.T) {
	f := big.NewFloat(2.5)

	t.Log(f.Int64())
}
