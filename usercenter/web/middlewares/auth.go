package middlewares

import (
	"fmt"
	"strings"

	"github.com/codelieche/microservice/usercenter/datamodels"
	"github.com/codelieche/microservice/usercenter/datasources"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
)

// 给ctx设置user
func CtxSetUserMiddleware(ctx iris.Context) {
	session := sessions.Get(ctx)
	// 判断是否有userID，而且用户是有效的
	userID := session.GetIntDefault("userID", 0)
	if userID > 0 {
		// 根据ID获取用户，判断用户是否存在
		db := datasources.GetDb()
		user := datamodels.User{}
		if db.Model(datamodels.User{}).First(&user, "id = ?", userID).Error != nil {
			// 获取用户出错
			//log.Println("获取用户出错")
			session.Destroy()
			// 无需设置
			goto NEXT
		} else {
			// 判断用户是否是禁用了，如果禁用了，那么也销毁session
			if !user.IsActive {
				session.Destroy()
			} else {
				ctx.Values().Set("user", user)
			}
			goto NEXT
		}
	} else {
		// 无需设置
	}
NEXT:
	ctx.Next()
}

// 检查是否是登录的中间件【不会跳转】
// 如果是登录了的，那么就执行下一步
// 未登录，就返回401错误
// 注意别把登录页面也设置了这个了，否则登录不了【会显示Unauthorized】
func IsAuthenticatedMiddleware(ctx iris.Context) {
	session := sessions.Get(ctx)

	// 判断是否有userID，而且用户是有效的
	userID := session.GetIntDefault("userID", 0)

	// 如果没有user，那么就表示当前用户未登录
	u := ctx.Values().Get("user")
	var user datamodels.User
	if u != nil {
		user = u.(datamodels.User)
	}
	// log.Println(user)
	if userID > 0 && user.IsActive {
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
	// 获取session中用户的ID
	userID := session.GetIntDefault("userID", 0)

	//log.Println(userID)
	if userID > 0 {
		ctx.Next()
	} else {
		//log.Println("未登录的，需要跳转sso系统")
		urlPath := ctx.Request().URL.Path
		if urlPath == "/api/v1/user/auth" || urlPath == "/api/v1/user/login" || urlPath == "/user/login" {
			ctx.Next()
		} else {
			// 获取到当前的URL
			//log.Println(requestUrl)
			//log.Println(ctx.Request().Host, ctx.Host())
			//log.Println(ctx.Request().Proto)

			var url string

			requestUrl := ctx.Request().RequestURI
			if strings.Contains(ctx.Request().Proto, "HTTPS") {
				url = fmt.Sprintf("https://%s%s", ctx.Host(), requestUrl)
			} else {
				url = fmt.Sprintf("http://%s%s", ctx.Host(), requestUrl)
			}

			//log.Println("账号未登录")
			//ctx.StatusCode(401)
			//ctx.StatusCode(iris.StatusUnauthorized)

			// 也可让其跳转去登录页面
			ctx.StatusCode(302)
			rediretUrl := fmt.Sprintf("http://localhost:9000/user/login?returnUrl=%s", url)
			ctx.Redirect(rediretUrl)
		}
	}
}
