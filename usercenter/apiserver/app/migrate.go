package app

import (
	"github.com/codelieche/microservice/usercenter/core"
	"github.com/codelieche/microservice/usercenter/store"
	"log"
)

func (app *App) autoMigrate() {
	// 执行auto migrate
	if db := store.GetMySQLDB(app.Config); db != nil {
		core.AutoMigrate(db)
	} else {
		log.Panic("获取gorm.DB为空")
	}
}
