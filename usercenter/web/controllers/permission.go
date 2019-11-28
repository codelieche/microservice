package controllers

import (
	"github.com/codelieche/microservice/usercenter/datamodels"
	"github.com/codelieche/microservice/usercenter/web/services"
	"github.com/kataras/iris/v12"
)

type PermissionController struct {
	Service services.PermissionService
}

func (c *PermissionController) GetBy(id int64) (permission *datamodels.Permission, success bool) {
	if permission, err := c.Service.GetById(id); err != nil {
		return nil, false
	} else {
		return permission, true
	}
}

func (c *PermissionController) GetByBy(appID int, code string) (permission *datamodels.Permission, success bool) {
	if permission, err := c.Service.GetByAppIDAndCode(appID, code); err != nil {
		return nil, false
	} else {
		return permission, true
	}
}

func (c *PermissionController) GetList(ctx iris.Context) (permissions []*datamodels.Permission, success bool) {
	return c.GetListBy(1, ctx)
}

// 获取权限的列表页：注意以前的版本，ctx是可放前面的
func (c *PermissionController) GetListBy(page int, ctx iris.Context) (permissions []*datamodels.Permission, success bool) {
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

	//	获取权限
	//log.Println(offset, limit)
	if permissions, err = c.Service.List(offset, limit); err != nil {
		return nil, false
	} else {
		return permissions, true
	}
}
