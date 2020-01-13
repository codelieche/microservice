package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/codelieche/microservice/usercenter/common"

	"github.com/codelieche/microservice/usercenter/datasources"
	"github.com/codelieche/microservice/usercenter/repositories"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
)

var userRepo = repositories.NewUserRepository(datasources.GetDb())

func CheckUserPermissionLocal(ctx iris.Context, codes ...string) bool {
	session := sessions.Get(ctx)
	userID := session.GetInt64Default("userID", 0)

	// 获取用户的权限缓存
	if permissionsMap, err := userRepo.GetOrSetUserPermissionsCache(userID, false); err != nil {
		ctx.StatusCode(403)
		return false
	} else {
		// 判断权限是否在map中
		for _, code := range codes {
			checkKey := fmt.Sprintf("app_%d_%s", common.GetConfig().App.ID, code)
			if _, isExist := permissionsMap[checkKey]; isExist {
				//ctx.StatusCode(200)
				// codes列表全部为true，才可以
			} else {
				ctx.StatusCode(403)
				return false
			}
		}
		// 到这里，返回true
		return true
	}
}

/**
权限相关中间件
先实现功能，后续再优化，鉴权响应时间目标：10ms内
1. 权限检查函数
*/
func PostCheckUserPermission(ctx iris.Context, app string, codes ...string) bool {
	// 1. 定义变量
	var (
		session      *sessions.Session
		codeListStr  string
		ssoSessionID string
		tokenStr     string
	)

	//route := ctx.GetCurrentRoute()
	//log.Println(route)

	// 获取变量
	session = sessions.Get(ctx)
	log.Println(ctx.Request().Cookies())

	if ctx.Request().Host != ssoServerHost {
		ssoSessionID = session.GetString("ssoSessionID")
		// 如果未找到ssoSessionID的值，就实用当前session的ID
		if ssoSessionID == "" {
			ssoSessionID = session.ID()
		}
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
	codeListStr = strings.Join(codes, "&codes=")
	url := fmt.Sprintf("http://%s/api/v1/user/permission/check?app=%s&codes=%s",
		ssoServerHost, app, codeListStr)

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
