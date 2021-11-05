package util

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var MysqlDSN string
var DB *sql.DB

func init() {
	// InitConf()
	MysqlDSN = GetConfStr("mysqlDSN")
	DB = GetDB()
}

func GetDB() *sql.DB {
	DB, err := sql.Open("mysql", MysqlDSN)
	if err != nil {
		DB = nil
	}
	return DB
}
