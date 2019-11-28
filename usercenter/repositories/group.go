package repositories

import (
	"github.com/codelieche/microservice/usercenter/common"
	"github.com/codelieche/microservice/usercenter/datamodels"
	"github.com/jinzhu/gorm"
)

type GroupRepository interface {
	// 保存分组
	Save(group *datamodels.Group) (*datamodels.Group, error)
	// 获取分组的列表
	List(offset int, limit int) ([]*datamodels.Group, error)
	// 获取分组信息
	Get(id int64) (*datamodels.Group, error)
	// 根据ID或者Name获取分组信息
	GetByIdOrName(idOrName string) (*datamodels.Group, error)
	// 获取分组的用户列表
	GetUserList(group *datamodels.Group, offset int, limit int) ([]*datamodels.User, error)
	// 获取分组的权限列表
	GetPermissionList(group *datamodels.Group, offset int, limit int) ([]*datamodels.Permission, error)
}

func NewGroupRepository(db *gorm.DB) GroupRepository {
	return &groupRepository{
		db:               db,
		infoFields:       []string{"id", "created_at", "updated_at", "name"},
		userFields:       []string{"id", "created_at", "updated_at", "username", "email", "mobile", "is_active"},
		permissionFields: []string{"id", "created_at", "updated_at", "app_id", "code"},
	}
}

type groupRepository struct {
	db               *gorm.DB
	infoFields       []string // 分组信息字段
	userFields       []string // 分组用户信息字段
	permissionFields []string // 分组权限字段
}

// 保存Group
func (r *groupRepository) Save(group *datamodels.Group) (*datamodels.Group, error) {
	if group.ID > 0 {
		// 是更新操作
		if err := r.db.Model(&datamodels.Group{}).Update(group).Error; err != nil {
			return nil, err
		} else {
			return group, nil
		}
	} else {
		// 是创建操作
		if err := r.db.Create(group).Error; err != nil {
			return nil, err
		} else {
			return group, nil
		}

	}
}

// 获取Group的列表
func (r *groupRepository) List(offset int, limit int) (groups []*datamodels.Group, err error) {
	query := r.db.Model(&datamodels.Group{}).Select(r.infoFields).Offset(offset).Limit(limit).Find(&groups)
	if query.Error != nil {
		return nil, query.Error
	} else {
		return groups, nil
	}
	return
}

// 根据ID获取分组
func (r *groupRepository) Get(id int64) (group *datamodels.Group, err error) {

	group = &datamodels.Group{}
	r.db.Select(r.infoFields).First(group, "id = ?", id)
	if group.ID > 0 {
		return group, nil
	} else {
		return nil, common.NotFountError
	}
}

// 根据ID/Name获取分组
func (r *groupRepository) GetByIdOrName(idOrName string) (group *datamodels.Group, err error) {

	group = &datamodels.Group{}
	r.db.Select(r.infoFields).First(group, "id = ? or name = ?", idOrName, idOrName)
	if group.ID > 0 {
		return group, nil
	} else {
		return nil, common.NotFountError
	}
}

// 获取分组的用户
func (r *groupRepository) GetUserList(
	group *datamodels.Group, offset int, limit int) (users []*datamodels.User, err error) {
	query := r.db.Model(group).Select(r.userFields).Offset(offset).Limit(limit).Related(&users, "Users")
	if query.Error != nil {
		return nil, query.Error
	} else {
		return users, nil
	}
}

// 获取分组的权限列表
func (r *groupRepository) GetPermissionList(group *datamodels.Group, offset int, limit int) (permissions []*datamodels.Permission, err error) {
	query := r.db.Model(group).Select(r.permissionFields).Preload("Application", func(d *gorm.DB) *gorm.DB {
		return d.Select("id, name, code")
	}).Offset(offset).Limit(limit).Related(&permissions, "Permissions")
	if query.Error != nil {
		return nil, query.Error
	} else {
		return permissions, nil
	}
}
