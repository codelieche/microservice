package services

import (
	"errors"
	"log"

	"github.com/codelieche/microservice/usercenter/common"

	"github.com/codelieche/microservice/usercenter/datamodels"
	"github.com/codelieche/microservice/usercenter/repositories"
)

// User Service Interface
type UserService interface {
	CreateUser(user *datamodels.User) (*datamodels.User, error)
	Save(user *datamodels.User) (*datamodels.User, error)
	GetById(id int64) (user *datamodels.User, err error)
	Update(user *datamodels.User, fields map[string]interface{}) (*datamodels.User, error)
	GetByIdOrName(idOrName string) (user *datamodels.User, err error)
	List(offset int, limit int) (users []*datamodels.User, err error)
	ChangeUserPassword(user *datamodels.User, password string) (*datamodels.User, error)
	// 检查用户的密码
	CheckUserPassword(user *datamodels.User, password string) (bool, error)
	SaveTicket(ticket *datamodels.Ticket) (*datamodels.Ticket, error)
	DeleteUserByIdOrName(idOrName string) (success bool, err error)
}

// 实例化User Service
func NewUserService(repo repositories.UserRepository, tRepo repositories.TicketRepository) UserService {
	return &userService{repo: repo, ticketRepo: tRepo}
}

// user Service
type userService struct {
	repo       repositories.UserRepository
	ticketRepo repositories.TicketRepository
}

func (s *userService) DeleteUserByIdOrName(idOrName string) (success bool, err error) {
	// 1. 先获取用户
	if user, err := s.GetByIdOrName(idOrName); err != nil {
		//log.Println(err)
		return false, err
	} else {
		// 2. 禁用用户
		isActive := user.IsActive
		updateFields := map[string]interface{}{
			"IsActive": false,
			//"DeletedAt": time.Now(),
		}
		if user, err = s.repo.UpdateByID(int64(user.ID), updateFields); err != nil {
			log.Println(err)
			if err == common.NotFountError {
				return true, nil
			} else {
				return false, err
			}
		} else {
			log.Println(user)
			if isActive && !user.IsActive {
				return true, nil
			} else {
				return false, errors.New("删除失败")

			}
		}
	}
}

func (s *userService) CreateUser(user *datamodels.User) (*datamodels.User, error) {
	return s.repo.Save(user)
}

func (s *userService) Save(user *datamodels.User) (*datamodels.User, error) {
	return s.repo.Save(user)
}

func (s *userService) Update(user *datamodels.User, fields map[string]interface{}) (*datamodels.User, error) {
	return s.repo.UpdateByID(int64(user.ID), fields)
}

func (s *userService) GetById(id int64) (user *datamodels.User, err error) {
	return s.repo.Get(id)
}

func (s *userService) GetByIdOrName(idOrName string) (user *datamodels.User, err error) {
	return s.repo.GetByIdOrName(idOrName)
}

func (s *userService) List(offset int, limit int) (users []*datamodels.User, err error) {
	return s.repo.List(offset, limit)
}

func (s *userService) ChangeUserPassword(user *datamodels.User, password string) (*datamodels.User, error) {
	return s.repo.SetUserPassword(user, password)
}

func (s *userService) CheckUserPassword(user *datamodels.User, password string) (bool, error) {
	return s.repo.CheckUserPassword(user, password)
}

func (s *userService) SaveTicket(ticket *datamodels.Ticket) (*datamodels.Ticket, error) {
	return s.ticketRepo.Save(ticket)
}
