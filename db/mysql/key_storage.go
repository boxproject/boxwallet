package mysql

import (
	"github.com/boxproject/boxwallet/bcconfig/mysql"
	"github.com/boxproject/boxwallet/db"
	"github.com/jinzhu/gorm"
)

type Key struct {
	gorm.Model
	Name    string `gorm:"size:255;not null"`
	Address string `gorm:"size:255;not null`
	BCType  int    `gorm:not null`
}

var keyStorageIntance *KeyStorage

func NewKeyStorage(conf mysql.MySqlConf) *KeyStorage {
	if keyStorageIntance == nil {
		keyStorageIntance = &KeyStorage{
			DB: db.MysqlConn(conf.Link, conf.Limit),
		}
	}
	return keyStorageIntance
}

type KeyStorage struct {
	*gorm.DB
}

func (k *KeyStorage) AddKey(name, address string, bcType int) bool {
	key := &Key{
		Name:    name,
		Address: address,
		BCType:  bcType,
	}
	if !k.DB.NewRecord(&key) {
		return false
	}
	k.DB.Create(&key)
	return !k.DB.NewRecord(&key)
}
