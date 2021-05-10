package store

import (
	"github.com/codelieche/microservice/usercenter/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var mysqldb *gorm.DB

func GetMySQLDB(cfg *config.GlobalConfig) *gorm.DB {
	var err error
	if mysqldb != nil {
		return mysqldb
	} else {
		mysqlDsn := cfg.MySQL.GetDSN()
		log.Printf(mysqlDsn)

		cfg := &gorm.Config{}
		mysqldb, err = gorm.Open(mysql.Open(mysqlDsn), cfg)
		if err != nil {
			panic(err)
			return nil
		} else {
			return mysqldb
		}
	}
}
