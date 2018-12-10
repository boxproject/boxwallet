package token

import (
	"github.com/boxproject/boxwallet/bccoin"
	"github.com/boxproject/boxwallet/bccore"
)

type TokenLawer interface {
	GetTokenInfo(contract bccore.Token) (*bccoin.CoinInfo, error)
}
