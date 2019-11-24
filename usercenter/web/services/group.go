package services

import (
	"github.com/codelieche/microservice/datamodels"
	"github.com/codelieche/microservice/repositories"
)

// Group Service Interface
type GroupService interface {
	GetById(id int64) (group *datamodels.Group, err error)
	GetByIdOrName(idOrName string) (group *datamodels.Group, err error)
	List(offset int, limit int) (groups []*datamodels.Group, err error)
	GetGroupUserList(group *datamodels.Group, offset int, limit int) (users []*datamodels.User, err error)
	GetGroupPermissionList(group *datamodels.Group, offset int, limit int) (permissions []*datamodels.Permission, err error)
}

// 实例化Group Service
func NewGroupService(repo repositories.GroupRepository) GroupService {
	return &groupService{repo: repo}
}

// group Service
type groupService struct {
	repo repositories.GroupRepository
}

func (s *groupService) GetById(id int64) (group *datamodels.Group, err error) {
	return s.repo.Get(id)
}

func (s *groupService) GetByIdOrName(idOrName string) (group *datamodels.Group, err error) {
	return s.repo.GetByIdOrName(idOrName)
}

// 获取用户分组列表
func (s *groupService) List(offset int, limit int) (groups []*datamodels.Group, err error) {
	return s.repo.List(offset, limit)
}

// 获取分组的用户列表
func (s *groupService) GetGroupUserList(group *datamodels.Group, offset int, limit int) (users []*datamodels.User, err error) {
	return s.repo.GetUserList(group, offset, limit)
}

// 获取分组的权限列表
func (s *groupService) GetGroupPermissionList(group *datamodels.Group, offset int, limit int) (permissions []*datamodels.Permission, err error) {
	return s.repo.GetPermissionList(group, offset, limit)
}
