package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/codelieche/microservice/usercenter/common"

	"github.com/codelieche/microservice/usercenter/web/forms"
	"github.com/go-playground/validator"

	"github.com/kataras/iris/v12/mvc"

	"github.com/codelieche/microservice/usercenter/datamodels"
)

// 用户登录登出相关

// 判断用户是否登录
func (c *UserController) isLoginIn() bool {
	userID := c.Session.GetIntDefault("userID", 0)
	return userID > 0
}

// 检查要跳转的网址是否可以跳转
func (c *UserController) checkReturnUrl(returnUrl string) bool {
	if u, err := url.Parse(returnUrl); err != nil {
		log.Println(err)
		return false
	} else {
		// 判断url的host
		host := u.Host
		config := common.GetConfig()
		for _, domain := range config.Http.Domains {
			log.Println(domain)
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

// Get登录页面
func (c *UserController) GetLogin() {
	if c.isLoginIn() {
		if u, err := c.getCurrentUser(); err != nil || u == nil {
			//	用户未登录
		} else {
			// 判断是否需要跳转
			// 获取returnUrl
			returnUrl := c.Ctx.URLParam("returnUrl")
			if returnUrl != "" && c.checkReturnUrl(returnUrl) {

				c.Ctx.Redirect(returnUrl)
			} else {
				// 填写信息
				msg := fmt.Sprintf("已经登录，用户名是：%s", u.Username)
				c.Ctx.ViewData("msg", msg)
				c.Ctx.ViewData("isLogin", true)
			}
		}
	}

	//c.Ctx.ViewData("msg", "提示消息")
	c.Ctx.View("user/login.html")
}

// Post登录用户
//func (c *UserController) PostLogin() (user *datamodels.User, err error) {
func (c *UserController) PostLogin() mvc.Result {
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

// 判断用户是否登录
func (c *UserController) GetAuth() string {
	if c.isLoginIn() {
		return "logined"
	} else {
		return "no login"
	}
}

// 退出登录
func (c *UserController) GetLogout() mvc.Result {
	c.Session.Destroy()
	return mvc.Response{
		Path: "/api/v1/user/auth",
	}
}

// 获取当前session中的用户
func (c *UserController) getCurrentUser() (user *datamodels.User, err error) {
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

// 获取当前登录的用户信息
func (c *UserController) GetInfo() (user *datamodels.User, err error) {
	userID := c.Session.GetIntDefault("userID", 0)
	if userID > 0 {
		if user, err = c.Service.GetById(int64(userID)); err != nil {
			return nil, err
		} else {
			return user, nil
		}
	} else {
		return nil, nil
	}
}

// 修改用户密码
func (c *UserController) PostChangePassword() mvc.Result {
	if !c.isLoginIn() {
		return mvc.Response{
			Path: "/api/v1/user/auth",
		}
	}
	// 获取信息
	var (
		password   = c.Ctx.FormValue("password")
		repassword = c.Ctx.FormValue("repassword")
		userForm   forms.UserChangePasswrodForm
		user       *datamodels.User
		err        error
		v          *validator.Validate
	)

	// 验证表单
	userForm = forms.UserChangePasswrodForm{
		Password:   password,
		Repassword: repassword,
	}
	v = validator.New()

	if err = v.Struct(userForm); err != nil {
		return mvc.Response{
			Err: err,
		}
	}

	if password != repassword {
		err = errors.New("密码和确认密码不一样")
		return mvc.Response{
			Err: err,
		}
	}

	// 获取用户
	if user, err = c.getCurrentUser(); err != nil {
		return mvc.Response{
			Err: err,
		}
	} else {
		// 修改密码
		if user, err = c.Service.ChangeUserPassword(user, password); err != nil {
			return mvc.Response{
				Err: err,
			}
		} else {
			return mvc.Response{
				Object: user,
			}
		}
	}
}
