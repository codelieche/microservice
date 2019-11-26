package repositories

import (
	"github.com/codelieche/microservice/common"
	"github.com/codelieche/microservice/datamodels"
	"github.com/jinzhu/gorm"
)

type TicketRepository interface {
	// 保存Ticket
	Save(ticket *datamodels.Ticket) (*datamodels.Ticket, error)
	// 获取Ticket的列表
	List(offset int, limit int) ([]*datamodels.Ticket, error)
	// 获取Ticket信息
	Get(id int64) (*datamodels.Ticket, error)
	// 根据Name获取Ticket信息
	GetByName(name string) (*datamodels.Ticket, error)
	// 根据ID或者Name获取Ticket信息
	GetByIdOrName(idOrName string) (*datamodels.Ticket, error)
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{db: db}
}

type ticketRepository struct {
	db *gorm.DB
}

// 保存Ticket
func (r *ticketRepository) Save(ticket *datamodels.Ticket) (*datamodels.Ticket, error) {
	if ticket.ID > 0 {
		// 是更新操作
		if err := r.db.Model(&datamodels.Ticket{}).Update(ticket).Error; err != nil {
			return nil, err
		} else {
			return ticket, nil
		}
	} else {
		// 是创建操作
		if err := r.db.Create(ticket).Error; err != nil {
			return nil, err
		} else {
			return ticket, nil
		}

	}
}

// 获取Ticket的列表
func (r *ticketRepository) List(offset int, limit int) (tickets []*datamodels.Ticket, err error) {
	query := r.db.Model(&datamodels.Ticket{}).Offset(offset).Limit(limit).Find(&tickets)
	if query.Error != nil {
		return nil, query.Error
	} else {
		return tickets, nil
	}
	return
}

// 根据ID获取Ticket
func (r *ticketRepository) Get(id int64) (ticket *datamodels.Ticket, err error) {

	ticket = &datamodels.Ticket{}
	r.db.First(ticket, "id = ?", id)
	if ticket.ID > 0 {
		return ticket, nil
	} else {
		return nil, common.NotFountError
	}
}

func (r *ticketRepository) GetByName(name string) (ticket *datamodels.Ticket, err error) {
	ticket = &datamodels.Ticket{}
	r.db.First(ticket, "name = ?", name)
	if ticket.ID > 0 {
		return ticket, nil
	} else {
		return nil, common.NotFountError
	}
}

// 根据ID/Name获取Ticket
func (r *ticketRepository) GetByIdOrName(idOrName string) (ticket *datamodels.Ticket, err error) {

	ticket = &datamodels.Ticket{}
	r.db.First(ticket, "id = ? or name = ?", idOrName, idOrName)
	if ticket.ID > 0 {
		return ticket, nil
	} else {
		return nil, common.NotFountError
	}
}
