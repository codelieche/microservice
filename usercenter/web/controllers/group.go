package controllers

import (
	"github.com/codelieche/microservice/datamodels"
	"github.com/codelieche/microservice/web/services"
	"github.com/kataras/iris"
)

type GroupController struct {
	Service services.GroupService
}

//func (c *GroupController) GetBy(id int64) (group *datamodels.Group, success bool) {
//	if group, err := c.Service.GetById(id); err != nil {
//		return nil, false
//	} else {
//		return group, true
//	}
//}

func (c *GroupController) GetBy(idOrName string) (group *datamodels.Group, success bool) {
	if group, err := c.Service.GetByIdOrName(idOrName); err != nil {
		return nil, false
	} else {
		return group, true
	}
}

func (c *GroupController) GetList(ctx iris.Context) (groups []*datamodels.Group, success bool) {
	return c.GetListBy(1, ctx)
}

// 获取用户的列表页：注意以前的版本，ctx是可放前面的
func (c *GroupController) GetListBy(page int, ctx iris.Context) (groups []*datamodels.Group, success bool) {
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
	if groups, err = c.Service.List(offset, limit); err != nil {
		return nil, false
	} else {
		return groups, true
	}
}

// 获取分组的用户列表
func (c *GroupController) GetByUserList(idOrName string, ctx iris.Context) (users []*datamodels.User, success bool) {
	//	定义变量
	var (
		group    *datamodels.Group
		page     int
		pageSize int
		offset   int
		limit    int
		err      error
	)

	// 获取分组
	if group, err = c.Service.GetByIdOrName(idOrName); err != nil {
		return nil, false
	}

	//	获取分组的用户
	// 获取到变量
	pageSize = ctx.URLParamIntDefault("pageSize", 10)
	page = ctx.URLParamIntDefault("page", 1)
	if page > 1 {
		offset = (page - 1) * pageSize
	}
	limit = pageSize

	// 获取用户
	if users, err = c.Service.GetGroupUserList(group, offset, limit); err != nil {
		return nil, false
	} else {
		return users, true
	}
}

// 获取分组的权限列表
func (c *GroupController) GetByPermissionList(idOrName string, ctx iris.Context) (permissions []*datamodels.Permission, success bool) {
	//	定义变量
	var (
		group    *datamodels.Group
		page     int
		pageSize int
		offset   int
		limit    int
		err      error
	)

	// 获取分组
	if group, err = c.Service.GetByIdOrName(idOrName); err != nil {
		return nil, false
	}

	//	获取分组的用户
	// 获取到变量
	pageSize = ctx.URLParamIntDefault("pageSize", 10)
	page = ctx.URLParamIntDefault("page", 1)
	if page > 1 {
		offset = (page - 1) * pageSize
	}
	limit = pageSize

	// 获取权限
	if permissions, err = c.Service.GetGroupPermissionList(group, offset, limit); err != nil {
		return nil, false
	} else {
		return permissions, true
	}
}
