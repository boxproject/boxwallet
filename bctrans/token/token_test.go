package token_test

import (
	"testing"

	"github.com/boxproject/boxwallet/bccore"
	"github.com/boxproject/boxwallet/bctrans/token"
	"github.com/boxproject/boxwallet/mock"
)

func Test_GetTokenInfo(t *testing.T) {
	info, err := mock.Erc20Token.GetTokenInfo("0xFCA050c8369ED68d38034746eBA8Ad16aF89FB6A")
	if err != nil {
		t.Fail()
		t.Error(err)
	} else {
		t.Log(info)
	}
}

func TestGetToken(t *testing.T) {
	arr := token.GetToken(bccore.BC_ERC20)
	for _, v := range *arr {
		t.Log(v)
	}
}
