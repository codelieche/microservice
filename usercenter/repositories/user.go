package repositories

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/codelieche/microservice/usercenter/common"
	"github.com/codelieche/microservice/usercenter/datamodels"
	"github.com/jinzhu/gorm"
)

type UserRepository interface {
	// 保存User: 不会修改username和password
	Save(user *datamodels.User) (*datamodels.User, error)
	// 获取User的列表
	List(offset int, limit int) ([]*datamodels.User, error)
	// 获取User信息
	Get(id int64) (*datamodels.User, error)
	// 根据ID或者Name获取User信息
	GetByIdOrName(idOrName string) (*datamodels.User, error)
	// 获取User的用户列表
	GetUserList(user *datamodels.User, offset int, limit int) ([]*datamodels.User, error)
	// 设置用户密码
	SetUserPassword(user *datamodels.User, password string) (*datamodels.User, error)
	// 检查用户的密码
	CheckUserPassword(user *datamodels.User, password string) (bool, error)
	//	获取用户的分组列表
	GetUserGroups(user *datamodels.User) (groups []*datamodels.Group, err error)
	// 获取用户的角色列表
	GetUserRoles(user *datamodels.User) (roles []*datamodels.Role, err error)
	// 用户更新操作
	UpdateByID(id int64, fields map[string]interface{}) (*datamodels.User, error)
	// 创建管理员用户
	CheckOrCreateAdminUser() (user *datamodels.User, err error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db:               db,
		infoFields:       []string{"id", "created_at", "updated_at", "username", "email", "mobile", "is_active"},
		groupFields:      []string{"id", "created_at", "name", "user_id"},
		roleFields:       []string{"id", "created_at", "name", "user_id"},
		permissionFields: []string{"id", "created_at", "name", "code", "app_id", "user_id"},
	}
}

// 用户Repository
// 查询数据的时候，如果不指定Select，会Select *
type userRepository struct {
	db               *gorm.DB
	infoFields       []string // 基本信息字段
	groupFields      []string // 分组字段
	roleFields       []string // 角色字段
	permissionFields []string // 权限字段
}

// 检查用户的密码
func (r *userRepository) CheckUserPassword(user *datamodels.User, password string) (success bool, err error) {
	r.db.Select("id,username,password").First(user)
	if user.ID > 0 {
		// 检查用户的密码
		return user.CheckPassword(password)
	} else {
		return false, common.NotFountError
	}
}

// 保存User
func (r *userRepository) Save(user *datamodels.User) (*datamodels.User, error) {
	// 保存账号

	if user.ID > 0 {
		// 是更新操作
		user.Username = ""
		//user.Password = ""

		tx := r.db.Begin()
		// 判断是否需要更新Groups
		if len(user.Groups) > 0 {
			if err := tx.Model(user).Association("Groups").Replace(user.Groups).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		} else {
			// 删除掉所有的Groups
			if err := tx.Model(user).Association("Groups").Clear().Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		// 判断是否需要更新Roles
		if len(user.Roles) > 0 {
			if err := tx.Model(user).Association("Roles").Replace(user.Roles).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		} else {
			// 删除掉所有的Roles
			if err := tx.Model(user).Association("Roles").Clear().Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		// 是否需要更新权限
		if len(user.Permissions) > 0 {
			if err := tx.Model(user).Association("Permissions").Replace(user.Permissions).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		} else {
			// 删除掉所有的Permissions
			if err := tx.Model(user).Association("Permissions").Clear().Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		// 关于空值
		// r.db.Update: 不会更新空值
		// tx.Update()：是会更新空值的
		if err := tx.Model(&datamodels.User{}).Where("id=?", user.ID).Limit(1).
			Update(map[string]interface{}{"email": user.Email, "mobile": user.Mobile}).Error; err != nil {
			tx.Rollback()
			return nil, err
		} else {
			tx.Commit()
			//return user, nil
			return r.Get(int64(user.ID))
		}
	} else {
		// 是创建操作
		// 判断密码
		if user.Password != "" && len(user.Password) < 40 {
			// 密码不是加密了的，我们给它加密一下
			if err := user.SetPassword(user.Password); err != nil {
				log.Println(err.Error())
				return nil, err
			}
		}

		if err := r.db.Create(user).Error; err != nil {
			return nil, err
		} else {
			return user, nil
		}

	}
}

// 设置用户的密码
func (r *userRepository) SetUserPassword(user *datamodels.User, password string) (u *datamodels.User, err error) {
	if err = user.SetPassword(password); err != nil {
		return nil, err
	} else {
		if user.ID <= 0 {
			err = fmt.Errorf("传入的ID小于等于0")
			return nil, err
		}
		if err := r.db.Model(user).Where("id = ?", user.ID).Limit(1).Update("password", user.Password).Error; err != nil {
			return nil, err
		} else {
			return user, nil
		}
	}
}

// 获取用户的列表
func (r *userRepository) List(offset int, limit int) (users []*datamodels.User, err error) {
	query := r.db.Model(&datamodels.User{}).Select(r.infoFields).Offset(offset).Limit(limit).Find(&users)
	if query.Error != nil {
		return nil, query.Error
	} else {
		return users, nil
	}
}

// 根据ID获取User
func (r *userRepository) Get(id int64) (user *datamodels.User, err error) {
	user = &datamodels.User{}
	r.db.Select(r.infoFields).Preload("Groups", func(db *gorm.DB) *gorm.DB {
		return db.Select(r.groupFields)
	}).Preload("Roles", func(db *gorm.DB) *gorm.DB {
		return db.Select(r.roleFields)
	}).Preload("Permissions", func(db *gorm.DB) *gorm.DB {
		return db.Select(r.permissionFields)
	}).First(user, "id = ?", id)
	if user.ID > 0 {
		return user, nil
	} else {
		return nil, common.NotFountError
	}
}

// 根据ID/Name获取User
func (r *userRepository) GetByIdOrName(idOrName string) (user *datamodels.User, err error) {
	user = &datamodels.User{}
	r.db.Select(r.infoFields).Preload("Groups", func(db *gorm.DB) *gorm.DB {
		return db.Select(r.groupFields)
	}).Preload("Roles", func(db *gorm.DB) *gorm.DB {
		return db.Select(r.roleFields)
	}).Preload("Permissions", func(db *gorm.DB) *gorm.DB {
		return db.Select(r.permissionFields)
	}).First(user, "id = ? or username = ?", idOrName, idOrName)
	if user.ID > 0 {
		return user, nil
	} else {
		return nil, common.NotFountError
	}
}

// 获取User的用户
func (r *userRepository) GetUserList(
	user *datamodels.User, offset int, limit int) (users []*datamodels.User, err error) {
	query := r.db.Model(user).Select(r.infoFields).Offset(offset).Limit(limit).Related(&users, "Users")
	if query.Error != nil {
		return nil, query.Error
	} else {
		return users, nil
	}
}

// 获取用户的User
func (r *userRepository) GetUserGroups(user *datamodels.User) (groups []*datamodels.Group, err error) {
	query := r.db.Model(user).Select(r.groupFields).Related(&groups, "Groups")
	if query.Error != nil {
		return nil, query.Error
	} else {
		return groups, nil
	}
}

// 获取用户的角色
func (r *userRepository) GetUserRoles(user *datamodels.User) (roles []*datamodels.Role, err error) {
	query := r.db.Model(user).Select(r.roleFields).Related(&roles, "Roles")
	if query.Error != nil {
		return nil, query.Error
	} else {
		return roles, nil
	}
}

// 更新
func (r *userRepository) UpdateByID(id int64, fields map[string]interface{}) (user *datamodels.User, err error) {
	// 判断ID
	if id <= 0 {
		err = errors.New("传入ID为0,会更新全部数据")
		return nil, err
	}
	// 因为指定了ID了，所以这里可不判断这个ID
	// 丢弃ID/id/Id/iD
	//idKeys := []string{"ID", "id", "Id", "iD"}
	//for _, k := range idKeys {
	//	if _, exist := fields[k]; exist {
	//		delete(fields, k)
	//	}
	//}

	// 密码是不更新的: id也是不用传的
	// username也不可更新
	for key := range fields {
		if strings.ToLower(key) == "password" {
			delete(fields, key)
		}
		// username也不可更新
		if strings.ToLower(key) == "username" {
			delete(fields, key)
		}
	}

	// 更新操作
	if err = r.db.Model(&datamodels.User{}).Where("id = ?", id).Limit(1).Update(fields).Error; err != nil {
		return nil, err
	} else {
		return r.Get(id)
		//user = &datamodels.User{}
		//if err = r.db.First(user, "id = ?", id).Error; err != nil {
		//	return nil, err
		//} else {
		//	return user, nil
		//}
	}
}

// 创建超级用户
func (r *userRepository) CheckOrCreateAdminUser() (user *datamodels.User, err error) {
	var count int
	if err := r.db.Model(&datamodels.User{}).Count(&count).Error; err != nil {
		return nil, err
	} else {
		// 判断count
		if count < 1 {
			// 创建admin用户
			user := &datamodels.User{
				Username:    "admin",
				Password:    "changeme",
				Mobile:      "",
				Email:       "",
				IsSuperuser: true,
				IsActive:    true,
			}
			return r.Save(user)
		} else {
			// log.Println("当前用户个数为：", count)
			return nil, nil
		}
	}
}
