package bcconfig

type WalletPath struct {
	LOCAL    string
	OFFIC    string
	KVDB     string
	COINJSON string
	COMMON   string
	HEIGHT   string
}

func InitPath(path string, fileType string) *WalletPath {
	provide, err := FromConfigString(path, fileType)
	if err != nil {
		panic(err)
	}
	cnf := &WalletPath{
		LOCAL:    provide.GetString("LOCAL"),
		OFFIC:    provide.GetString("OFFIC"),
		KVDB:     provide.GetString("KVDB"),
		COINJSON: provide.GetString("COINJSON"),
		COMMON:   provide.GetString("COMMON"),
		HEIGHT:   provide.GetString("HEIGHT"),
	}
	return cnf
}
