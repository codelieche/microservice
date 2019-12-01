package controllers

import (
	"github.com/codelieche/microservice/usercenter/common"
	"github.com/codelieche/microservice/usercenter/datamodels"
	"github.com/codelieche/microservice/usercenter/web/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
)

// Safe Log Controller
type SafeLogController struct {
	Session *sessions.Session
	Ctx     iris.Context
	Service services.SafeLogService
}

// 获取安全日志
func (c *SafeLogController) GetBy(id int64) (safeLog *datamodels.SafeLog, err error) {
	userID := c.Session.GetIntDefault("userID", 0)
	if userID > 0 {
		if safeLog, err = c.Service.GetById(id); err != nil {
			return nil, err
		} else {
			if safeLog.ID != uint(userID) {
				err := common.NotFountError
				return nil, err
			} else {
				return safeLog, err
			}
		}
	} else {
		c.Ctx.StatusCode(iris.StatusUnauthorized)
		return
	}
}

// 安全日志列表
func (c *SafeLogController) GetList() (safeLogs []*datamodels.SafeLog, err error) {
	return c.GetListBy(1)
}

func (c *SafeLogController) GetListBy(page int) (safeLogs []*datamodels.SafeLog, err error) {
	// 定义变量
	var (
		offset   int
		limit    int
		pageSize int
	)

	// 获取变量
	pageSize = c.Ctx.URLParamIntDefault("pageSize", 10)
	limit = pageSize

	if page > 1 {
		offset = (page - 1) * pageSize
	}
	return c.Service.List(offset, limit)
}
