package controllers

import (
	"fmt"
	"strings"

	"github.com/codelieche/microservice/usercenter/datamodels"
	"github.com/codelieche/microservice/usercenter/web/forms"
	"github.com/codelieche/microservice/usercenter/web/services"
	"github.com/kataras/iris/v12"
)

type PermissionController struct {
	Service    services.PermissionService
	AppService services.ApplicationService
}

func (c *PermissionController) PostCreate(ctx iris.Context) (permission *datamodels.Permission, err error) {
	// 定义变量
	var (
		contentType string
		form        *forms.PermissionCreateForm
		app         *datamodels.Application
	)
	contentType = ctx.Request().Header.Get("Content-Type")

	form = &forms.PermissionCreateForm{}

	if strings.Contains(contentType, "application/json") {
		err = ctx.ReadJSON(form)
	} else {
		err = ctx.ReadForm(form)
	}

	if err != nil {
		return nil, err
	}

	// 实例化Permission
	permission = &datamodels.Permission{
		Name:  form.Name,
		Code:  form.Code,
		AppID: uint(form.AppID),
	}
	// 判断app_id是否存在
	if app, err = c.AppService.Get(form.AppID); err != nil {
		err = fmt.Errorf("应用(id:%d)出错:%s", form.AppID, err)
		return nil, err
	}

	if permission, err = c.Service.Create(permission); err != nil {
		return nil, err
	} else {
		return permission, nil
	}
}

func (c *PermissionController) PutBy(id int64, ctx iris.Context) (permission *datamodels.Permission, err error) {
	// 定义变量
	var (
		app         *datamodels.Application
		contentType string
		form        *forms.PermissionCreateForm
	)
	contentType = ctx.Request().Header.Get("Content-Type")

	form = &forms.PermissionCreateForm{}

	if strings.Contains(contentType, "application/json") {
		err = ctx.ReadJSON(form)
	} else {
		err = ctx.ReadForm(form)
	}

	if err != nil {
		return nil, err
	} else {
		//log.Println(form)
	}

	// 获取permission
	if permission, err = c.Service.GetById(id); err != nil {
		return nil, err
	}

	// 判断app_id是否存在
	if app, err = c.AppService.Get(form.AppID); err != nil {
		err = fmt.Errorf("获取应用(id:%d)出错:%s", form.AppID, err)
		return nil, err
	}

	permission.Application = app
	permission.AppID = uint(form.AppID)
	permission.Code = form.Code
	permission.Name = form.Name

	// 对Application进行更新
	if permission, err = c.Service.Save(permission); err != nil {
		return nil, err
	} else {
		return permission, nil
	}
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
