package services

import (
	"github.com/codelieche/microservice/usercenter/datamodels"
	"github.com/codelieche/microservice/usercenter/repositories"
)

type ApplicationService interface {
	// 创建App
	Create(app *datamodels.Application) (*datamodels.Application, error)
	// 保存App
	Save(app *datamodels.Application) (*datamodels.Application, error)
	// 获取App的列表
	List(offset int, limit int) ([]*datamodels.Application, error)
	// 获取App信息
	Get(id int64) (*datamodels.Application, error)
	// 根据ID或者Code获取App信息
	GetByIdOrCode(idOrCode string) (app *datamodels.Application, err error)
	// 获取App的Permission列表
	GetPermissionList(app *datamodels.Application, offset int, limit int) ([]*datamodels.Permission, error)
}

func NewApplicationService(repo repositories.ApplicationRepository) ApplicationService {
	return &applicationService{repo: repo}
}

type applicationService struct {
	repo repositories.ApplicationRepository
}

func (s *applicationService) Create(app *datamodels.Application) (*datamodels.Application, error) {
	return s.repo.Save(app)
}

func (s *applicationService) Save(app *datamodels.Application) (*datamodels.Application, error) {
	return s.repo.Save(app)
}

func (s *applicationService) List(offset int, limit int) ([]*datamodels.Application, error) {
	return s.repo.List(offset, limit)
}

func (s *applicationService) Get(id int64) (*datamodels.Application, error) {
	return s.repo.Get(id)
}

func (s *applicationService) GetByIdOrCode(idOrCode string) (app *datamodels.Application, err error) {
	return s.repo.GetByIdOrCode(idOrCode)
}

func (s *applicationService) GetPermissionList(app *datamodels.Application, offset int, limit int) ([]*datamodels.Permission, error) {
	return s.repo.GetPermissionList(app, offset, limit)
}
