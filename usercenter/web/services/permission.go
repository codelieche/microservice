package services

import (
	"github.com/codelieche/microservice/datamodels"
	"github.com/codelieche/microservice/repositories"
)

// Permission Service Interface
type PermissionService interface {
	GetById(id int64) (permission *datamodels.Permission, err error)
	GetByAppIDAndCode(appID int, code string) (permission *datamodels.Permission, err error)
	List(offset int, limit int) (permissions []*datamodels.Permission, err error)
}

// 实例化Permission Service
func NewPermissionService(repo repositories.PermissionRepository) PermissionService {
	return &permissionService{repo: repo}
}

// permission Service
type permissionService struct {
	repo repositories.PermissionRepository
}

func (s *permissionService) GetById(id int64) (permission *datamodels.Permission, err error) {
	return s.repo.Get(id)
}

func (s *permissionService) GetByAppIDAndCode(appID int, code string) (permission *datamodels.Permission, err error) {
	return s.repo.GetByAppIDAndCode(appID, code)
}

func (s *permissionService) List(offset int, limit int) (permissions []*datamodels.Permission, err error) {
	return s.repo.List(offset, limit)
}
