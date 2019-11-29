package middlewares

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/kataras/iris/v12/sessions"

	"github.com/codelieche/microservice/usercenter/web/forms"

	"github.com/kataras/iris/v12"
)

// 检查Ticket的中间件
func CheckTicketMiddleware(ctx iris.Context) {
	// 如果URL中传递了Ticket参数: 就向sso检查ticket
	ticket := ctx.URLParam("ticket")
	// 所以所有的接口，都不再使用ticket这个名字传参了
	if ticket != "" {
		// 向sso发起校验ticket
		// SSO_SERVER_HOST := "sso.codelieche.com"
		// 传入ticket检查：如果是登录了的，那么设置本系统为登录，且获取到用户信息
		// RPC检测最好，但是先用http方式校验
		if result, err := CheckTicketFromSSOServer(ticket); err != nil {
			log.Println(err.Error())
			ctx.Text("验证ticket出错")
		} else {
			//log.Println(sessionID)
			session := sessions.Get(ctx)
			userID := fmt.Sprintf("%d", result.User.ID)
			//log.Println(result.User.ID, result.User.Username, userID)
			session.Set("userID", userID)
			session.Set("ssoSessionID", result.Session)
			session.Set("username", result.User.Username)

			ctx.SetCookie(&http.Cookie{
				Name:    "username",
				Value:   result.User.Username,
				Expires: time.Now().Add(time.Hour),
			})

			ctx.SetCookie(&http.Cookie{
				Name:    "userID",
				Value:   userID,
				Expires: time.Now().Add(time.Hour),
			})

			// 把user设置到context中
			//ctx.Values().Set("user", result.User)

			//ctx.SetCookie(&http.Cookie{
			//	Name:    "ssoSessionID",
			//	Value:   result.Session,
			//	Expires: time.Now().Add(time.Minute),
			//})
			//log.Println("设置cookie完毕")

			// 设置系统的cookie，然后跳转Url
			url := ctx.Request().URL
			urlStr := fmt.Sprintf("%s", url)
			urlStr = strings.Split(urlStr, "?ticket=")[0]
			//log.Println(urlStr)
			ctx.StatusCode(iris.StatusFound)
			ctx.Redirect(urlStr)
		}
	} else {
		ctx.Next()
	}
}

// 向sso server发起检查ticket的操作
func CheckTicketFromSSOServer(ticket string) (ticketResponse *forms.TicketValidateResponse, err error) {
	// 检查2次
	ssoServerHost := "localhost:9000"
	ticketValidateUrl := fmt.Sprintf("http://%s/api/v1/ticket/validate/%s", ssoServerHost, ticket)

	// 方式一：
	client := &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Second * 5,
	}

	if req, err := http.NewRequest("GET", ticketValidateUrl, nil); err != nil {
		log.Println(err)
		return nil, err
	} else {
		//req.Header.Add("Host", "www.codelieche.com")
		//req.SetBasicAuth("user01", "password01")

		if resp, err := client.Do(req); err != nil {
			log.Println(err)
			return nil, err
		} else {
			defer resp.Body.Close()
			if body, err := ioutil.ReadAll(resp.Body); err != nil {
				log.Println(err)
				return nil, err
			} else {
				//log.Println(string(body))
				result := forms.TicketValidateResponse{}
				if err = json.Unmarshal(body, &result); err != nil {
					return nil, err
				} else {
					if result.Session != "" {
						return &result, nil
					} else {
						log.Println(result)
						return nil, errors.New("session为空")
					}
				}
			}
		}
	}

	// 方式二：不设置超时时间还是不ok的
	//if resp, err := http.Get(ticketValidateUrl); err != nil {
	//	log.Println(err)
	//} else {
	//	defer resp.Body.Close()
	//	if body, err := ioutil.ReadAll(resp.Body); err != nil {
	//		log.Println(err)
	//	} else {
	//		log.Println(string(body))
	//	}
	//}
}
