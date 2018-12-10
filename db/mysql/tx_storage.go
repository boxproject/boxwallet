package mysql

import (
	"time"

	"database/sql/driver"

	"encoding/json"

	"github.com/boxproject/boxwallet/bcconfig/mysql"
	"github.com/boxproject/boxwallet/db"
	"github.com/jinzhu/gorm"
)

func (c *TxObj) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *TxObj) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}

//10 attribute
type TxInfo struct {
	TxId          string `gorm:"column:txId;primary_key"`
	Confirmations uint64 `gorm:"column:confirm"`
	Token         string `gorm:"column:token"`
	BCType        int    `gorm:"column:type"`
	BlockH        uint64 `gorm:"column:height"`
	Target        uint64 `gorm:"column:target"`

	Fee   string `gorm:"size:50"`
	TxObj *TxObj `gorm:"type:json"`

	ExtValid bool `gorm:"colum:extValid"`

	CreatedAt  time.Time `gorm:"index"`
	UpdatedAt  time.Time
	DeleteddAt *time.Time
}

type TxObj struct {
	In     []*AddrAmount
	Out    []*AddrAmount
	ExtIn  []*AddrAmount
	ExtOut []*AddrAmount
}

type AddrAmount struct {
	Addr string
	Amt  string
}

var txStorageIntance *TxStorage

func NewTxStorage(conf mysql.MySqlConf) *TxStorage {
	if txStorageIntance == nil {
		txStorageIntance = &TxStorage{
			db: db.MysqlConn(conf.Link, conf.Limit),
		}
		if !txStorageIntance.HasTable(&TxInfo{}) {
			txStorageIntance.CreateTable(&TxInfo{})
		}
	}
	return txStorageIntance
}
func GetTxStorageInstance() *TxStorage {
	return txStorageIntance
}

type TxStorage struct {
	db *gorm.DB
}

func (t *TxStorage) GetTx(txId string) *TxInfo {
	var tx TxInfo
	t.db.Where(&TxInfo{TxId: txId}).First(&tx)
	if tx.TxId == "" {
		return nil
	}
	return &tx
}

func (t *TxStorage) AddTx(tx *TxInfo) bool {
	tx.CreatedAt = time.Now()
	tx.UpdatedAt = time.Now()
	tx.DeleteddAt = nil
	return t.db.Create(&tx).Error == nil
}

func (t *TxStorage) UpdateTx(txId string, confirm uint64, extValid bool) bool {
	tx := &TxInfo{
		TxId: txId,
	}
	result := t.db.Model(&tx).Update(TxInfo{Confirmations: confirm, ExtValid: extValid})
	return result.RowsAffected == 1
}

func (t *TxStorage) HasTable(v interface{}) bool {
	return t.db.HasTable(v)
}

func (t *TxStorage) CreateTable(v interface{}) bool {
	return t.db.CreateTable(v).Error == nil
}
