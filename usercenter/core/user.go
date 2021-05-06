package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/codelieche/microservice/codelieche/utils"
	"gorm.io/gorm"
	"log"
	"regexp"
)

type (
	// User 用户Model username要以字符开头，且长度是6+(小写字母和数字和-)
	User struct {
		gorm.Model
		Username    string `gorm:"size:60;unique" json:"username" form:"username"`                     // 用户名，英文名小写+数字，且唯一
		Nickname    string `gorm:"size:60" json:"nickname" form:"nickname"`                            // 昵称-中文名
		Password    string `gorm:"size:256;" json:"-" form:"password"`                                 // 密码
		Phone       string `gorm:"size:40" json:"phone" form:"phone"`                                  // 手机号，为空或者unique
		Email       string `gorm:"size:100" json:"email" form:"email"`                                 // 邮箱，为空或者unique
		IsSuperuser bool   `gorm:"type:boolean;default:false" json:"is_superuser" form:"is_superuser"` // 是否超级用户
		IsActive    bool   `gorm:"type:boolean;default:true" json:"is_active" form:"is_active"`        // 是否启用
		Avatar      string `gorm:"size:512" json:"avatar" form:"avatar"`                               // 头像地址
		Signature   string `gorm:"size:256;default:null" json:"signature" form:"signature"`            // 个性签名
	}

	// UserStore 用户Store
	UserStore interface {
		// Find returns a user from the database
		Find(context.Context, int64) (*User, error)

		//	FindByUsername returns a user from the database by username
		FindByUsername(context.Context, string) (*User, error)

		// Create persists a new user to the database
		Create(context.Context, *User) (*User, error)

		// Update persists an updated user to the database
		Update(context.Context, *User) error

		// Delete deletes a user from the database
		Delete(context.Context, *User) error

		// Count returns 用户的数量
		Count(context.Context) (int64, error)

		// CountByTeam 根据团队获取总共的用户数量
		CountByTeam(context.Context, int64) (int64, error)
	}

	// UserService 用户服务接口
	UserService interface {
		// Find 查找用户
		Find(context.Context, int64) (*User, error)

		// Create 创建用户
		Create(context.Context, *User) (*User, error)
	}
)

func (u *User) ValidateUsername() (bool, error) {
	// 检查用户的用户名是否存在，检查用户的用户名是否符合规范
	reg := regexp.MustCompile(`^[a-z][a-z0-9\-]{5,59}$`)
	if reg.MatchString(u.Username) {
		return true, nil
	} else {
		err := fmt.Errorf("用户名需要是字母开头，由小写字母/-/数字组成，长度6-60")
		return false, err
	}
}

// BeforeCreate 创建用户之前先检查用户是否已经存在
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	// 1. 检查用户名是否合法
	if isOk, err := u.ValidateUsername(); err != nil {
		return err
	} else {
		if !isOk {
			msg := fmt.Sprintf("用户名不合法:%s", u.Username)
			return errors.New(msg)
		}
	}

	// 2. 检查用户是否已经存在
	var count int64
	// 2-1: 根据username判断
	if err := tx.Model(u).Find(nil, "username=?", u.Username).Count(&count).Error; err != nil {
		return err
	} else {
		if count > 0 {
			// 不用等数据库报错，我们这里直接判断重复了
			msg := fmt.Sprintf("用户%s已经存在", u.Username)
			return errors.New(msg)
		}
	}
	// 2-2：根据手机号判断
	if u.Phone != "" {
		if err := tx.Model(u).Find(nil, "phone=?", u.Phone).Count(&count).Error; err != nil {
			return err
		} else {
			if count > 0 {
				// 不用等数据库报错，我们这里直接判断重复了
				msg := fmt.Sprintf("手机号%s已经注册", u.Phone)
				return errors.New(msg)
			}
		}
	}
	// 2-3：根据邮箱判断
	if u.Email != "" {
		if err := tx.Model(u).Find(nil, "email=?", u.Email).Count(&count).Error; err != nil {
			return err
		} else {
			if count > 0 {
				// 不用等数据库报错，我们这里直接判断重复了
				msg := fmt.Sprintf("邮箱%s已经注册", u.Email)
				return errors.New(msg)
			}
		}
	}

	// 3. 初始化密码
	if u.Password == "" {
		randPassword := utils.RandomString(16)
		if hashedPassword, err := utils.HashPassword(randPassword); err != nil {
			u.Password = ""
		} else {
			u.Password = hashedPassword
			log.Printf("即将创建用户:%s, 未配置密码，我们初始化密码为：%s\n", u.Username, randPassword)
		}
	}
	// 原始密码，我们只配置最长40位的, 加密后一般是60位长度
	if u.Password != "" && len(u.Password) <= 40 {
		if hashedPassword, err := utils.HashPassword(u.Password); err == nil {
			u.Password = hashedPassword
		}
	}

	// 4. 如果昵称为空就设置其未username
	if u.Nickname == "" {
		u.Nickname = u.Username
	}
	return nil
}
