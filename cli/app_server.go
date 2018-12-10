package cli

import (
	"github.com/boxproject/boxwallet/bccoin"
	"github.com/boxproject/boxwallet/bccore"
	"github.com/boxproject/boxwallet/bckey"
	"github.com/boxproject/boxwallet/bctrans"
	"github.com/boxproject/boxwallet/bctrans/token"
	"github.com/boxproject/boxwallet/pipeline"
	"github.com/boxproject/boxwallet/signature"
)

type AppServer struct {
	trans       map[bccore.BlockChainSign]*bctrans.Trans
	token       *token.Erc20Token
	keyUtil     *bckey.KeyUtil
	heightCache *pipeline.HeightCache
}

func NewAppServer() *AppServer {
	app := &AppServer{}
	var err error
	app.heightCache = pipeline.GetHeightCacheInstance()
	app.trans = make(map[bccore.BlockChainSign]*bctrans.Trans)
	app.trans[bccore.STR_BTC], err = bctrans.NewTrans(bccore.STR_BTC, DaemonCnf.Btc.Unlock != 0)
	if err != nil {
		panic(err)
	}
	app.trans[bccore.STR_ETH], err = bctrans.NewTrans(bccore.STR_ETH, DaemonCnf.Eth.Unlock != 0)
	if err != nil {
		panic(err)
	}
	app.trans[bccore.STR_ERC20], err = bctrans.NewTrans(bccore.STR_ERC20, DaemonCnf.Eth.Unlock != 0)
	if err != nil {
		panic(err)
	}
	app.trans[bccore.STR_USDT], err = bctrans.NewTrans(bccore.STR_USDT, DaemonCnf.Btc.Unlock != 0)
	if err != nil {
		panic(err)
	}
	app.trans[bccore.STR_LTC], err = bctrans.NewTrans(bccore.STR_LTC, DaemonCnf.Ltc.Unlock != 0)
	if err != nil {
		panic(err)
	}
	app.token, err = token.GetErc20TokenInstance()
	if err != nil {
		panic(err)
	}
	app.keyUtil = bckey.GetKeyUtilInstance()
	return app
}

func (app *AppServer) GetBalance(sign bccore.BlockChainSign, address string, token string) (balance bccoin.CoinAmounter, err error) {
	return app.trans[sign].GetBalance(address, token)
}

func (app *AppServer) CreateTx(sign bccore.BlockChainSign, addrsFrom []string, token string, addrsTo []*bccoin.AddressAmount, feeCeo float64) (uuid string, txu signature.TxUtil, err error) {
	return app.trans[sign].CreateTx(addrsFrom, token, addrsTo, feeCeo)
}

func (app *AppServer) SendTx(sign bccore.BlockChainSign, txu signature.TxUtil, uuid string) error {
	return app.trans[sign].SendTx(txu, uuid)
}

func (app *AppServer) GetCoinInfo(chainType bccore.BloclChainType, token string) (*bccoin.CoinInfo, error) {
	return app.token.GetTokenInfo(bccore.Token(token))
}

func (app *AppServer) GetMasterKey() {
	app.keyUtil.GetMasterKey()
}

func (app *AppServer) SaveMasterKey(pubkey string) error {
	key, _ := app.keyUtil.GetMasterKey()
	if key == nil {
		err := app.keyUtil.SaveMasterKey(pubkey)
		if err != nil {
			return err
		}
	}
	for k, v := range app.trans {
		mkey, err := app.keyUtil.GetMasterGenericKey(bccore.BCMap[k])
		if err != nil {
			return err
		}
		v.Walleter.ImportAddress(mkey.Address(), true)
	}
	return nil
}

func (app *AppServer) GetMasterAddress(bct bccore.BloclChainType) (bckey.GenericKey, error) {
	return app.keyUtil.GetMasterGenericKey(bct)
}

func (app *AppServer) GeneraterKeys(num int, sign bccore.BlockChainSign) (keys []bckey.GenericKey, err error) {
	keys = make([]bckey.GenericKey, 0, num)

	for i := 0; i < num; i++ {
		key, err := app.trans[sign].Walleter.GetNewAddress()
		if err != nil {
			return keys, err
		}
		keys = append(keys, key)
	}

	return
}

//Gets the non-primary public key
//withCurNum:true=>current，false=>parent
func (app *AppServer) GetKey(customDeep []uint32, bc bccore.BloclChainType, curNum uint32, withCurNum bool) (key bckey.GenericKey, err error) {
	return app.keyUtil.GetKey(customDeep, bc, curNum, withCurNum)
}

//Gets the total number of public keys that have been generated（Not including masterkey）
func (app *AppServer) GetChildKeyCount() (int, error) {
	return app.keyUtil.GetChildKeyCount()
}

func (app *AppServer) GetHeights() map[bccore.BloclChainType]*pipeline.Heights {
	return app.heightCache.LoadAll()
}
