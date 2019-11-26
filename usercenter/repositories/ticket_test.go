package repositories

import (
	"crypto/md5"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/codelieche/microservice/common"
	"github.com/codelieche/microservice/datasources"

	"github.com/codelieche/microservice/datamodels"
)

func generateTicket(sessionID string, returnUrl string) *datamodels.Ticket {

	Md5Inst := md5.New()
	Md5Inst.Write([]byte(fmt.Sprintf("%s-%d", sessionID, time.Now().UnixNano())))

	Result := Md5Inst.Sum([]byte(""))
	ticketName := fmt.Sprintf("%x", Result)
	ticket := datamodels.Ticket{
		Name:        ticketName,
		Session:     sessionID,
		IsActive:    true,
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
			for _, ticket := range tickets {
				log.Println(ticket.ID, ticket.Name, ticket.TimeExpired)
			}
		}
	}
}
