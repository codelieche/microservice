package app

import (
	"log"
	"os"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12/sessions/sessiondb/redis"
)

var sess *sessions.Sessions
var redisDB *redis.Database

func initSession() {
	// 1. 连接Redis
	cfg := redis.Config{
		Network:   "tcp",
		Addr:      "127.0.0.1:6379",
		Clusters:  nil,
		Password:  "",
		Database:  "0",
		MaxActive: 10,
		Timeout:   time.Second * 20,
		Prefix:    "",
		Delim:     "-",
		Driver:    nil,
	}
	redisDB = redis.New(cfg)
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

	//	2. 实例化session
	sess = sessions.New(sessions.Config{
		Cookie:                      "sessionid",
		CookieSecureTLS:             false,
		AllowReclaim:                true,
		Encode:                      nil,
		Decode:                      nil,
		Encoding:                    nil,
		Expires:                     time.Minute * 60,
		SessionIDGenerator:          nil,
		DisableSubdomainPersistence: false,
		//Expires:                     time.Second * 10,
	})

	// 3. use database
	sess.UseDatabase(redisDB)
}

func useSessionMiddleware(app *iris.Application) {
	if sess == nil {
		initSession()
	}

	app.Use(sess.Handler())
}
