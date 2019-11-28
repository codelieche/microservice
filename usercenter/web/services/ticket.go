package services

import (
	"github.com/codelieche/microservice/usercenter/datamodels"
	"github.com/codelieche/microservice/usercenter/repositories"
)

// Ticket Service Interface
type TicketService interface {
	Save(ticket *datamodels.Ticket) (*datamodels.Ticket, error)
	Update(ticket *datamodels.Ticket, fields map[string]interface{}) (*datamodels.Ticket, error)
	UpdateByID(id int64, fields map[string]interface{}) (*datamodels.Ticket, error)
	GetById(id int64) (ticket *datamodels.Ticket, err error)
	GetByName(idOrName string) (ticket *datamodels.Ticket, err error)
	GetByIdOrName(idOrName string) (ticket *datamodels.Ticket, err error)
	List(offset int, limit int) (tickets []*datamodels.Ticket, err error)
}

// 实例化Ticket Service
func NewTicketService(repo repositories.TicketRepository) TicketService {
	return &ticketService{repo: repo}
}

// ticket Service
type ticketService struct {
	repo repositories.TicketRepository
}

func (s *ticketService) UpdateByID(id int64, fields map[string]interface{}) (*datamodels.Ticket, error) {
	return s.repo.UpdateByID(id, fields)
}

func (s *ticketService) Save(ticket *datamodels.Ticket) (*datamodels.Ticket, error) {
	return s.repo.Save(ticket)
}

func (s *ticketService) Update(ticket *datamodels.Ticket, fields map[string]interface{}) (*datamodels.Ticket, error) {
	return s.repo.Update(ticket, fields)
}

func (s *ticketService) GetById(id int64) (ticket *datamodels.Ticket, err error) {
	return s.repo.Get(id)
}

func (s *ticketService) GetByName(name string) (ticket *datamodels.Ticket, err error) {
	return s.repo.GetByName(name)
}

func (s *ticketService) GetByIdOrName(idOrName string) (ticket *datamodels.Ticket, err error) {
	return s.repo.GetByIdOrName(idOrName)
}

// 获取用户Ticket列表
func (s *ticketService) List(offset int, limit int) (tickets []*datamodels.Ticket, err error) {
	return s.repo.List(offset, limit)
}
