package signature

import (
	"crypto/ecdsa"

	"encoding/json"

	"math/big"

	"github.com/boxproject/boxwallet/errors"

	"encoding/hex"
	"strings"

	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type EthTx struct {
	Nonce    uint64         `josn:"nonce"`
	To       common.Address `json:"to"`
	Amount   *big.Int       `json:"amount"`
	GasLimit uint64         `json:"gaslimit"`
	GasPrice *big.Int       `json:"gasprice"`
	Data     []byte         `json:"data"`
	Signed   bool           `json:"signed"`
	From     common.Address `json:"from"`
	T        *types.Transaction
	Loc      bool
}

const tranferId = "a9059cbb"

func NewEthTx(nonce uint64, from, to common.Address, amount *big.Int, gaslimit uint64, gasprice *big.Int, data []byte, local bool) *EthTx {
	return &EthTx{Nonce: nonce, To: to, Amount: amount, GasLimit: gaslimit, GasPrice: gasprice, Data: data, From: from, Loc: local}
}
func (tx *EthTx) FromAddresses() (froms []string) {
	froms = make([]string, 0, 1)
	froms = append(froms, tx.From.String())
	return
}

func (tx *EthTx) Info(addrFrom string) (to []*AddressAmount, err error) {
	addrTo := ""
	amount := new(big.Int)
	if len(tx.Data) == 0 {
		//Ordinary transfer
		addrTo = tx.To.String()
		amount = tx.Amount
	} else {
		//tokens transfer
		//MethodID: 0xa9059cbb
		//address: 64 byte
		//amount: 64 byte
		hexStr := hex.EncodeToString(tx.Data)
		if hexStr[:8] != tranferId {
			//todo
			return nil, nil
		}
		addrTo = strings.ToLower("0x" + hexStr[32:72])
		value := hexStr[72:]

		val := big.NewInt(0)
		val, _ = val.SetString(value, 16)
		amount = val
	}

	to = []*AddressAmount{
		{
			Amount:  amount,
			Address: addrTo,
		}}
	return to, nil
}
func (tx *EthTx) IsSign() bool {
	return tx.Signed
}
func (tx *EthTx) TxId() string {
	if !tx.IsSign() {
		return ""
	}
	return tx.T.Hash().String()
}

func (tx *EthTx) Sign(prvKeys []string) error {
	prvkey, err := tx.stringToEthPrvKey(prvKeys[0])
	if err != nil {
		return err
	}
	signTx, err := types.SignTx(tx.T, types.HomesteadSigner{}, prvkey)
	if err != nil {
		return err
	}
	tx.T = signTx
	tx.Signed = true
	return nil
}
func (tx *EthTx) TxForSend() (v interface{}, err error) {
	if tx.IsSign() {
		return tx.T, nil
	}
	return nil, errors.ERR_TX_WITHOUT_SGIN
}

func (tx *EthTx) stringToEthPrvKey(prvkey string) (*ecdsa.PrivateKey, error) {
	hdPrvkey, err := hdkeychain.NewKeyFromString(prvkey)
	if err != nil {
		return nil, err
	}
	prv, err := hdPrvkey.ECPrivKey()
	if err != nil {
		return nil, err
	}
	ethkey := new(ecdsa.PrivateKey)
	ethkey.D = prv.D
	ethkey.PublicKey.Curve = prv.Curve
	ethkey.PublicKey.X, ethkey.PublicKey.Y = prv.X, prv.Y
	return ethkey, nil
}

func (tx *EthTx) Marshal() ([]byte, error) {
	return json.Marshal(tx)
}
func (tx *EthTx) Local() bool {
	return tx.Loc

}
func (tx *EthTx) UnMarshal(data []byte) error {
	err := json.Unmarshal(data, tx)
	if err != nil {
		return err
	}
	if tx.T == nil {
		tx.T = types.NewTransaction(
			tx.Nonce,
			tx.To,
			tx.Amount,
			tx.GasLimit,
			tx.GasPrice,
			tx.Data,
		)
	}
	return nil
}
