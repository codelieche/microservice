package datamodels

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	//gorm.Model
	BaseFields
	Username    string        `gorm:"type:varchar(40);NOT NULL;UNIQUE_INDEX" json:"username"`  // 用户名
	Password    string        `gorm:"type:varchar(256);NULL" json:"-"`                         // 用户密码
	Mobile      string        `gorm:"size:100" json:"mobile"`                                  // 用户手机号
	Email       string        `gorm:"size:100" json:"email"`                                   // 用户邮箱
	IsSuperuser bool          `gorm:"type:boolean" json:"is_superuser"`                        // 是否是超级用户
	IsActive    bool          `gorm:"type:boolean" json:"is_active"`                           // 是否启用
	Groups      []*Group      `gorm:"many2many:group_users" json:"groups,omitempty"`           // 用户分组
	Roles       []*Role       `gorm:"many2many:role_users" json:"roles,omitempty"`             // 用户角色
	Permissions []*Permission `gorm:"many2many:user_permissions" json:"permissions,omitempty"` // 用户权限
}

// 设置用户的密码
func (u *User) SetPassword(password string) (err error) {
	password = strings.TrimSpace(password)
	if password == "" {
		err = errors.New("密码不可为空")
		return err
	} else {
		// 判断密码长度
		length := len(password)
		if length < 8 {
			msg := fmt.Sprintf("密码最少长度为8, 当前长度为%d", length)
			err = errors.New(msg)
			return err
		}
		if length > 40 {
			msg := fmt.Sprintf("密码最大长度为40, 当前长度为%d", length)
			err = errors.New(msg)
			return err
		}
	}

	// 加密密码
	if hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8); err != nil {
		return err
	} else {
		u.Password = string(hashedPassword)
		return nil
	}
}

// 检查用户的密码
func (u *User) CheckPassword(password string) (success bool, err error) {
	password = strings.TrimSpace(password)
	if password == "" {
		err = errors.New("传入的密码为空")
		return false, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		//	crypto/bcrypt: hashedPassword is not the hash of the given password
		return false, err
	} else {
		return true, nil
	}
}
