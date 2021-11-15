package dao

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"ssq-spider/configure"
)

type MysqlDB struct {
	dbUrl    string
	dbClient *gorm.DB
}

var mysqlDB MysqlDB

func init() {
	mysqlDB.dbUrl = configure.GlobalConfig.Mysql.Url
	mysqlDB.dbClient = nil
}

func NewMysqlDBClient() (*gorm.DB, error) {
	if mysqlDB.dbClient != nil {
		if err := mysqlDB.dbClient.DB().Ping(); err == nil {
			return mysqlDB.dbClient, nil
		}
		_ = mysqlDB.dbClient.Close()
	}

	db, err := gorm.Open("mysql", configure.GlobalConfig.Mysql.Url)
	if err != nil {
		panic("open mysql connect error: " + err.Error())
	}

	mysqlDB.dbClient = db

	return mysqlDB.dbClient, nil
}
