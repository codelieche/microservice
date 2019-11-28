package controllers

import (
	"github.com/codelieche/microservice/usercenter/datamodels"
	"github.com/codelieche/microservice/usercenter/web/services"
	"github.com/kataras/iris/v12"
)

type RoleController struct {
	Service services.RoleService
}

//func (c *RoleController) GetBy(id int64) (role *datamodels.Role, success bool) {
//	if role, err := c.Service.GetById(id); err != nil {
//		return nil, false
//	} else {
//		return role, true
//	}
//}

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
