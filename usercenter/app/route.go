package app

import (
	"time"

	"github.com/codelieche/microservice/datasources"
	"github.com/codelieche/microservice/repositories"
	"github.com/codelieche/microservice/web/controllers"
	"github.com/codelieche/microservice/web/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func setAppRoute(app *iris.Application) {

	// 首页
	mvc.Configure(app.Party("/"), func(app *mvc.Application) {
		// session
		// 注册控制器需要的Session和StartTime
		app.Register(
			time.Now(),
		)
		app.Handle(new(controllers.IndexController))
	})

	// /api/v1 相关的api
	apiV1 := app.Party("/api/v1")

	// 用户相关api
	mvc.Configure(apiV1.Party("/user"), func(app *mvc.Application) {
		// 实例化User的Repository
		db := datasources.GetDb()
		repo := repositories.NewUserRepository(db)
		// 实例化User的Service
		uService := services.NewUserService(repo)
		// 注册Service
		app.Register(uService)
		// 添加Crontroller
		app.Handle(new(controllers.UserController))
	})

	// 分组相关api
	mvc.Configure(apiV1.Party("/group"), func(app *mvc.Application) {
		// 实例化User的Repository
		db := datasources.GetDb()
		repo := repositories.NewGroupRepository(db)
		// 实例化Group的Service
		gService := services.NewGroupService(repo)
		// 注册Service
		app.Register(gService)
		// 添加Crontroller
		app.Handle(new(controllers.GroupController))
	})

	// 角色相关api
	mvc.Configure(apiV1.Party("/role"), func(app *mvc.Application) {
		// 实例化User的Repository
		db := datasources.GetDb()
		repo := repositories.NewRoleRepository(db)
		// 实例化Group的Service
		gService := services.NewRoleService(repo)
		// 注册Service
		app.Register(gService)
		// 添加Crontroller
		app.Handle(new(controllers.RoleController))
	})

	// 权限相关api
	mvc.Configure(apiV1.Party("/permission"), func(app *mvc.Application) {
		// 实例化User的Repository
		db := datasources.GetDb()
		repo := repositories.NewPermissionRepository(db)
		// 实例化Service
		pService := services.NewPermissionService(repo)
		// 注册Service
		app.Register(pService)
		// 添加Crontroller
		app.Handle(new(controllers.PermissionController))
	})

}
