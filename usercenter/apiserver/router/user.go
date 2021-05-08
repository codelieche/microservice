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
	// 校验Token
	api.GET("/user/auth/", userController.Auth)
	// 创建用户
	api.POST("/user/", userController.Create)
	// 用户列表
	api.GET("/user/", userController.List)
	// 用户信息
	api.GET("/user/:id/info/", userController.Find)
	// 修改密码
	api.POST("/user/password/change/", userController.ChangePassword)
	// 重置密码
	api.POST("/user/password/reset/", userController.ResetPassword)
}
