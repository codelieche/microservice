package datasources

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/codelieche/microservice/usercenter/datamodels"

	"github.com/codelieche/microservice/usercenter/common"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var config *common.Config

func initDb() {
	var (
		err      error
		mysqlUri string
	)
	// 1. 先获取配置
	if config == nil {
		if err = common.ParseConfig(); err != nil {
			log.Println(err.Error())
			os.Exit(1)
		} else {
			config = common.GetConfig()
		}
	}
	//	2. 连接数据库
	// 2-1: 获取mysqlUri
	mysqlUri = fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Database.MySQL.User, config.Database.MySQL.Password,
		config.Database.MySQL.Host, config.Database.MySQL.Port, config.Database.MySQL.Database)
	//log.Println(mysqlUri)
	//	2-2: 连接数据库
	db, err = gorm.Open("mysql", mysqlUri)

	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	} else {

	}

	// 3. Migrate the Schema
	db.AutoMigrate(&datamodels.User{})
	db.AutoMigrate(&datamodels.Group{})
	db.AutoMigrate(&datamodels.Role{})
	db.AutoMigrate(&datamodels.Application{})
	db.AutoMigrate(&datamodels.Permission{})
	db.AutoMigrate(&datamodels.Ticket{})
	db.AutoMigrate(&datamodels.Token{})

	//	4. 是否显示Model的SQL
	db.LogMode(config.Debug)

	db.DB().SetConnMaxLifetime(100 * time.Second)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetMaxIdleConns(20)
}

func GetDb() *gorm.DB {
	if db != nil {
		return db
	} else {
		initDb()
		return db
	}
}
