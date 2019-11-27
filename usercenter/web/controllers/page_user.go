package controllers

import (
	"crypto/md5"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/codelieche/microservice/common"

	"github.com/codelieche/microservice/datamodels"
	"github.com/codelieche/microservice/web/forms"
	"github.com/go-playground/validator"
	"github.com/kataras/iris/v12/mvc"

	"github.com/codelieche/microservice/web/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
)

// 用户登录、登出页面
// GET：/user/login 用户登录页面
// POST：/user/login POST登录
// GET: /user/logout 退出登录
type PageUserControler struct {
	Session *sessions.Session
	Ctx     iris.Context
	Service services.UserService
}

// 判断用户是否登录
func (c *PageUserControler) isLoginIn() bool {
	userID := c.Session.GetIntDefault("userID", 0)
	return userID > 0
}

// 获取当前session中的用户
func (c *PageUserControler) getCurrentUser() (user *datamodels.User, err error) {
	userID := c.Session.GetIntDefault("userID", 0)
	if userID > 0 {
		if user, err = c.Service.GetById(int64(userID)); err != nil {
			return nil, err
		} else {
			return user, nil
		}
	} else {
		// ID是0，等同用户不存在
		return nil, common.NotFountError
	}
}

// 检查要跳转的网址是否可以跳转
func (c *PageUserControler) checkReturnUrl(returnUrl string) bool {
	if u, err := url.Parse(returnUrl); err != nil {
		log.Println(err)
		return false
	} else {
		// 判断url的host
		host := u.Host
		config := common.GetConfig()
		for _, domain := range config.Http.Domains {
			//log.Println(domain)
			if domain != "" {
				if strings.HasSuffix(host, domain) {
					return true
				} else {
					// 继续检查下一个域名
					continue
				}
			}
		}
		// 循环完毕，返回false
		// 如果domains为空，就表示不可跳转域名
		//return strings.HasSuffix(host, "codelieche.com")
		return false
	}
}

func (c *PageUserControler) generateTicket(returnUrl string) *datamodels.Ticket {

	Md5Inst := md5.New()

	Md5Inst.Write([]byte(fmt.Sprintf("%s-%d", c.Session.ID(), time.Now().UnixNano())))

	Result := Md5Inst.Sum([]byte(""))
	ticketName := fmt.Sprintf("%x", Result)
	ticket := datamodels.Ticket{
		Name:        ticketName,
		Session:     c.Session.ID(),
		ReturnUrl:   returnUrl,
		TimeExpired: time.Now().Add(time.Minute),
	}

	return &ticket
}

func (c *PageUserControler) GetLogin() {

	// 生成ticket

	//session := sessions.Get(c.Ctx)
	userID := c.Session.GetIntDefault("userID", 0)
	//log.Println(userID)

	if user, err := c.getCurrentUser(); err != nil {
		// 获取用户失败
		//log.Println(err)

	} else {
		// 判断是否需要跳转
		if user != nil && user.ID > 0 {
			// 获取returnUrl
			returnUrl := c.Ctx.URLParam("returnUrl")
			//returnUrl := c.Ctx.URLParamDefault("returnUrl", "http://www.codelieche.com")

			if returnUrl != "" && c.checkReturnUrl(returnUrl) {
				// 可以跳转:Ticket可以生成一个另外的值【推荐】
				ticket := c.generateTicket(returnUrl)
				//log.Println(ticket)

				if strings.Contains(returnUrl, "?") {
					returnUrl = fmt.Sprintf("%s&ticket=%s", returnUrl, ticket.Name)
				} else {
					returnUrl = fmt.Sprintf("%s?ticket=%s", returnUrl, ticket.Name)
				}
				c.Ctx.Redirect(returnUrl)
			} else {
				//	填写提示信息
				msg := fmt.Sprintf("已经登录，用户名是：%s", user.Username)
				c.Ctx.ViewData("msg", msg)
				c.Ctx.ViewData("isLogin", true)
			}
		}
	}

	if userID > 0 {
		// 用户登录了的
	} else {
		// 未登录，显示登录页面
	}

	//c.Ctx.ViewData("msg", "提示消息")
	c.Ctx.View("user/login.html")
}

func (c *PageUserControler) GetLogout() {
	c.Session.Destroy()
	c.Ctx.Redirect("/user/login")
}

func (c *PageUserControler) PostLogin() mvc.Result {
	var (
		username  = c.Ctx.FormValue("username")
		password  = c.Ctx.FormValue("password")
		mobile    = c.Ctx.FormValue("mobile")
		email     = c.Ctx.FormValue("email")
		returnUrl = c.Ctx.URLParam("returnUrl")
		userForm  forms.UserLoginForm
		v         *validator.Validate
		user      *datamodels.User
		success   bool
		err       error
	)
	//log.Println(c.Ctx.Request().URL)

	// 如果是登录的，那么就退出登录先
	if c.isLoginIn() {
		//c.logout()
		c.Session.Destroy()
	}

	// 判断用户名和密码
	if username == "" {
		err = errors.New("username不可为空")
		goto ERR
	}

	// 判断用户名和密码
	if password == "" {
		err = errors.New("password不可为空")
		goto ERR
	}

	// 验证表单
	v = validator.New()
	userForm = forms.UserLoginForm{
		Username: username,
		Password: password,
		Mobile:   mobile,
		Email:    email,
	}
	if err = v.Struct(userForm); err != nil {
		goto ERR
	}

	// 获取用户
	user, err = c.Service.GetByIdOrName(username)
	if err != nil {
		goto ERR
	}

	// 判断用户密码是否正确
	if success, err = user.CheckPassword(password); err != nil {
		err = errors.New("输入的密码不正确")
		goto ERR
	} else {
		if success {
			// 登录成功
			c.Session.Set("userID", user.ID)

			// 判断是否需要跳转
			if returnUrl != "" && c.checkReturnUrl(returnUrl) {
				//log.Println(returnUrl)

				//c.Ctx.StatusCode(302)
				//c.Ctx.Redirect(returnUrl)
				return mvc.Response{
					Path: returnUrl,
				}
			} else {
				c.Ctx.JSON(user)
			}
		} else {
			err = errors.New("账号或者密码不正确")
			goto ERR
		}
	}
ERR:
	//c.Ctx.JSON(iris.Map{
	//	"error": err.Error(),
	//})
	return mvc.Response{
		Err: err,
	}
}

func (c *PageUserControler) GetInfo() string {
	return "This Is User Info Page"
}
