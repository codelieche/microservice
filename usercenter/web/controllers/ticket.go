package controllers

import (
	"errors"

	"github.com/codelieche/microservice/usercenter/web/forms"

	"github.com/codelieche/microservice/usercenter/datamodels"
	"github.com/codelieche/microservice/usercenter/web/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type TicketController struct {
	Service services.TicketService
	Ctx     iris.Context
}

//func (c *TicketController) GetBy(id int) (ticket *datamodels.Ticket, success bool) {
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

// 验证ticket是否ok
func (c *TicketController) GetValidateBy(name string) mvc.Result {
	if ticket, err := c.Service.GetByName(name); err != nil {
		return mvc.Response{
			Object: forms.TicketValidateResponse{
				Errcode:   4001,
				Errmsg:    err.Error(),
				ReturnUrl: "",
				Name:      name,
				Session:   "",
			},
		}
	} else {
		//log.Println(ticket.User)
		// 判断ticket是否校验过了
		if ticket.IsActive {
			// 设置为IsActive为False
			updateFields := map[string]interface{}{}
			updateFields["IsActive"] = false
			// 校验要跳转的url是否是这个

			// 保存一下User
			var user *forms.TicketValidateUser

			if ticket.User != nil {
				user = &forms.TicketValidateUser{
					ID:       ticket.UserID,
					Username: ticket.User.Username,
					Email:    ticket.User.Email,
					Mobile:   ticket.User.Mobile,
					IsActive: ticket.User.IsActive,
				}
			}
			// 保存ticket
			if ticket, err := c.Service.UpdateByID(int64(ticket.ID), updateFields); err != nil {
				return mvc.Response{
					Object: forms.TicketValidateResponse{
						Errcode:   4001,
						Errmsg:    err.Error(),
						ReturnUrl: "",
						Name:      name,
						Session:   "",
					},
				}
			} else {
				if user != nil {
					if user.IsActive {
						return mvc.Response{
							Object: forms.TicketValidateResponse{
								Errcode:   0,
								Errmsg:    "",
								ReturnUrl: ticket.ReturnUrl,
								Name:      ticket.Name,
								Session:   ticket.Session,
								User:      user,
							},
						}
					} else {
						err := errors.New("当前用户已经被禁用了")
						return mvc.Response{
							Object: forms.TicketValidateResponse{
								Errcode:   40001,
								Errmsg:    err.Error(),
								ReturnUrl: ticket.ReturnUrl,
								Name:      ticket.Name,
								Session:   ticket.Session,
								User:      user,
							},
						}
					}

				} else {
					err := errors.New("当前Ticket的用户为空")
					return mvc.Response{
						Object: forms.TicketValidateResponse{
							Errcode:   4001,
							Errmsg:    err.Error(),
							ReturnUrl: "",
							Name:      name,
							Session:   "",
						},
					}
				}

			}
		} else {
			// 已经核验过了，不可再使用了
			err := errors.New("当前Ticket已经被使用过了")
			return mvc.Response{
				Object: forms.TicketValidateResponse{
					Errcode:   4001,
					Errmsg:    err.Error(),
					ReturnUrl: "",
					Name:      name,
					Session:   "",
				},
			}
		}
	}
}
