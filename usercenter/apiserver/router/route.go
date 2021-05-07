package router

import (
	"github.com/codelieche/microservice/usercenter/internal/config"
	"github.com/codelieche/microservice/usercenter/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func InjectRouter(engin *gin.Engine) {
	// 1. 显示首页提示，只提供web api服务
	engin.Handle("GET", "/", func(context *gin.Context) {
		context.String(200, "Hello World！")
	})

	// 2. 健康检查的接口

	// 3. api相关的接口
	api := engin.Group("/api/v1")
	// 加入中间件
	api.Use(middlewares.JwtAuth(config.Config.Web.JWT))
	injectUserRoute(api)

	api.Handle("GET", "/", func(context *gin.Context) {
		context.String(200, "Hello World!")
	})
}
