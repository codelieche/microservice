package main

import (
	"log"
	"time"

	"github.com/codelieche/microservice/usercenter/web/middlewares"

	"github.com/kataras/iris/v12/sessions"

	"github.com/kataras/iris/v12"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	app := iris.New()

	sess := sessions.New(sessions.Config{
		Cookie:       "sessionid",
		Expires:      0, // defaults to 0: unlimited life. Another good value is: 45 * time.Minute,
		AllowReclaim: true,
	})
	app.Use(sess.Handler())
	app.Use(middlewares.CheckTicketMiddleware, middlewares.CheckLoginMiddleware)
	//app.Use(middlewares.CheckLoginMiddleware)

	app.Get("/", func(ctx iris.Context) {
		sess.Start(ctx)
		ctx.JSON(iris.Map{
			"time": time.Now(),
		})
	})

	// 运行程序
	app.Run(iris.Addr(":9001"), iris.WithoutServerError(iris.ErrServerClosed))
}
