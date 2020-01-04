package app

import (
	"fmt"
	"log"

	"github.com/codelieche/microservice/usercenter/web/middlewares"

	"github.com/codelieche/microservice/usercenter/common"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func newApp() *iris.Application {
	app := iris.New()

	// 配置应用
	appConfigure(app)

	// 设置auth
	//appAddBasictAuth(app)

	// 使用中间件，添加logger
	app.Use(logger.New(logger.Config{
		Status:             true,
		IP:                 true,
		Method:             true,
		Path:               true,
		Query:              true,
		Columns:            false,
		MessageContextKeys: nil,
		MessageHeaderKeys:  nil,
		LogFunc:            nil,
		LogFuncCtx:         nil,
		Skippers:           nil,
	}))
	useSessionMiddleware(app)
	//app.Use(checkLogin)
	//app.Use(checkLogin)

	// 给context设置User
	// 设置下sso的域名
	middlewares.SetSsoServerHost("0.0.0.0:9000")
	app.Use(middlewares.CtxSetUserMiddleware)

	// 处理错误页面
	handleAppOnError(app)

	// 设置View的路径
	viewEngine := iris.HTML("./web/templates", ".html")
	app.RegisterView(viewEngine)

	// 当执行kill的时候执行操作：关闭数据库等
	iris.RegisterOnInterrupt(handleAppInterupt)

	// 设置路由：重点
	setAppRoute(app)

	// 静态文件
	app.HandleDir("/static", "./web/public")

	return app
}

func Run() {
	app := newApp()

	config := common.GetConfig()
	addr := fmt.Sprintf("%s:%d", config.Http.Host, config.Http.Port)

	// 运行程序
	app.Run(iris.Addr(addr), iris.WithoutServerError(iris.ErrServerClosed))
}
