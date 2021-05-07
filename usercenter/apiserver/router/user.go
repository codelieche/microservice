package router

import (
	"github.com/codelieche/microservice/usercenter/controllers"
	"github.com/codelieche/microservice/usercenter/services"
	"github.com/codelieche/microservice/usercenter/store"
	"github.com/gin-gonic/gin"
)

func injectUserRoute(api *gin.RouterGroup) {
	db := store.GetMySQLDB(nil)
	s := store.NewUserStore(db)
	userService := services.NewUserService(s)
	userController := controllers.NewUserController(userService)

	// 创建用户
	api.POST("/user/login/", userController.Login)
	api.GET("/user/auth/", userController.Auth)
	api.POST("/user/", userController.Create)
	api.GET("/user/:id/info/", userController.Find)
}
