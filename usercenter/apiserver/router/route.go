package router

import "github.com/gin-gonic/gin"

func InjectRouter(engin *gin.Engine) {
	// 1. 显示首页提示，只提供web api服务
	engin.Handle("GET", "/", func(context *gin.Context) {
		context.String(200, "Hello World！")
	})

	// 2. 健康检查的接口

	// 3. api相关的接口
	api := engin.Group("/api/v1")
	injectUserRoute(api)

	api.Handle("GET", "/", func(context *gin.Context) {
		context.String(200, "Hello World!")
	})
}
