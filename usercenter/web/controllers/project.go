package controllers

import (
	"strings"

	"github.com/kataras/iris/v12/mvc"

	"github.com/codelieche/microservice/usercenter/datamodels"
	"github.com/codelieche/microservice/usercenter/web/forms"
	"github.com/codelieche/microservice/usercenter/web/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
)

type ProjectController struct {
	Session           *sessions.Session
	Ctx               iris.Context
	Service           services.ProjectService
	PermissionService services.PermissionService
}

func (c *ProjectController) PostCreate(ctx iris.Context) (project *datamodels.Project, err error) {
	var (
		form        *forms.ProjectCreateForm
		contentType string
	)

	contentType = ctx.Request().Header.Get("Content-Type")

	form = &forms.ProjectCreateForm{}
	if strings.Contains(contentType, "projectlication/json") {
		err = ctx.ReadJSON(form)
	} else {
		err = ctx.ReadForm(form)
	}

	if err != nil {
		return nil, err
	} else {
		//log.Print(form)
	}

	// 创建项目
	project = &datamodels.Project{
		Name:        form.Name,
		Code:        form.Code,
		Description: form.Description,
	}

	return c.Service.Create(project)
}

func (c *ProjectController) PutBy(id int64, ctx iris.Context) (project *datamodels.Project, err error) {
	var (
		form        *forms.ProjectCreateForm
		contentType string
	)

	contentType = ctx.Request().Header.Get("Content-Type")

	form = &forms.ProjectCreateForm{}
	if strings.Contains(contentType, "projectlication/json") {
		err = ctx.ReadJSON(form)
	} else {
		err = ctx.ReadForm(form)
	}

	if err != nil {
		return nil, err
	} else {
		//log.Print(form)
	}

	// 获取Project
	if project, err = c.Service.Get(id); err != nil {
		return nil, err
	} else {
		project.Name = form.Name
		project.Description = form.Description
		return c.Service.Save(project)
	}
}

// 获取项目
func (c *ProjectController) GetBy(idOrCode string) (project *datamodels.Project, success bool) {
	if project, err := c.Service.GetByIdOrCode(idOrCode); err != nil {
		return nil, false
	} else {
		return project, true
	}
}

func (c *ProjectController) GetList(ctx iris.Context) (projects []*datamodels.Project, success bool) {
	return c.GetListBy(1, ctx)
}

// 获取应用的列表页：注意以前的版本，ctx是可放前面的
func (c *ProjectController) GetListBy(page int, ctx iris.Context) (projects []*datamodels.Project, success bool) {
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

	//	获取projects
	//log.Println(offset, limit)
	if projects, _, err = c.Service.List(offset, limit); err != nil {
		return nil, false
	} else {
		return projects, true
	}
}

// 获取权限的列表页：注意以前的版本，ctx是可放前面的
func (c *ProjectController) GetByPermissionList(idOrCode string, ctx iris.Context) (permissions []*datamodels.Permission, success bool) {
	return c.GetByPermissionListBy(idOrCode, 1, ctx)
}

// 获取权限的列表页：注意以前的版本，ctx是可放前面的
func (c *ProjectController) GetByPermissionListBy(idOrCode string, page int, ctx iris.Context) (permissions []*datamodels.Permission, success bool) {
	//	定义变量
	var (
		project  *datamodels.Project
		pageSize int
		offset   int
		limit    int
		err      error
	)

	// 获取变量
	// 获取project
	if project, err = c.Service.GetByIdOrCode(idOrCode); err != nil {
		return nil, false
	}

	pageSize = ctx.URLParamIntDefault("pageSize", 10)
	limit = pageSize
	if page > 1 {
		offset = (page - 1) * pageSize
	}

	//	获取project
	//log.Println(offset, limit)
	if permissions, err = c.Service.GetPermissionList(project, offset, limit); err != nil {
		return nil, false
	} else {
		return permissions, true
	}
}

func (c *ProjectController) PostUserAdd(ctx iris.Context) (pUser *datamodels.ProjectUser, err error) {
	var (
		form        *forms.ProjectUserAddForm
		contentType string
	)
	contentType = ctx.Request().Header.Get("Content-Type")

	form = &forms.ProjectUserAddForm{}

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

	// 添加用户
	return c.Service.AddProjectUser(form.Project, form.Username, form.Role)
}

func (c *ProjectController) DeleteUser(ctx iris.Context) mvc.Response {
	var (
		form        *forms.ProjectUserAddForm
		contentType string
		err         error
	)
	contentType = ctx.Request().Header.Get("Content-Type")

	form = &forms.ProjectUserAddForm{}

	if strings.Contains(contentType, "application/json") {
		err = ctx.ReadJSON(form)
	} else {
		err = ctx.ReadForm(form)
	}

	if err != nil {
		return mvc.Response{
			Code: 400,
			Err:  err,
		}
	} else {
		//log.Println(form)
	}

	// 添加用户
	if success, err := c.Service.DeleteProjectUser(form.Project, form.Username); err != nil {
		return mvc.Response{
			Code: 400,
			Err:  err,
		}
	} else {
		if success {
			return mvc.Response{
				Code: 204,
			}
		} else {
			return mvc.Response{
				Code: 400,
			}
		}
	}
}

// 获取权限的列表页：注意以前的版本，ctx是可放前面的
func (c *ProjectController) GetByUsersList(idOrCode string, ctx iris.Context) (users []*datamodels.ProjectUser, success bool) {
	return c.GetByUsersListBy(idOrCode, 1, ctx)
}

// 获取权限的列表页：注意以前的版本，ctx是可放前面的
func (c *ProjectController) GetByUsersListBy(idOrCode string, page int, ctx iris.Context) (users []*datamodels.ProjectUser, success bool) {
	//	定义变量
	var (
		project  *datamodels.Project
		pageSize int
		offset   int
		limit    int
		err      error
	)

	// 获取变量
	// 获取project
	if project, err = c.Service.GetByIdOrCode(idOrCode); err != nil {
		return nil, false
	}

	pageSize = ctx.URLParamIntDefault("pageSize", 10)
	limit = pageSize
	if page > 1 {
		offset = (page - 1) * pageSize
	}

	//	获取project
	//log.Println(offset, limit)
	//var count int
	if users, _, err = c.Service.GetUsersList(project, offset, limit); err != nil {
		return nil, false
	} else {
		//log.Print("用户的总数", count)
		return users, true
	}
}
