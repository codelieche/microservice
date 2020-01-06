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
		userFields:       []string{"id", "created_at", "updated_at", "username", "email", "mobile", "is_active", "group_id"},
		permissionFields: []string{"id", "created_at", "updated_at", "name", "app_id", "code", "group_id"},
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
		// 更新的话用事务的方式来处理
		// 当编辑的时候，如果Users、Permissions为空，那么也需要处理
		tx := r.db.Begin()
		// 1. 更新users
		if len(group.Users) > 0 {
			users := group.Users
			if err := tx.Model(group).Association("Users").Replace(users).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		} else {
			// 删除掉groups的所有Users
			//log.Println("需要对用户置空")
			if err := tx.Model(group).Association("Users").Clear().Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
		// 2. 更新Permissions
		if len(group.Permissions) > 0 {
			permissions := group.Permissions
			if err := tx.Model(group).Association("Permissions").Replace(permissions).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		} else {
			// 需要把permissions置空
			if err := tx.Model(group).Association("Permissions").Clear().Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		// 3. 更新group
		group.Users = nil
		if err := tx.Model(&datamodels.Group{}).Update(group).Error; err != nil {
			tx.Rollback()
			return nil, err
		} else {
			tx.Commit()
			return r.Get(int64(group.ID))
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
}

// 根据ID获取分组
func (r *groupRepository) Get(id int64) (group *datamodels.Group, err error) {
	group = &datamodels.Group{}
	r.db.Model(group).Preload("Users", func(db *gorm.DB) *gorm.DB {
		return db.Select(r.userFields)
	}).Preload("Permissions", func(db *gorm.DB) *gorm.DB {
		return db.Select(r.permissionFields)
	}).Select(r.infoFields).First(group, "id = ?", id)

	if group != nil && group.ID > 0 {
		return group, nil
	} else {
		return nil, common.NotFountError
	}
}

// 根据ID/Name获取分组
func (r *groupRepository) GetByIdOrName(idOrName string) (group *datamodels.Group, err error) {

	group = &datamodels.Group{}
	r.db.Model(group).Preload("Users", func(db *gorm.DB) *gorm.DB {
		return db.Select(r.userFields)
	}).
		Preload("Permissions", func(db *gorm.DB) *gorm.DB {
			return db.Select(r.permissionFields)
		}).
		Select(r.infoFields).First(group, "id = ? or name = ?", idOrName, idOrName)
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
