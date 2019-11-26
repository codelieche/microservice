package controllers

import (
	"github.com/codelieche/microservice/datamodels"
	"github.com/codelieche/microservice/web/services"
	"github.com/kataras/iris/v12"
)

type TicketController struct {
	Service services.TicketService
	Ctx     iris.Context
}

//func (c *TicketController) GetBy(id int64) (ticket *datamodels.Ticket, success bool) {
//	if ticket, err := c.Service.GetById(id); err != nil {
//		return nil, false
//	} else {
//		return ticket, true
//	}
//}

func (c *TicketController) GetBy(idOrName string) (ticket *datamodels.Ticket, success bool) {
	if ticket, err := c.Service.GetByIdOrName(idOrName); err != nil {
		return nil, false
	} else {
		return ticket, true
	}
}

func (c *TicketController) GetList(ctx iris.Context) (tickets []*datamodels.Ticket, success bool) {
	return c.GetListBy(1)
}

// 获取用户的列表页：注意以前的版本，ctx是可放前面的
func (c *TicketController) GetListBy(page int) (tickets []*datamodels.Ticket, success bool) {
	//	定义变量
	var (
		pageSize int
		offset   int
		limit    int
		err      error
	)

	// 获取变量
	pageSize = c.Ctx.URLParamIntDefault("pageSize", 10)
	limit = pageSize
	if page > 1 {
		offset = (page - 1) * pageSize
	}

	//	获取用户
	//log.Println(offset, limit)
	if tickets, err = c.Service.List(offset, limit); err != nil {
		return nil, false
	} else {
		return tickets, true
	}
}
