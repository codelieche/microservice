package controllers

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/codelieche/microservice/usercenter/web/forms"
	"github.com/go-playground/validator"

	"github.com/codelieche/microservice/usercenter/datamodels"
	"github.com/codelieche/microservice/usercenter/web/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
)

type UserController struct {
	Session           *sessions.Session
	Ctx               iris.Context
	Service           services.UserService
	GroupService      services.GroupService
	RoleService       services.RoleService
	PermissionService services.PermissionService
}

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
	// 判断session是否登录
	if c.Session.GetIntDefault("userID", 0) != 0 {
		err = errors.New("当前登录了用户，不可创建用户")
		return nil, err
	}
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
		IsActive:    true,
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

// 用户更新
func (c *UserController) PutBy(id int64, ctx iris.Context) (user *datamodels.User, err error) {
	// 1. 定义变量
	var (
		contentType string
		form        *forms.UserUpdateform
		//updateFields map[string]interface{}
		groups      []*datamodels.Group
		group       *datamodels.Group
		roles       []*datamodels.Role
		role        *datamodels.Role
		permissions []*datamodels.Permission
		permission  *datamodels.Permission
	)

	// 2. 获取数据
	// 2-1: 获取用户
	if user, err = c.Service.GetById(id); err != nil {
		return nil, err
	}

	// 2-2: 获取Content-Type
	contentType = ctx.Request().Header.Get("Content-Type")

	// 2-3：解析表单
	form = &forms.UserUpdateform{}
	if strings.Contains(contentType, "application/json") {
		err = ctx.ReadJSON(form)
	} else {
		err = ctx.ReadForm(form)
	}

	// 2-4：判断获取form是否成功
	if err != nil {
		return nil, err
	} else {
		//log.Println(form)
	}

	// 3. 处理更新字段
	//updateFields = make(map[string]interface{})
	//updateFields["email"] = form.Email
	//updateFields["mobile"] = form.Mobile
	user.Email = form.Email
	user.Mobile = form.Mobile

	// 处理many2many字段
	groups = []*datamodels.Group{}
	if len(form.Groups) > 0 {
		for _, i := range form.Groups {
			if group, err = c.GroupService.GetById(i); err != nil {
				err = fmt.Errorf("分组(ID:%d)：%s", i, err)
				return nil, err
			} else {
				groups = append(groups, group)
			}
		}
		//updateFields["groups"] = groups
	}
	user.Groups = groups

	roles = []*datamodels.Role{}
	if len(form.Roles) > 0 {
		for _, i := range form.Roles {
			if role, err = c.RoleService.GetById(i); err != nil {
				err = fmt.Errorf("角色(ID:%d)：%s", i, err)
				return nil, err
			} else {
				roles = append(roles, role)
			}
		}
	}
	user.Roles = roles

	permissions = []*datamodels.Permission{}
	if len(form.Permissions) > 0 {
		for _, i := range form.Permissions {
			if permission, err = c.PermissionService.GetById(i); err != nil {
				err = fmt.Errorf("权限(ID:%d)：%s", i, err)
				return nil, err
			} else {
				permissions = append(permissions, permission)
			}
		}
	}
	user.Permissions = permissions

	// 4. 更新字段
	return c.Service.Save(user)
}

func (c *UserController) DeleteBy(idOrName string) {
	// 删掉用户的session
	if success, err := c.Service.DeleteUserByIdOrName(idOrName); err != nil {
		log.Println(err)
		c.Ctx.StatusCode(iris.StatusBadRequest)
	} else {
		if success {
			c.Ctx.StatusCode(iris.StatusNoContent)
		} else {
			c.Ctx.StatusCode(iris.StatusBadRequest)
		}
	}
}
