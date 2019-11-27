package app

import (
	"time"

	"github.com/codelieche/microservice/web/middlewares"

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
			sess.Start,
			time.Now(),
		)
		app.Handle(new(controllers.IndexController))
	})

	// User Page
	pageUser := app.Party("/user")
	// /user开头的页面都需要采用CheckLoginMiddleware中间件【如果未登录会跳转到登录页面】
	pageUser.Use(middlewares.CheckLoginMiddleware)
	mvc.Configure(pageUser, func(app *mvc.Application) {
		// 实例化User的Repository
		db := datasources.GetDb()
		repo := repositories.NewUserRepository(db)
		tRepo := repositories.NewTicketRepository(db)
		// 实例化User的Service
		uService := services.NewUserService(repo, tRepo)

		// 注册Service
		app.Register(uService, sess.Start)

		app.Handle(new(controllers.PageUserControler))
	})

	// /api/v1 相关的api
	apiV1 := app.Party("/api/v1")
	// /api/v1开头的url都需要使用IsAuthenticatedMiddleware的中间件
	apiV1.Use(middlewares.IsAuthenticatedMiddleware)

	// 用户相关api
	mvc.Configure(apiV1.Party("/user"), func(app *mvc.Application) {
		// 实例化User的Repository
		db := datasources.GetDb()
		repo := repositories.NewUserRepository(db)
		tRepo := repositories.NewTicketRepository(db)
		// 实例化User的Service
		uService := services.NewUserService(repo, tRepo)
		// 注册Service
		app.Register(uService, sess.Start)
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

	// Ticket相关api
	mvc.Configure(apiV1.Party("/ticket"), func(app *mvc.Application) {
		// 实例化User的Repository
		db := datasources.GetDb()
		repo := repositories.NewTicketRepository(db)
		// 实例化Service
		tService := services.NewTicketService(repo)
		// 注册Service
		app.Register(tService, sess.Start)
		// 添加Crontroller
		app.Handle(new(controllers.TicketController))
	})

	// 添加测试相关api
	mvc.Configure(app.Party("/demo"), func(app *mvc.Application) {
		// 注册
		app.Register(sess)

		// 添加Controller
		app.Handle(new(controllers.DemoController))
	})

}
