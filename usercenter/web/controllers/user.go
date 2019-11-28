package controllers

import (
	"errors"

	"github.com/codelieche/microservice/usercenter/web/forms"
	"github.com/go-playground/validator"

	"github.com/codelieche/microservice/usercenter/datamodels"
	"github.com/codelieche/microservice/usercenter/web/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
)

type UserController struct {
	Session *sessions.Session
	Ctx     iris.Context
	Service services.UserService
}

//func (c *UserController) GetBy(id int64) (user *datamodels.User, success bool) {
//	if user, err := c.Service.GetById(id); err != nil {
//		return nil, false
//	} else {
//		return user, true
//	}
//}

func (c *UserController) GetBy(idOrName string) (user *datamodels.User, success bool) {
	if user, err := c.Service.GetByIdOrName(idOrName); err != nil {
		return nil, false
	} else {
		return user, true
	}
}

func (c *UserController) GetList(ctx iris.Context) (users []*datamodels.User, success bool) {
	return c.GetListBy(1, ctx)
}

// 获取用户的列表页：注意以前的版本，ctx是可放前面的
func (c *UserController) GetListBy(page int, ctx iris.Context) (users []*datamodels.User, success bool) {
	//	定义变量
	var (
		pageSize int
		offset   int
		limit    int
		err      error
	)

	// 获取变量
	pageSize = ctx.URLParamIntDefault("pageSize", 10)
	limit = pageSize
	if page > 1 {
		offset = (page - 1) * pageSize
	}

	//	获取用户
	//log.Println(offset, limit)
	if users, err = c.Service.List(offset, limit); err != nil {
		return nil, false
	} else {
		return users, true
	}
}

// 注册用户
func (c *UserController) PostCreate() (user *datamodels.User, err error) {
	//	定义变量
	var (
		username   = c.Ctx.FormValue("username")
		password   = c.Ctx.FormValue("password")
		repassword = c.Ctx.FormValue("repassword")
		mobile     = c.Ctx.FormValue("mobile")
		email      = c.Ctx.FormValue("email")
		userForm   forms.UserCreateForm
		v          *validator.Validate
	)

	if password != repassword {
		err = errors.New("密码和确认密码不一样")
		return nil, err
	}

	if len(username) < 6 {
		err = errors.New("用户名长度小于6")
		return nil, err
	}
	if len(password) < 8 {
		err = errors.New("用户密码长度小于8")
		return nil, err
	}

	// 验证表单
	v = validator.New()
	userForm = forms.UserCreateForm{
		Username:   username,
		Password:   password,
		Repassword: repassword,
		Mobile:     mobile,
		Email:      email,
	}
	if err = v.Struct(userForm); err != nil {
		return nil, err
	}

	user = &datamodels.User{
		Username:    username,
		Password:    password,
		Mobile:      mobile,
		Email:       email,
		Groups:      nil,
		Roles:       nil,
		Permissions: nil,
	}

	// 保存用户
	if user, err = c.Service.CreateUser(user); err != nil {
		return nil, err
	} else {
		// 创建用户成功
		return user, nil
	}
}
