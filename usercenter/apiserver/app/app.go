package app

import (
	"fmt"
	"github.com/codelieche/microservice/usercenter/apiserver/router"
	"github.com/codelieche/microservice/usercenter/internal/config"
	"github.com/gin-gonic/gin"
)

type App struct {
	Config *config.GlobalConfig
	Engine *gin.Engine
}

func NewApp() *App {
	cfg := config.GetConfig()
	return &App{Config: cfg}
}

func (app *App) Run() {
	// 1. 实例化 gin Engine
	engin := gin.New()
	app.Engine = engin

	// 2. Auto Migrate
	app.autoMigrate()

	// 3. 配置路由
	router.InjectRouter(app.Engine)

	// 4. 启动web服务
	webAddr := fmt.Sprintf("%s:%d", app.Config.Web.Address, app.Config.Web.Port)
	app.Engine.Run(webAddr)
}
