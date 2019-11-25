package controllers

import (
	"github.com/codelieche/microservice/datamodels"
	"github.com/codelieche/microservice/web/services"
	"github.com/kataras/iris/v12"
)

type UserController struct {
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
