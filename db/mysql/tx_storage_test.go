package mysql_test

import (
	"testing"

	"time"

	mysqlDB "github.com/boxproject/boxwallet/db/mysql"
	"github.com/boxproject/boxwallet/mock"
)

func TestCreateTable(t *testing.T) {
	if !mock.TxStorage.HasTable(&mysqlDB.TxInfo{}) {
		mock.TxStorage.CreateTable(&mysqlDB.TxInfo{})
	}
}
func TestTxStorage_AddTx(t *testing.T) {
	txc := &mysqlDB.TxInfo{
		TxId:          "test1",
		Confirmations: 20,
		Token:         "",
		BCType:        1,
		BlockH:        21,

		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		DeleteddAt: nil,
		TxObj:      &mysqlDB.TxObj{},
		Fee:        "0.0001",
	}
	txc.TxObj.In = append(txc.TxObj.In, &mysqlDB.AddrAmount{
		Addr: "address from",
		Amt:  "12000",
	})
	txc.TxObj.Out = append(txc.TxObj.Out, &mysqlDB.AddrAmount{
		Addr: "address to",
		Amt:  "11999.9999",
	})

	t.Log(mock.TxStorage.AddTx(txc))
}
func TestTxStorage_GetTx(t *testing.T) {
	txc := mock.TxStorage.GetTx("test2")
	if txc != nil {
		t.Log(txc.Confirmations)
	}
}

func TestTxStorage_UpdateTx(t *testing.T) {
	result := mock.TxStorage.UpdateTx("test1", 200, false)
	if !result {
		t.Fail()
	}
}
