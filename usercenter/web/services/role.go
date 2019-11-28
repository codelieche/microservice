package services

import (
	"github.com/codelieche/microservice/usercenter/datamodels"
	"github.com/codelieche/microservice/usercenter/repositories"
)

// Role Service Interface
type RoleService interface {
	GetById(id int64) (role *datamodels.Role, err error)
	GetByIdOrName(idOrName string) (role *datamodels.Role, err error)
	List(offset int, limit int) (roles []*datamodels.Role, err error)
	GetRoleUserList(role *datamodels.Role, offset int, limit int) (users []*datamodels.User, err error)
	GetRolePermissionList(role *datamodels.Role, offset int, limit int) (permissions []*datamodels.Permission, err error)
}

// 实例化Role Service
func NewRoleService(repo repositories.RoleRepository) RoleService {
	return &roleService{repo: repo}
}

// role Service
type roleService struct {
	repo repositories.RoleRepository
}

func (s *roleService) GetById(id int64) (role *datamodels.Role, err error) {
	return s.repo.Get(id)
}

func (s *roleService) GetByIdOrName(idOrName string) (role *datamodels.Role, err error) {
	return s.repo.GetByIdOrName(idOrName)
}

// 获取用户角色列表
func (s *roleService) List(offset int, limit int) (roles []*datamodels.Role, err error) {
	return s.repo.List(offset, limit)
}

// 获取角色的用户列表
func (s *roleService) GetRoleUserList(role *datamodels.Role, offset int, limit int) (users []*datamodels.User, err error) {
	return s.repo.GetUserList(role, offset, limit)
}

// 获取角色的权限列表
func (s *roleService) GetRolePermissionList(role *datamodels.Role, offset int, limit int) (permissions []*datamodels.Permission, err error) {
	return s.repo.GetPermissionList(role, offset, limit)
}
