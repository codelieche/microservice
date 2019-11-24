package repositories

import (
	"github.com/codelieche/microservice/common"
	"github.com/codelieche/microservice/datamodels"
	"github.com/jinzhu/gorm"
)

type PermissionRepository interface {
	// 保存权限
	Save(permission *datamodels.Permission) (*datamodels.Permission, error)
	// 获取权限的列表
	List(offset int, limit int) ([]*datamodels.Permission, error)
	// 获取权限信息
	Get(id int64) (*datamodels.Permission, error)
	// 根据AppID或者Code获取权限信息
	GetByAppIDAndCode(appId int, code string) (*datamodels.Permission, error)
	// 获取权限的用户列表
	GetUserList(permission *datamodels.Permission, offset int, limit int) ([]*datamodels.User, error)
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

type permissionRepository struct {
	db *gorm.DB
}

// 保存Permission
func (r *permissionRepository) Save(permission *datamodels.Permission) (*datamodels.Permission, error) {
	if permission.ID > 0 {
		// 是更新操作
		if err := r.db.Model(&datamodels.Permission{}).Update(permission).Error; err != nil {
			return nil, err
		} else {
			return permission, nil
		}
	} else {
		// 是创建操作
		if err := r.db.Create(permission).Error; err != nil {
			return nil, err
		} else {
			return permission, nil
		}

	}
}

// 获取Permission的列表
func (r *permissionRepository) List(offset int, limit int) (permissions []*datamodels.Permission, err error) {
	query := r.db.Model(&datamodels.Permission{}).Offset(offset).Limit(limit).Find(&permissions)
	if query.Error != nil {
		return nil, query.Error
	} else {
		return permissions, nil
	}
	return
}

// 根据ID获取权限
func (r *permissionRepository) Get(id int64) (permission *datamodels.Permission, err error) {

	permission = &datamodels.Permission{}
	r.db.First(permission, "id = ?", id)
	if permission.ID > 0 {
		return permission, nil
	} else {
		return nil, common.NotFountError
	}
}

// 根据ID/Name获取权限
func (r *permissionRepository) GetByAppIDAndCode(appID int, code string) (permission *datamodels.Permission, err error) {

	permission = &datamodels.Permission{}
	r.db.First(permission, "app_id = ? and code = ?", appID, code)
	if permission.ID > 0 {
		return permission, nil
	} else {
		return nil, common.NotFountError
	}
}

// 获取权限的用户
func (r *permissionRepository) GetUserList(
	permission *datamodels.Permission, offset int, limit int) (users []*datamodels.User, err error) {
	query := r.db.Model(permission).Offset(offset).Limit(limit).Related(&users, "Users")
	if query.Error != nil {
		return nil, query.Error
	} else {
		return users, nil
	}
}