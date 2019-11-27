package repositories

import (
	"errors"

	"github.com/codelieche/microservice/common"
	"github.com/codelieche/microservice/datamodels"
	"github.com/jinzhu/gorm"
)

type TicketRepository interface {
	// 保存Ticket
	Save(ticket *datamodels.Ticket) (*datamodels.Ticket, error)
	// Update
	Update(ticket *datamodels.Ticket, fields map[string]interface{}) (*datamodels.Ticket, error)
	// Update
	UpdateByID(id int64, fields map[string]interface{}) (*datamodels.Ticket, error)
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
		if err := r.db.Model(ticket).Update(ticket).Error; err != nil {
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

// 更新对象
func (r *ticketRepository) Update(ticket *datamodels.Ticket, fields map[string]interface{}) (*datamodels.Ticket, error) {
	// 判断ID:
	// 如果传入的是0，那么会更新全部
	// 如果fields中传入ID，那么会更新ID是它的对象
	// 推荐加个limit(1), 确保只更新一条数据
	if ticket.ID <= 0 {
		err := errors.New("传入ID为0,会更新全部数据")
		return nil, err
	}
	// 丢弃ID/id/Id/iD
	idKeys := []string{"ID", "id", "Id", "iD"}
	for _, k := range idKeys {
		if _, exist := fields[k]; exist {
			delete(fields, k)
		}
	}

	// 更新操作
	if err := r.db.Model(ticket).Limit(1).Update(fields).Error; err != nil {
		return nil, err
	} else {
		return ticket, nil
	}
}

func (r *ticketRepository) UpdateByID(id int64, fields map[string]interface{}) (ticket *datamodels.Ticket, err error) {
	// 判断ID
	if id <= 0 {
		err := errors.New("传入ID为0,会更新全部数据")
		return nil, err
	}
	// 因为指定了ID了，所以这里可不判断这个ID
	// 丢弃ID/id/Id/iD
	//idKeys := []string{"ID", "id", "Id", "iD"}
	//for _, k := range idKeys {
	//	if _, exist := fields[k]; exist {
	//		delete(fields, k)
	//	}
	//}

	// 更新操作
	if err = r.db.Model(&datamodels.Ticket{}).Where("id = ?", id).Limit(1).Update(fields).Error; err != nil {
		return nil, err
	} else {
		return r.Get(id)
		//ticket = &datamodels.Ticket{}
		//if err = r.db.First(ticket, "id = ?", id).Error; err != nil {
		//	return nil, err
		//} else {
		//	return ticket, nil
		//}
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
