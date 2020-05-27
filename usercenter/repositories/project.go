package repositories

import (
	"errors"
	"fmt"
	"log"

	"github.com/codelieche/microservice/usercenter/common"
	"github.com/codelieche/microservice/usercenter/datamodels"
	"github.com/jinzhu/gorm"
)

type ProjectRepository interface {
	// 保存项目
	Save(project *datamodels.Project) (*datamodels.Project, error)
	// 获取Project的列表
	List(offset int, limit int) (projects []*datamodels.Project, count int, err error)
	// 获取Project信息
	Get(id int64) (project *datamodels.Project, err error)
	// 根据ID或者Code获取Project信息
	GetByIdOrCode(idOrCode string) (project *datamodels.Project, err error)
	// 获取Project的所有Permission列表
	GetPermissionList(project *datamodels.Project, offset int, limit int) ([]*datamodels.Permission, error)
	//	获取Project的所有用户列表
	GetUsersList(project *datamodels.Project, offset int, limit int) (users []*datamodels.ProjectUser, count int, err error)
	// 给Project添加用户
	AddProjectUser(project string, username string, role string) (projectUser *datamodels.ProjectUser, err error)
	// 给项目删除用户
	DeleteProjectUser(project string, username string) (success bool, err error)
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{
		db:               db,
		infoFields:       []string{"id", "created_at", "updated_at", "name", "code", "description"},
		userFields:       []string{"id", "role", "project", "user"},
		permissionFiedls: []string{"id", "created_at", "name", "code", "app_id"},
	}
}

type projectRepository struct {
	db               *gorm.DB
	infoFields       []string
	userFields       []string
	permissionFiedls []string
}

func (r *projectRepository) Save(project *datamodels.Project) (*datamodels.Project, error) {
	if project.Name == "" {
		err := errors.New("name不可为空")
		return nil, err
	}
	if project.Code == "" {
		err := errors.New("code不可为空")
		return nil, err
	}
	if project.ID > 0 {
		// 是更新操作
		tx := r.db.Begin()
		// 修改app的信息
		if err := tx.Model(&datamodels.Project{}).Update(project).Error; err != nil {
			tx.Rollback()
			return nil, err
		} else {
			// 提交事务
			tx.Commit()
			return r.Get(int64(project.ID))
		}
	} else {
		// 是创建操作
		// 随机生成Token
		if err := r.db.Create(project).Error; err != nil {
			return nil, err
		} else {
			return project, nil
		}
	}
}

// 获取Project的列表
func (r *projectRepository) List(offset int, limit int) (projects []*datamodels.Project, count int, err error) {
	//r.db.Model(&datamodels.Project{}).Count(&count)
	query := r.db.Model(&datamodels.Project{}).
		Select(r.infoFields).
		Count(&count).
		//Preload("Users", func(db *gorm.DB) *gorm.DB {
		//	return db.Select(r.userFields)
		//}).
		Offset(offset).Limit(limit).
		Find(&projects)
	if query.Error != nil {
		return nil, 0, query.Error
	} else {
		return projects, count, nil
	}
}

// 根据ID获取Project
func (r *projectRepository) Get(id int64) (project *datamodels.Project, err error) {
	project = &datamodels.Project{}
	r.db.Preload("Users", func(db *gorm.DB) *gorm.DB {
		return db.Select(r.userFields)
	}).First(project, "id = ?", id)

	if project.ID > 0 {
		return project, nil
	} else {
		return nil, common.NotFountError
	}
}

// 根据ID/Code获取Project
func (r *projectRepository) GetByIdOrCode(idOrCode string) (project *datamodels.Project, err error) {
	project = &datamodels.Project{}
	r.db.Preload("Users", func(db *gorm.DB) *gorm.DB {
		return db.Select(r.userFields)
	}).First(project, "id = ? or code = ?", idOrCode, idOrCode)

	if project.ID > 0 {
		return project, nil
	} else {
		return nil, common.NotFountError
	}
}

// 获取权限列表
func (r *projectRepository) GetPermissionList(project *datamodels.Project, offset int, limit int) (permissions []*datamodels.Permission, err error) {
	query := r.db.Model(project).Offset(offset).Limit(limit).Related(&permissions, "Permissions")
	if query.Error != nil {
		return nil, query.Error
	} else {
		return permissions, nil
	}
}

// 获取项目的用户
func (r *projectRepository) GetUsersList(project *datamodels.Project, offset int, limit int) (users []*datamodels.ProjectUser, count int, err error) {
	// 获取项目用户的个数
	r.db.Model(&datamodels.ProjectUser{}).Select("id").Where("project = ?", project.Code).Count(&count)
	query := r.db.Model(project).
		Offset(offset).Limit(limit).
		Related(&users, "Users")
	if query.Error != nil {
		return nil, 0, query.Error
	} else {
		return users, count, nil
	}
}

// 给项目添加用户
func (r *projectRepository) AddProjectUser(projectCode string, username string, role string) (projectUser *datamodels.ProjectUser, err error) {
	// 1. 判断项目和用户是否存在
	var project = &datamodels.Project{}
	var user = &datamodels.User{}
	if project, err = r.GetByIdOrCode(projectCode); err != nil {
		message := fmt.Sprintf("%s项目不存在", projectCode)
		return nil, errors.New(message)
	} else {
		if project.ID <= 0 {
			message := fmt.Sprintf("%s项目不存在", projectCode)
			return nil, errors.New(message)
		}
	}

	if err = r.db.Model(&datamodels.User{}).First(&user, "username = ?", username).Error; err != nil {
		log.Print(err.Error())
		message := fmt.Sprintf("%s用户不存在", username)
		return nil, errors.New(message)
	} else {
		if user.ID <= 0 {
			message := fmt.Sprintf("%s用户不存在", username)
			return nil, errors.New(message)
		}
	}

	// 2. 开始添加用户
	// log.Println(project.ID, project.Code, user.ID, user.Username)
	projectUser = &datamodels.ProjectUser{}
	err = r.db.FirstOrCreate(
		&projectUser,
		&datamodels.ProjectUser{Project: project.Code, User: user.Username}).
		Error

	if err != nil {
		return nil, err
	} else {
		// 设置角色
		updatedFields := map[string]interface{}{}
		needUpdated := false
		if projectUser.Role != role {
			needUpdated = true
			updatedFields["role"] = role
		}
		if projectUser.IsDeleted {
			needUpdated = true
			updatedFields["IsDeleted"] = false
		}
		// 保存用户角色
		if needUpdated {
			r.db.Model(projectUser).Limit(1).Update(updatedFields)
		}
		return projectUser, nil
	}

}

func (r *projectRepository) DeleteProjectUser(project string, username string) (success bool, err error) {
	pUser := &datamodels.ProjectUser{}

	if err = r.db.Model(&datamodels.ProjectUser{}).
		Select(r.userFields).
		First(&pUser, "project = ? and user = ?", project, username).Error; err != nil {
		return false, err
	}

	if pUser.ID > 0 {
		if pUser.IsDeleted {
			return true, nil
		} else {
			if err = r.db.Model(&pUser).Update("IsDeleted", true).Error; err != nil {
				return false, err
			} else {
				return true, nil
			}
		}
	} else {
		return false, common.NotFountError
	}
}
