package app

import (
	"time"

	"github.com/kataras/iris/v12/sessions"

	"github.com/codelieche/microservice/common"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/basicauth"
)

// 设置Basic Auth仅供测试使用
func appAddBasictAuth(app *iris.Application) {
	config := common.GetConfig()

	authConfig := basicauth.Config{
		//Users:   map[string]string{"user01": "password01", "user02": "password02"},
		Users:   config.Http.BasicAuth,
		Realm:   "Atuhorization Required",
		Expires: time.Duration(10) * time.Second,
		OnAsk:   nil,
	}

	authentication := basicauth.New(authConfig)

	//设置authentication
	app.Use(authentication)
}

func checkLogin(ctx iris.Context) {
	session := sessions.Get(ctx)
	// 判断是否登录了
	if session.GetIntDefault("userID", 0) > 0 {
		ctx.Next()
	} else {
		urlPath := ctx.Request().URL.Path
		if urlPath == "/api/v1/user/auth" || urlPath == "/api/v1/user/login" {
			ctx.Next()
		} else {
			ctx.StatusCode(401)
			//ctx.Redirect("/api/v1/user/auth")
		}
	}
}
