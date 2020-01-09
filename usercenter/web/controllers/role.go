package controllers

import (
	"fmt"
	"strings"

	"github.com/codelieche/microservice/usercenter/datamodels"
	"github.com/codelieche/microservice/usercenter/web/forms"
	"github.com/codelieche/microservice/usercenter/web/services"
	"github.com/kataras/iris/v12"
)

type RoleController struct {
	Service           services.RoleService
	UserService       services.UserService
	PermissionService services.PermissionService
}

func (c *RoleController) PostCreate(ctx iris.Context) (role *datamodels.Role, err error) {

	// 1. 定义变量
	var (
		contentType string
		form        *forms.RoleCreateForm
		users       []*datamodels.User
		user        *datamodels.User
		permissions []*datamodels.Permission
		permission  *datamodels.Permission
	)

	// 2. 获取变量
	contentType = ctx.Request().Header.Get("Content-Type")
	form = &forms.RoleCreateForm{}

	if strings.Contains(contentType, "application/json") {
		err = ctx.ReadJSON(form)
	} else {
		err = ctx.ReadForm(form)
	}

	if err != nil {
		return nil, err
	}

	// 3. 实例化role
	role = &datamodels.Role{Name: form.Name}
	// 处理users
	users = []*datamodels.User{}
	if len(form.Users) > 0 {
		for _, i := range form.Users {
			// 判断传入的userid是否ok
			if user, err = c.UserService.GetById(i); err != nil {
				err = fmt.Errorf("用户(id:%d): %s", i, err)
				return nil, err
			} else {
				users = append(users, user)
			}
		}
	}
	role.Users = users
	// 处理Permissions
	permissions = []*datamodels.Permission{}
	if len(form.Permissions) > 0 {
		for _, i := range form.Permissions {
			// 判断传入的userid是否ok
			if permission, err = c.PermissionService.GetById(i); err != nil {
				err = fmt.Errorf("权限(id:%d): %s", i, err)
				return nil, err
			} else {
				permissions = append(permissions, permission)
			}
		}
	}
	role.Permissions = permissions

	// 创建Role
	return c.Service.Create(role)

}

func (c *RoleController) PutBy(id int64, ctx iris.Context) (role *datamodels.Role, err error) {
	// 1. 定义变量
	var (
		contentType string
		form        *forms.RoleCreateForm
		users       []*datamodels.User
		user        *datamodels.User
		permissions []*datamodels.Permission
		permission  *datamodels.Permission
	)

	// 2. 获取变量
	if role, err = c.Service.GetById(id); err != nil {
		return nil, err
	}

	contentType = ctx.Request().Header.Get("Content-Type")
	form = &forms.RoleCreateForm{}

	if strings.Contains(contentType, "application/json") {
		err = ctx.ReadJSON(form)
	} else {
		err = ctx.ReadForm(form)
	}

	if err != nil {
		return nil, err
	}

	// 3. 更新role
	role.Name = form.Name

	// 处理users
	users = []*datamodels.User{}
	if len(form.Users) > 0 {
		for _, i := range form.Users {
			// 判断传入的userid是否ok
			if user, err = c.UserService.GetById(i); err != nil {
				err = fmt.Errorf("用户(id:%d): %s", i, err)
				return nil, err
			} else {
				users = append(users, user)
			}
		}
	}
	role.Users = users

	// 处理Permissions
	permissions = []*datamodels.Permission{}
	if len(form.Permissions) > 0 {
		for _, i := range form.Permissions {
			// 判断传入的userid是否ok
			if permission, err = c.PermissionService.GetById(i); err != nil {
				err = fmt.Errorf("权限(id:%d): %s", i, err)
				return nil, err
			} else {
				permissions = append(permissions, permission)
			}
		}
	}
	role.Permissions = permissions

	// 修改Role
	return c.Service.Save(role)

}

func (c *RoleController) GetBy(idOrName string) (role *datamodels.Role, success bool) {
	if role, err := c.Service.GetByIdOrName(idOrName); err != nil {
		return nil, false
	} else {
		return role, true
	}
}

func (c *RoleController) GetList(ctx iris.Context) (roles []*datamodels.Role, success bool) {
	return c.GetListBy(1, ctx)
}

// 获取用户的列表页：注意以前的版本，ctx是可放前面的
func (c *RoleController) GetListBy(page int, ctx iris.Context) (roles []*datamodels.Role, success bool) {
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
	if roles, err = c.Service.List(offset, limit); err != nil {
		return nil, false
	} else {
		return roles, true
	}
}

// 获取角色的用户列表
func (c *RoleController) GetByUserList(idOrName string, ctx iris.Context) (users []*datamodels.User, success bool) {
	//	定义变量
	var (
		role     *datamodels.Role
		page     int
		pageSize int
		offset   int
		limit    int
		err      error
	)

	// 获取角色
	if role, err = c.Service.GetByIdOrName(idOrName); err != nil {
		return nil, false
	}

	//	获取角色的用户
	// 获取到变量
	pageSize = ctx.URLParamIntDefault("pageSize", 10)
	page = ctx.URLParamIntDefault("page", 1)
	if page > 1 {
		offset = (page - 1) * pageSize
	}
	limit = pageSize

	// 获取用户
	if users, err = c.Service.GetRoleUserList(role, offset, limit); err != nil {
		return nil, false
	} else {
		return users, true
	}
}

// 获取角色的权限列表
func (c *RoleController) GetByPermissionList(idOrName string, ctx iris.Context) (permissions []*datamodels.Permission, success bool) {
	//	定义变量
	var (
		role     *datamodels.Role
		page     int
		pageSize int
		offset   int
		limit    int
		err      error
	)

	// 获取角色
	if role, err = c.Service.GetByIdOrName(idOrName); err != nil {
		return nil, false
	}

	//	获取角色的用户
	// 获取到变量
	pageSize = ctx.URLParamIntDefault("pageSize", 10)
	page = ctx.URLParamIntDefault("page", 1)
	if page > 1 {
		offset = (page - 1) * pageSize
	}
	limit = pageSize

	// 获取权限
	if permissions, err = c.Service.GetRolePermissionList(role, offset, limit); err != nil {
		return nil, false
	} else {
		return permissions, true
	}
}
