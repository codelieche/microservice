package services

import (
	"github.com/codelieche/microservice/datamodels"
	"github.com/codelieche/microservice/repositories"
)

// User Service Interface
type UserService interface {
	CreateUser(user *datamodels.User) (*datamodels.User, error)
	GetById(id int64) (user *datamodels.User, err error)
	GetByIdOrName(idOrName string) (user *datamodels.User, err error)
	List(offset int, limit int) (users []*datamodels.User, err error)
	ChangeUserPassword(user *datamodels.User, password string) (*datamodels.User, error)
}

// 实例化User Service
func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

// user Service
type userService struct {
	repo repositories.UserRepository
}

func (s *userService) CreateUser(user *datamodels.User) (*datamodels.User, error) {
	return s.repo.Save(user)
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
