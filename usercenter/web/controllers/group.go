package controllers

import (
	"fmt"
	"log"
	"strings"

	"github.com/codelieche/microservice/usercenter/datamodels"
	"github.com/codelieche/microservice/usercenter/web/forms"
	"github.com/codelieche/microservice/usercenter/web/services"
	"github.com/kataras/iris/v12"
)

type GroupController struct {
	Service           services.GroupService
	UserService       services.UserService
	PermissionService services.PermissionService
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

// 创建分组
func (c *GroupController) PostCreate(ctx iris.Context) (group *datamodels.Group, err error) {
	// 1. 定义变量
	var (
		contentType string
		form        *forms.GroupCreateFrom
		user        *datamodels.User
		users       []*datamodels.User
		permission  *datamodels.Permission
		permissions []*datamodels.Permission
		id          int64
	)
	// 2. 读取表单
	// 2-1: 获取Content-Type
	contentType = ctx.Request().Header.Get("Content-Type")

	// 2-2: 根据不同的Content-Type用不同方式解析
	form = &forms.GroupCreateFrom{}

	if strings.Contains(contentType, "application/json") {
		// 传递的是json方式的
		err = ctx.ReadJSON(form)
	} else {
		// multipart/form-data; boundary=----xxx 或者 application/x-www-form-urlencoded
		err = ctx.ReadForm(form)
	}
	// 2-3：判断是否有错误
	if err != nil {
		log.Println(err)
		return nil, err
	} else {
		//log.Println(form)
	}

	// 3. 创建Group
	// 3-1: 实例化Group
	group = &datamodels.Group{Name: form.Name}
	// 3-2: 获取users
	if len(form.Users) > 0 {
		users = []*datamodels.User{}
		for _, id = range form.Users {
			if user, err = c.UserService.GetById(id); err != nil {
				err = fmt.Errorf("users(ID: %d): %s", id, err)
				return nil, err
			} else {
				users = append(users, user)
			}
		}
		// 给group赋值users
		group.Users = users
	}
	// 3-3: 获取Permissions
	if len(form.Permissions) > 0 {
		permissions = []*datamodels.Permission{}
		for _, id = range form.Permissions {
			if permission, err = c.PermissionService.GetById(id); err != nil {
				err = fmt.Errorf("permission(ID:%d): %s", id, err)
				return nil, err
			} else {
				permissions = append(permissions, permission)
			}
		}
		// 给group赋值Permissions
		group.Permissions = permissions
	}

	// 4. 创建Group
	return c.Service.Create(group)
}

// 创建分组
func (c *GroupController) PutBy(id int64, ctx iris.Context) (group *datamodels.Group, err error) {
	// 1. 定义变量
	var (
		contentType string
		form        *forms.GroupCreateFrom
		user        *datamodels.User
		users       []*datamodels.User
		permission  *datamodels.Permission
		permissions []*datamodels.Permission
	)
	// 2. 读取表单
	// 2-1: 获取Content-Type
	contentType = ctx.Request().Header.Get("Content-Type")

	// 2-2: 根据不同的Content-Type用不同方式解析
	form = &forms.GroupCreateFrom{}

	if strings.Contains(contentType, "application/json") {
		// 传递的是json方式的
		err = ctx.ReadJSON(form)
	} else {
		// multipart/form-data; boundary=----xxx 或者 application/x-www-form-urlencoded
		err = ctx.ReadForm(form)
	}
	// 2-3：判断是否有错误
	if err != nil {
		log.Println(err)
		return nil, err
	} else {
		//log.Println(form)
	}

	// 3. 更新Group
	// 3-1: Group
	if group, err = c.Service.GetById(id); err != nil {
		return nil, err
	}
	// 3-2: 获取users
	if len(form.Users) > 0 {
		users = []*datamodels.User{}
		for _, i := range form.Users {
			if user, err = c.UserService.GetById(i); err != nil {
				err = fmt.Errorf("users(ID: %d): %s", i, err)
				return nil, err
			} else {
				users = append(users, user)
			}
		}
		// 给group赋值users
		group.Users = users
	} else {
		group.Users = nil
	}
	// 3-3: 获取Permissions
	if len(form.Permissions) > 0 {
		permissions = []*datamodels.Permission{}
		for _, i := range form.Permissions {
			if permission, err = c.PermissionService.GetById(i); err != nil {
				err = fmt.Errorf("permission(ID:%d): %s", i, err)
				return nil, err
			} else {
				permissions = append(permissions, permission)
			}
		}
		// 给group赋值Permissions
		group.Permissions = permissions
	} else {
		group.Permissions = nil
	}

	// 3-4: 给group.name赋新的值
	group.Name = form.Name

	// 4. 更新Group
	return c.Service.Update(group)
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
