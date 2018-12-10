package db

type PrefixKey []byte

var (
	//key prifex
	Pfk_PubKey       PrefixKey = []byte("pubkey_")  //for pubkey
	Pfk_Pubkey_Count PrefixKey = []byte("pubkeyc_") //for pubkey count

	//coininfo prifex
	Coin_Info PrefixKey = []byte("ci_")
)
