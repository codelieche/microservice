package middlewares

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/codelieche/microservice/usercenter/web/forms"

	"github.com/codelieche/microservice/usercenter/datamodels"
	"github.com/codelieche/microservice/usercenter/datasources"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
)

var ssoServerHost string = "sso.codelieche.com"

func SetSsoServerHost(host string) {
	ssoServerHost = host
}

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
	var needReturnUrl bool

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
		needReturnUrl = true
	}

	// 跳转登录页面
	if needReturnUrl {
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
	var needReturnUrl bool
	if userID > 0 {
		// 调用sso server的：/api/v1/user/auth接口
		// 判断是否需要验证用户在sso server是否登录:
		if ctx.Host() != ssoServerHost {
			// 当前中间件请求的服务host不是sso server的host
			go checkSessionIsOkRsync(session, userID)
			ctx.Next()

			// 后续优化把检查ssoSession改成异步的
			// 同步方式检查session
			//ssoSessionID := session.GetString("ssoSessionID")
			//if success := CheckSessionIsOK(ssoSessionID, userID); success {
			//	ctx.Next()
			//} else {
			//	// 未登录
			//	//log.Println("未登录")
			//	// 摧毁session
			//	session.Destroy()
			//	needReturnUrl = true
			//}
		} else {
			// 不需要校验: 是sso的系统
			ctx.Next()
		}

	} else {
		needReturnUrl = true
	}
	// 跳转登录页面
	if needReturnUrl {
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
			rediretUrl := fmt.Sprintf("http://%s/user/login?returnUrl=%s", ssoServerHost, url)
			ctx.Redirect(rediretUrl)
		}
	}

}

func CheckSessionIsOK(ssoSessionID string, userID int) bool {
	// 后续应该用rpc检查

	// 发起http请求
	//ssoServerHost := "localhost:9000"
	userAuthUrl := fmt.Sprintf("http://%s/api/v1/user/auth", ssoServerHost)

	// 方式一：
	client := &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Second * 2,
	}

	if req, err := http.NewRequest("GET", userAuthUrl, nil); err != nil {
		log.Println(err)
		return false
	} else {
		req.Header.Add("Cookie", fmt.Sprintf("sessionid=%s", ssoSessionID))
		//req.Header.Add("Host", "sso.codelieche.com")

		if resp, err := client.Do(req); err != nil {
			log.Println(err)
			return false
		} else {
			//log.Println(resp.StatusCode, sessionID)
			//d, e := ioutil.ReadAll(resp.Body)
			//log.Println(string(d), e)
			if resp.StatusCode != 200 {
				return false
			}
			defer resp.Body.Close()
			if body, err := ioutil.ReadAll(resp.Body); err != nil {
				log.Println(err)
				return false
			} else {
				//log.Println(string(body))
				result := forms.TicketValidateUser{}
				if err = json.Unmarshal(body, &result); err != nil {
					return false
				} else {
					if result.ID == uint(userID) {
						return true
					} else {
						log.Println(result)
						return false
					}
				}
			}
		}
	}
}

// 异步检查session是否ok
// 请用go checkSessionIsOkRsync(session, userID)方式运行
// 这样假如sso已经退出了，但是当前系统未退出，这次检查出未登录，销毁当前系统的session
// 下一次登录的时候，就会跳转到登录页面了
func checkSessionIsOkRsync(session *sessions.Session, userID int) {
	ssoSessionID := session.GetString("ssoSessionID")
	//time.Sleep(time.Second * 1)
	if ssoSessionID != "" {
		if CheckSessionIsOK(ssoSessionID, userID) {
			//log.Println("检查session成功")
			// 成功
		} else {
			// 不成功，销毁当前服务的seesion
			msg := fmt.Sprintf("检查sso sesseion:%s已经退出登录, 需要销毁当前系统的session:%s", ssoSessionID, session.ID())
			log.Println(msg)
			session.Destroy()
		}
	}
}
