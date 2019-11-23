package repositories

import (
	"log"

	"github.com/codelieche/microservice/common"
	"github.com/codelieche/microservice/datamodels"
	"github.com/jinzhu/gorm"
)

type UserRepository interface {
	// 保存分组
	Save(user *datamodels.User) (*datamodels.User, error)
	// 设置用户密码
	SetUserPassword(user *datamodels.User, password string) (*datamodels.User, error)
	// 获取分组的列表
	List(offset int, limit int) ([]*datamodels.User, error)
	// 获取分组信息
	Get(id int64) (*datamodels.User, error)
	// 根据ID或者Name获取分组信息
	GetByIdOrName(idOrName string) (*datamodels.User, error)
	// 获取分组的用户列表
	GetUserList(user *datamodels.User, offset int, limit int) ([]*datamodels.User, error)
	//	获取用户的分组列表
	GetUserGroups(user *datamodels.User) (groups []*datamodels.Group, err error)
	GetUserRoles(user *datamodels.User) (roles []*datamodels.Role, err error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

type userRepository struct {
	db *gorm.DB
}

// 保存User
func (r *userRepository) Save(user *datamodels.User) (*datamodels.User, error) {
	// 判断密码
	if user.Password != "" && len(user.Password) < 40 {
		// 密码不是加密了的，我们给它加密一下
		if err := user.SetPassword(user.Password); err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}

	// 保存账号
	if user.ID > 0 {
		// 是更新操作
		if err := r.db.Model(&datamodels.User{}).Update(user).Error; err != nil {
			return nil, err
		} else {
			return user, nil
		}
	} else {
		// 是创建操作
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
		if user, err := r.Save(user); err != nil {
			return nil, err
		} else {
			return user, nil
		}
	}
}

// 获取用户的列表
func (r *userRepository) List(offset int, limit int) (users []*datamodels.User, err error) {
	query := r.db.Model(&datamodels.User{}).Offset(offset).Limit(limit).Find(&users)
	if query.Error != nil {
		return nil, query.Error
	} else {
		return users, nil
	}
	return
}

// 根据ID获取分组
func (r *userRepository) Get(id int64) (user *datamodels.User, err error) {

	user = &datamodels.User{}
	r.db.First(user, "id = ?", id)
	if user.ID > 0 {
		return user, nil
	} else {
		return nil, common.NotFountError
	}
}

// 根据ID/Name获取分组
func (r *userRepository) GetByIdOrName(idOrName string) (user *datamodels.User, err error) {

	user = &datamodels.User{}
	r.db.First(user, "id = ? or name = ?", idOrName, idOrName)
	if user.ID > 0 {
		return user, nil
	} else {
		return nil, common.NotFountError
	}
}

// 获取分组的用户
func (r *userRepository) GetUserList(
	user *datamodels.User, offset int, limit int) (users []*datamodels.User, err error) {
	query := r.db.Model(user).Offset(offset).Limit(limit).Related(&users, "Users")
	if query.Error != nil {
		return nil, query.Error
	} else {
		return users, nil
	}
}

// 获取用户的分组
func (r *userRepository) GetUserGroups(user *datamodels.User) (groups []*datamodels.Group, err error) {
	query := r.db.Model(user).Related(&groups, "Groups")
	if query.Error != nil {
		return nil, query.Error
	} else {
		return groups, nil
	}
}

// 获取用户的角色
func (r *userRepository) GetUserRoles(user *datamodels.User) (roles []*datamodels.Role, err error) {
	query := r.db.Model(user).Related(&roles, "Roles")
	if query.Error != nil {
		return nil, query.Error
	} else {
		return roles, nil
	}
}
