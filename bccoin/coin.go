package bccoin

import (
	"bytes"
	"math/big"

	"strconv"

	"encoding/json"

	"github.com/boxproject/boxwallet/bccore"
	"github.com/boxproject/boxwallet/errors"
)

//Details of the coin
type CoinInfo struct {
	CT       bccore.BloclChainType `json:"coinType"` //main chain type
	Token    bccore.Token          `json:"token"`    //contract address
	Symbol   string                `json:"symbol"`
	Decimals int                   `json:"decimals"` //precision
	Name     string                `json:"name"`
}

type CoinSign interface {
	Sign() []byte
	Bytes() []byte
}

type AddressAmount struct {
	Address string
	Amount  CoinAmounter
}

//暴露给外部的金额计算接口
type CoinAmounter interface {
	Val() *big.Int
	String() string
	///////----Mathematical calculations----start----///////
	Cmp(amount CoinAmounter) (int, error)
	Add(amounts ...CoinAmounter) error
	Sub(amounts ...CoinAmounter) error
	Mul(amounts ...CoinAmounter) error
	Set(amount CoinAmounter)
	///////----Mathematical calculations----end----///////
	CoinSign
}

func (c *CoinInfo) Sign() []byte {
	buff := &bytes.Buffer{}
	buff.WriteString(strconv.FormatInt(int64(c.CT), 10))
	if c.Token != "" {
		buff.WriteString("_")
		buff.WriteString(string(c.Token))
	}
	return buff.Bytes()
}
func (c *CoinInfo) Bytes() []byte {
	v, _ := json.Marshal(c)
	return v
}

type coinAmount struct {
	amount *big.Int
	*CoinInfo
}

func NewCoinAmount(bc bccore.BloclChainType, t bccore.Token, amount string) (CoinAmounter, error) {
	ci, err := defualtCoinCache.GetCoinInfo(bc, t)
	if err != nil {
		return nil, err
	}
	//按照类别实例化不同的amount
	err = regutil.CanPraseBigFloat(amount)
	if err != nil {
		return nil, err
	}
	ca := &coinAmount{
		CoinInfo: ci,
	}
	return ca.stringToAmount(amount)
}
func NewCoinAmountFromInt(bc bccore.BloclChainType, t bccore.Token, amount int64) (CoinAmounter, error) {
	ci, err := defualtCoinCache.GetCoinInfo(bc, t)
	if err != nil {
		return nil, err
	}
	ca := &coinAmount{
		CoinInfo: ci,
	}
	ca.amount = big.NewInt(amount)
	return ca, nil
}
func NewCoinAmountFromBigInt(bc bccore.BloclChainType, t bccore.Token, amount *big.Int) (CoinAmounter, error) {
	ci, err := defualtCoinCache.GetCoinInfo(bc, t)
	if err != nil {
		return nil, err
	}
	ca := &coinAmount{
		CoinInfo: ci,
	}
	ca.amount = big.NewInt(0).Set(amount)
	return ca, nil
}
func NewCoinAmountFromFloat(bc bccore.BloclChainType, t bccore.Token, amount float64) (CoinAmounter, error) {
	ci, err := defualtCoinCache.GetCoinInfo(bc, t)
	if err != nil {
		return nil, err
	}
	ca := &coinAmount{
		CoinInfo: ci,
	}
	return ca.stringToAmount(strconv.FormatFloat(amount, 'f', ci.Decimals, 64))
}

func (c *coinAmount) stringToAmount(amount string) (ca CoinAmounter, err error) {
	l, r, err := strutil.SplitStrToNum(amount, false)
	if err != nil {
		return
	}
	rlen := len(r)
	if rlen > c.Decimals {
		err = errors.ERR_COIN_PREC_OVERFLOW
		return
	}
	buff := &bytes.Buffer{}
	buff.WriteString(l)
	buff.WriteString(r)
	for i := 0; i < c.Decimals-rlen; i++ {
		buff.WriteString("0")
	}
	amt := big.NewInt(0)
	amt.SetString(buff.String(), 10)
	c.amount = amt
	return c, nil
}

func (c *coinAmount) Val() *big.Int {
	return c.amount
}

func (c *coinAmount) String() string {
	str := c.amount.String()
	length := len(str)
	buff := &bytes.Buffer{}
	//将要左移多少位
	if c.Decimals >= length {
		buff.WriteString("0.")
		for i := 0; i < c.Decimals-length; i++ {
			buff.WriteString("0")
		}
		buff.WriteString(str)
	} else {

		buff.WriteString(str[:length-c.Decimals])
		if c.Decimals != 0 {
			buff.WriteString(".")
			buff.WriteString(str[length-c.Decimals:])
		}
	}
	return buff.String()
}

func (c *coinAmount) Set(amount CoinAmounter) {
	c.amount.Set(amount.Val())
}

//c<amount return -1；c>amount return +1；c==amount return 0。
func (c *coinAmount) Cmp(amount CoinAmounter) (int, error) {
	if bytes.Compare(c.Sign(), amount.Sign()) != 0 {
		return 0, errors.ERR_DIFF_UNIT
	}
	return c.amount.Cmp(amount.Val()), nil
}

func (c *coinAmount) Add(amounts ...CoinAmounter) error {
	if amounts == nil {
		return errors.ERR_NIL_REFERENCE
	}
	for _, v := range amounts {
		if bytes.Compare(c.Sign(), v.Sign()) != 0 {
			return errors.ERR_DIFF_UNIT
		}
	}
	for _, v := range amounts {
		c.amount.Add(c.amount, v.Val())
	}
	return nil
}

func (c *coinAmount) Sub(amounts ...CoinAmounter) error {
	if amounts == nil {
		return errors.ERR_NIL_REFERENCE
	}
	for _, v := range amounts {
		if bytes.Compare(c.Sign(), v.Sign()) != 0 {
			return errors.ERR_DIFF_UNIT
		}
	}
	for _, v := range amounts {
		c.amount.Sub(c.amount, v.Val())
	}
	return nil
}

func (c *coinAmount) Mul(amounts ...CoinAmounter) error {
	if amounts == nil {
		return errors.ERR_NIL_REFERENCE
	}
	for _, v := range amounts {
		if bytes.Compare(c.Sign(), v.Sign()) != 0 {
			return errors.ERR_DIFF_UNIT
		}
	}
	for _, v := range amounts {
		c.amount.Mul(c.amount, v.Val())
	}
	return nil
}
func (c *coinAmount) Sign() []byte {
	return c.CoinInfo.Sign()
}
func (c *coinAmount) Bytes() []byte {
	return c.CoinInfo.Bytes()
}
