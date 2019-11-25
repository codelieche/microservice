package controllers

import (
	"errors"

	"github.com/codelieche/microservice/web/forms"
	"github.com/go-playground/validator"

	"github.com/kataras/iris/v12/mvc"

	"github.com/codelieche/microservice/datamodels"
)

// 用户登录登出相关

// 判断用户是否登录
func (c *UserController) isLoginIn() bool {
	userID := c.Session.GetIntDefault("userID", 0)
	return userID > 0
}

// Post登录用户
func (c *UserController) PostLogin() (user *datamodels.User, err error) {
	var (
		username = c.Ctx.FormValue("username")
		password = c.Ctx.FormValue("password")
		mobile   = c.Ctx.FormValue("mobile")
		email    = c.Ctx.FormValue("email")
		userForm forms.UserLoginForm
		v        *validator.Validate
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
		return nil, err
	}

	// 判断用户名和密码
	if password == "" {
		err = errors.New("password不可为空")
		return nil, err
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
		return nil, err
	}

	// 获取用户
	u, err := c.Service.GetByIdOrName(username)
	if err != nil {
		return nil, err
	}

	// 判断用户密码是否正确
	if success, err := u.CheckPassword(password); err != nil {
		err = errors.New("输入的密码不正确")
		return nil, err
	} else {
		if success {
			// 登录成功
			c.Session.Set("userID", u.ID)

			return u, nil
		} else {
			err = errors.New("账号或者密码不正确")
			return nil, err
		}
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
		return nil, nil
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
