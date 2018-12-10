package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func MysqlConn(link string, limit int) *gorm.DB {
	//db, err := gorm.Open("mysql", "root:qwe123456@/godb?charset=utf8&parseTime=True")
	db, err := gorm.Open("mysql", link)
	//defer db.Close()
	if err != nil {
		panic(err)
	}
	db.DB().SetMaxIdleConns(limit)
	return db
}
