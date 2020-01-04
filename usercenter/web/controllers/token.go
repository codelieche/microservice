package controllers

import (
	"github.com/codelieche/microservice/usercenter/datamodels"
	"github.com/codelieche/microservice/usercenter/web/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
)

type TokenController struct {
	Session *sessions.Session
	Ctx     iris.Context
	Service services.TokenService
}

// 获取Token的详情
// 由于token.Token包含了一些特殊符号(eg:#)，所以用Token有可能获取不到
func (c *TokenController) GetBy(idOrToken string) (token *datamodels.Token, success bool) {
	if token, err := c.Service.GetByIdOrToken(idOrToken); err != nil {
		return nil, false
	} else {
		return token, true
	}
}

func (c *TokenController) GetList(ctx iris.Context) (tokens []*datamodels.Token, success bool) {
	return c.GetListBy(1, ctx)
}

// 获取用户的列表页：注意以前的版本，ctx是可放前面的
func (c *TokenController) GetListBy(page int, ctx iris.Context) (tokens []*datamodels.Token, success bool) {
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
	if tokens, err = c.Service.List(offset, limit); err != nil {
		return nil, false
	} else {
		return tokens, true
	}
}
