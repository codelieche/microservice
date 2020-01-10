package repositories

import (
	"errors"
	"log"

	"github.com/codelieche/microservice/usercenter/common"
	"github.com/codelieche/microservice/usercenter/datamodels"
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
	GetAllPermissionByUserID(id int64) (permissions []*datamodels.Permission, err error)
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{
		db:         db,
		infoFields: []string{"id", "created_at", "updated_at", "name", "code", "app_id"},
		appFields:  []string{"id", "created_at", "updated_at", "name", "code", "info"},
	}
}

type permissionRepository struct {
	db         *gorm.DB
	infoFields []string
	appFields  []string
}

// 保存Permission
func (r *permissionRepository) Save(permission *datamodels.Permission) (*datamodels.Permission, error) {

	if permission.Name == "" {
		err := errors.New("name不可为空")
		return nil, err
	}
	if permission.Code == "" {
		err := errors.New("code不可为空")
		return nil, err
	}

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
	query := r.db.Model(&datamodels.Permission{}).Select(r.infoFields).
		Offset(offset).Limit(limit).Find(&permissions)
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
	r.db.Preload("Application", func(db *gorm.DB) *gorm.DB {
		return db.Select(r.appFields)
	}).Select(r.infoFields).First(permission, "id = ?", id)
	if permission.ID > 0 {
		return permission, nil
	} else {
		return nil, common.NotFountError
	}
}

// 根据ID/Name获取权限
func (r *permissionRepository) GetByAppIDAndCode(appID int, code string) (permission *datamodels.Permission, err error) {

	permission = &datamodels.Permission{}
	r.db.Preload("Application", func(db *gorm.DB) *gorm.DB {
		return db.Select(r.appFields)
	}).Select(r.infoFields).First(permission, "app_id = ? and code = ?", appID, code)
	if permission.ID > 0 {
		return permission, nil
	} else {
		return nil, common.NotFountError
	}
}

// 获取权限的用户
func (r *permissionRepository) GetUserList(
	permission *datamodels.Permission, offset int, limit int) (users []*datamodels.User, err error) {
	query := r.db.Model(permission).Select(r.infoFields).Offset(offset).Limit(limit).Related(&users, "Users")
	if query.Error != nil {
		return nil, query.Error
	} else {
		return users, nil
	}
}

// 获取用户的所有权限
func (r *permissionRepository) GetAllPermissionByUserID(id int64) (permissions []*datamodels.Permission, err error) {
	permissions = []*datamodels.Permission{}

	// 查询permissions的条件语句
	sql := `ID in (
SELECT permission_id from role_permissions 
	WHERE role_id in 
	(SELECT role_id FROM role_users WHERE user_id=?)
UNION

SELECT permission_id from group_permissions 
	WHERE group_id in 
	(SELECT group_id FROM group_users WHERE user_id=?)
UNION
SELECT permission_id from user_permissions WHERE user_id = ?
)`

	if err = r.db.Model(&datamodels.Permission{}).Select(r.infoFields).Where(sql, id, id, id).
		Find(&permissions).Error; err != nil {
		log.Println(err)
		return nil, err
	} else {
		return permissions, nil
	}
}
