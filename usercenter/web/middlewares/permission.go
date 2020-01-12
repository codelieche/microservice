package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
)

/**
权限相关中间件
先实现功能，后续再优化，鉴权响应时间目标：10ms内
1. 权限检查函数
*/
func CheckUserPermission(ctx iris.Context, app string, code string) bool {
	// 1. 定义变量
	var (
		session      *sessions.Session
		ssoSessionID string
		tokenStr     string
	)

	//route := ctx.GetCurrentRoute()
	//log.Println(route)

	// 获取变量
	session = sessions.Get(ctx)
	if ctx.Request().Host != ssoServerHost {
		ssoSessionID = session.GetString("ssoSessionID")
	} else {
		ssoSessionID = session.ID()
	}

	if ssoSessionID == "" {
		// 判断是否有传递Token
		if tokenStr = ctx.Request().Header.Get("Authorization"); tokenStr != "" {

		} else {
			// 返回False
			return false
		}
	}

	url := fmt.Sprintf("http://%s/api/v1/user/permission/check?app=%s&code=%s",
		ssoServerHost, app, code)

	// 方式一：
	client := &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Second * 2,
	}

	// 发起请求
	if req, err := http.NewRequest("POST", url, nil); err != nil {
		log.Println(err)
		return false
	} else {
		req.Header.Add("Cookie", fmt.Sprintf("sessionid=%s", ssoSessionID))
		if tokenStr != "" {
			req.Header.Add("Authorization", tokenStr)
		}
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
			} else {
				return true
			}
		}
	}

}
