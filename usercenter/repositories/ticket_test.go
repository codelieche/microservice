package repositories

import (
	"crypto/md5"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/codelieche/microservice/usercenter/common"
	"github.com/codelieche/microservice/usercenter/datasources"

	"github.com/codelieche/microservice/usercenter/datamodels"
)

func generateTicket(sessionID string, returnUrl string) *datamodels.Ticket {

	Md5Inst := md5.New()
	Md5Inst.Write([]byte(fmt.Sprintf("%s-%d", sessionID, time.Now().UnixNano())))

	Result := Md5Inst.Sum([]byte(""))
	ticketName := fmt.Sprintf("%x", Result)
	var isActive = true
	ticket := datamodels.Ticket{
		Name:        ticketName,
		Session:     sessionID,
		IsActive:    isActive,
		ReturnUrl:   returnUrl,
		TimeExpired: time.Now().Add(time.Minute),
	}

	return &ticket
}

func TestTicketRepository_Save(t *testing.T) {
	// 1. get db
	db := datasources.GetDb()

	db.AutoMigrate(&datamodels.Ticket{})

	// 2. init ticket repository
	r := NewTicketRepository(db)

	// 创建10条ticket
	i := 0
	for i < 10 {
		i++
		sessionID := common.RandString(32)
		ticket := generateTicket(sessionID, "http://www.codelieche.com")
		if ticket, err := r.Save(ticket); err != nil {
			t.Error(err.Error())
		} else {
			log.Println(ticket.ID, ticket.Name, ticket.CreatedAt, ticket.TimeExpired, ticket.Session)
		}
	}
}

func TestTicketRepository_List(t *testing.T) {
	// 1. get db
	db := datasources.GetDb()

	// 2. init ticket repository
	r := NewTicketRepository(db)

	// 3. list ticket
	haveNext := true
	offset := 0
	limit := 5
	for haveNext {
		if tickets, err := r.List(offset, limit); err != nil {
			t.Error(err.Error())
			haveNext = false
		} else {
			if len(tickets) == limit && limit > 0 {
				haveNext = true
				offset += limit
			} else {
				haveNext = false
			}
			// 输出分组
			updateFields := map[string]interface{}{}
			for _, ticket := range tickets {
				if ticket.ID == 12 {
					updateFields["id"] = 19
					updateFields["Name"] = "update_ticket_namedd"
					updateFields["IsActive"] = true
					updateFields["Times"] = 110
					if ti, err := r.Update(ticket, updateFields); err != nil {
						t.Error(err.Error())
					} else {
						log.Println(ti.ID, ti.Name)
					}
					//r.UpdateByID(12, updateFields)
				}
				log.Println(ticket.ID, ticket.Name, ticket.TimeExpired)
			}
		}
	}
}
