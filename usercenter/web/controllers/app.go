package controllers

import (
	"strings"

	"github.com/codelieche/microservice/usercenter/datamodels"
	"github.com/codelieche/microservice/usercenter/web/forms"
	"github.com/codelieche/microservice/usercenter/web/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
)

type ApplicationController struct {
	Session           *sessions.Session
	Ctx               iris.Context
	Service           services.ApplicationService
	PermissionService services.PermissionService
}

func (c *ApplicationController) PostCreate(ctx iris.Context) (app *datamodels.Application, err error) {
	var (
		form        *forms.ApplicationCreateForm
		contentType string
	)

	contentType = ctx.Request().Header.Get("Content-Type")

	form = &forms.ApplicationCreateForm{}
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

	// 创建app
	app = &datamodels.Application{
		Name: form.Name,
		Code: form.Code,
		Info: form.Info,
	}
	return c.Service.Create(app)

}

func (c *ApplicationController) PutBy(id int64, ctx iris.Context) (app *datamodels.Application, err error) {
	var (
		form        *forms.ApplicationUpdateForm
		contentType string
		//permissions []*datamodels.Permission
		//permission  *datamodels.Permission
	)

	contentType = ctx.Request().Header.Get("Content-Type")

	form = &forms.ApplicationUpdateForm{}
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

	// 创建app
	if app, err = c.Service.Get(id); err != nil {
		return nil, err
	}
	// 对新的字段重新赋值
	// 权限
	//permissions = []*datamodels.Permission{}
	//if len(form.Permissions) > 0 {
	//	for _, i := range form.Permissions {
	//		if permission, err = c.PermissionService.GetById(i); err != nil {
	//			err = fmt.Errorf("权限(id:%d)：%s", i, err)
	//			return nil, err
	//		} else {
	//			permissions = append(permissions, permission)
	//		}
	//	}
	//}
	//app.Permissions = permissions

	// 对name、code、info进行赋值
	//app.Name = form.Name
	//app.Code = form.Code
	app.Info = form.Info

	return c.Service.Save(app)

}

// 获取应用的列表
func (c *ApplicationController) GetBy(idOrCode string) (app *datamodels.Application, success bool) {
	if app, err := c.Service.GetByIdOrCode(idOrCode); err != nil {
		return nil, false
	} else {
		return app, true
	}
}

func (c *ApplicationController) GetList(ctx iris.Context) (apps []*datamodels.Application, success bool) {
	return c.GetListBy(1, ctx)
}

// 获取应用的列表页：注意以前的版本，ctx是可放前面的
func (c *ApplicationController) GetListBy(page int, ctx iris.Context) (apps []*datamodels.Application, success bool) {
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

	//	获取app
	//log.Println(offset, limit)
	if apps, err = c.Service.List(offset, limit); err != nil {
		return nil, false
	} else {
		return apps, true
	}
}

// 获取权限的列表页：注意以前的版本，ctx是可放前面的
func (c *ApplicationController) GetByPermissionList(idOrCode string, ctx iris.Context) (permissions []*datamodels.Permission, success bool) {
	return c.GetByPermissionListBy(idOrCode, 1, ctx)
}

// 获取权限的列表页：注意以前的版本，ctx是可放前面的
func (c *ApplicationController) GetByPermissionListBy(idOrCode string, page int, ctx iris.Context) (permissions []*datamodels.Permission, success bool) {
	//	定义变量
	var (
		app      *datamodels.Application
		pageSize int
		offset   int
		limit    int
		err      error
	)

	// 获取变量
	// 获取app
	if app, err = c.Service.GetByIdOrCode(idOrCode); err != nil {
		return nil, false
	}

	pageSize = ctx.URLParamIntDefault("pageSize", 10)
	limit = pageSize
	if page > 1 {
		offset = (page - 1) * pageSize
	}

	//	获取app
	//log.Println(offset, limit)
	if permissions, err = c.Service.GetPermissionList(app, offset, limit); err != nil {
		return nil, false
	} else {
		return permissions, true
	}
}
