package repositories

import (
	"github.com/codelieche/microservice/usercenter/common"
	"github.com/codelieche/microservice/usercenter/datamodels"
	"github.com/jinzhu/gorm"
)

type RoleRepository interface {
	// 保存Role
	Save(role *datamodels.Role) (*datamodels.Role, error)
	// 获取Role的列表
	List(offset int, limit int) ([]*datamodels.Role, error)
	// 获取Role信息
	Get(id int64) (*datamodels.Role, error)
	// 根据ID或者Name获取Role信息
	GetByIdOrName(idOrName string) (*datamodels.Role, error)
	// 获取Role的用户列表
	GetUserList(role *datamodels.Role, offset int, limit int) ([]*datamodels.User, error)
	// 获取Role的权限列表
	GetPermissionList(role *datamodels.Role, offset int, limit int) ([]*datamodels.Permission, error)
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

type roleRepository struct {
	db *gorm.DB
}

// 保存Role
func (r *roleRepository) Save(role *datamodels.Role) (*datamodels.Role, error) {
	if role.ID > 0 {
		// 是更新操作
		if err := r.db.Model(&datamodels.Role{}).Update(role).Error; err != nil {
			return nil, err
		} else {
			return role, nil
		}
	} else {
		// 是创建操作
		if err := r.db.Create(role).Error; err != nil {
			return nil, err
		} else {
			return role, nil
		}

	}
}

// 获取Role的列表
func (r *roleRepository) List(offset int, limit int) (roles []*datamodels.Role, err error) {
	query := r.db.Model(&datamodels.Role{}).Offset(offset).Limit(limit).Find(&roles)
	if query.Error != nil {
		return nil, query.Error
	} else {
		return roles, nil
	}
	return
}

// 根据ID获取Role
func (r *roleRepository) Get(id int64) (role *datamodels.Role, err error) {

	role = &datamodels.Role{}
	r.db.First(role, "id = ?", id)
	if role.ID > 0 {
		return role, nil
	} else {
		return nil, common.NotFountError
	}
}

// 根据ID/Name获取Role
func (r *roleRepository) GetByIdOrName(idOrName string) (role *datamodels.Role, err error) {

	role = &datamodels.Role{}
	r.db.First(role, "id = ? or name = ?", idOrName, idOrName)
	if role.ID > 0 {
		return role, nil
	} else {
		return nil, common.NotFountError
	}
}

// 获取Role的用户
func (r *roleRepository) GetUserList(
	role *datamodels.Role, offset int, limit int) (users []*datamodels.User, err error) {
	query := r.db.Model(role).Offset(offset).Limit(limit).Related(&users, "Users")
	if query.Error != nil {
		return nil, query.Error
	} else {
		return users, nil
	}
}

// 获取角色的权限列表
func (r *roleRepository) GetPermissionList(role *datamodels.Role, offset int, limit int) (permissions []*datamodels.Permission, err error) {
	query := r.db.Model(role).Offset(offset).Limit(limit).Related(&permissions, "Permissions")
	if query.Error != nil {
		return nil, query.Error
	} else {
		return permissions, nil
	}
}
