package services

import (
	"github.com/codelieche/microservice/usercenter/datamodels"
	"github.com/codelieche/microservice/usercenter/repositories"
)

type ProjectService interface {
	// 创建Project
	Create(project *datamodels.Project) (*datamodels.Project, error)
	// 保存Project
	Save(project *datamodels.Project) (*datamodels.Project, error)
	// 获取Project的列表
	List(offset int, limit int) ([]*datamodels.Project, int, error)
	// 获取Project的信息
	Get(id int64) (*datamodels.Project, error)
	// 根据ID或者Code获取Project信息
	GetByIdOrCode(idOrCode string) (project *datamodels.Project, err error)
	// 获取Project的权限列表
	GetPermissionList(project *datamodels.Project, offset int, limit int) ([]*datamodels.Permission, error)
	// 获取Project的用户列表
	GetUsersList(project *datamodels.Project, offset int, limit int) (users []*datamodels.ProjectUser, count int, err error)
	// 给Project添加用户
	AddProjectUser(project string, username string, role string) (projectUser *datamodels.ProjectUser, err error)
	// 给项目删除用户
	DeleteProjectUser(project string, username string) (success bool, err error)
}

func NewProjectService(repo repositories.ProjectRepository) ProjectService {
	return &projectService{
		repo: repo,
	}
}

type projectService struct {
	repo repositories.ProjectRepository
}

func (s *projectService) Create(project *datamodels.Project) (*datamodels.Project, error) {
	return s.repo.Save(project)
}

func (s *projectService) Save(project *datamodels.Project) (*datamodels.Project, error) {
	return s.repo.Save(project)
}

func (s *projectService) List(offset int, limit int) ([]*datamodels.Project, int, error) {
	return s.repo.List(offset, limit)
}

func (s *projectService) Get(id int64) (*datamodels.Project, error) {
	return s.repo.Get(id)
}

func (s *projectService) GetByIdOrCode(idOrCode string) (project *datamodels.Project, err error) {
	return s.repo.GetByIdOrCode(idOrCode)
}

func (s *projectService) GetPermissionList(project *datamodels.Project, offset int, limit int) ([]*datamodels.Permission, error) {
	return s.repo.GetPermissionList(project, offset, limit)
}

func (s *projectService) GetUsersList(project *datamodels.Project, offset int, limit int) (users []*datamodels.ProjectUser, count int, err error) {
	return s.repo.GetUsersList(project, offset, limit)
}

func (s *projectService) AddProjectUser(project string, username string, role string) (projectUser *datamodels.ProjectUser, err error) {
	return s.repo.AddProjectUser(project, username, role)
}

func (s *projectService) DeleteProjectUser(project string, username string) (success bool, err error) {
	return s.repo.DeleteProjectUser(project, username)
}
