package datasources

import (
	"fmt"
	"log"
	"os"

	"github.com/codelieche/microservice/datamodels"

	"github.com/codelieche/microservice/common"
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
		config.Database.User, config.Database.Password,
		config.Database.Host, config.Database.Port, config.Database.Database)
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
}

func GetDb() *gorm.DB {
	if db != nil {
		return db
	} else {
		initDb()
		return db
	}
}
