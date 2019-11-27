package middlewares

import (
	"fmt"
	"log"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
)

// 检查是否是登录的中间件【不会跳转】
// 如果是登录了的，那么就执行下一步
// 未登录，就返回401错误
// 注意别把登录页面也设置了这个了，否则登录不了【会显示Unauthorized】
func IsAuthenticatedMiddleware(ctx iris.Context) {
	session := sessions.Get(ctx)
	if session.GetIntDefault("userID", 0) > 0 {
		ctx.Next()
	} else {
		urlPath := ctx.Request().URL.Path
		if urlPath == "/api/v1/user/auth" || urlPath == "/api/v1/user/login" || urlPath == "/user/login" {
			ctx.Next()
		} else {
			//log.Println("账号未登录")
			// 判断是否是：/api/v1/ticket/validate/
			if strings.HasPrefix(urlPath, "/api/v1/ticket/validate/") {
				ctx.Next()
			}
			//ctx.StatusCode(401)
			ctx.StatusCode(iris.StatusUnauthorized)

			// 也可让其跳转去登录页面
			//ctx.StatusCode(302)
			//ctx.Redirect("/user/login")
		}
	}
}

// 检查是否是登录的中间件【会跳转登录页】
// 如果是登录了的，那么就执行下一步
// 未登录，就返回401错误
// 注意别把登录页面也设置了这个了，否则登录不了【会显示Unauthorized】
func CheckLoginMiddleware(ctx iris.Context) {
	session := sessions.Get(ctx)
	if session.GetIntDefault("userID", 0) > 0 {
		ctx.Next()
	} else {
		urlPath := ctx.Request().URL.Path
		if urlPath == "/api/v1/user/auth" || urlPath == "/api/v1/user/login" || urlPath == "/user/login" {
			ctx.Next()
		} else {
			log.Println("账号未登录")
			//ctx.StatusCode(401)
			//ctx.StatusCode(iris.StatusUnauthorized)

			// 也可让其跳转去登录页面
			ctx.StatusCode(302)
			rediretUrl := fmt.Sprintf("http://localhost:9000/user/login?returnUrl=%s", "http://localhost:9001/")
			ctx.Redirect(rediretUrl)
		}
	}
}
