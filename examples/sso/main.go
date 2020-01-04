package main

import (
	"log"
	"os"
	"time"

	"github.com/codelieche/microservice/middleware"

	"github.com/kataras/iris/v12/middleware/logger"

	"github.com/kataras/iris/v12/sessions/sessiondb/redis"

	"github.com/kataras/iris/v12/sessions"

	"github.com/kataras/iris/v12"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	app := iris.New()

	app.Use(logger.New())

	// 1. 连接Redis
	cfg := redis.Config{
		Network:   "tcp",
		Addr:      "127.0.0.1:6379",
		Clusters:  nil,
		Password:  "",
		Database:  "1",
		MaxActive: 10,
		Timeout:   time.Second * 20,
		Prefix:    "",
		Delim:     "-",
		Driver:    nil,
	}
	redisDB := redis.New(cfg)
	// 检查redis是否连接ok

	defer func() {
		if r := recover(); r != nil {
			log.Println("捕获到错误！连接Redis出错", r)
			os.Exit(1)
		}
	}()

	// 获取redis中：d03b09dd-f1f0-456e-945f-fd0588a577f6-test的值
	// 如果出错，会被recover捕获到，说明redis没起来
	redisDB.Get("d03b09dd-f1f0-456e-945f-fd0588a577f6", "test")

	sess := sessions.New(sessions.Config{
		Cookie:       "sessionid",
		Expires:      time.Second * 100, // defaults to 0: unlimited life. Another good value is: 45 * time.Minute,
		AllowReclaim: true,
	})

	sess.UseDatabase(redisDB)
	app.Use(sess.Handler())
	middleware.SetSsoServerHost("localhost:9000")
	app.Use(middleware.CheckTicketMiddleware, middleware.CheckLoginMiddleware)
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
