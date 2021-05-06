package datasources

import (
	"github.com/codelieche/microservice/usercenter/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

func GetMySQLDB() *gorm.DB {
	if db != nil {
		return db
	}
	// 1. 获取dsn
	dsn := config.Config.MySQL.GetDSN()

	// 2. 连接MySQL
	cfg := &gorm.Config{
		DryRun: false,
	}
	gormDB, err := gorm.Open(mysql.Open(dsn), cfg)
	if err != nil {
		log.Println("连接数据库出错", err)
		return nil
	}
	db = gormDB
	return db
}
