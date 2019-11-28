package repositories

import (
	"fmt"
	"log"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/codelieche/microservice/usercenter/datamodels"

	"github.com/codelieche/microservice/usercenter/datasources"
)

func TestUserRepository_Save(t *testing.T) {
	// 1. get db
	db := datasources.GetDb()

	// 2. init user repository
	r := NewUserRepository(db)

	// 3. 测试插入用户
	i := 0
	for i < 10 {
		i++
		name := fmt.Sprintf("user%d", i)
		password := fmt.Sprintf("password%d", i)
		user := &datamodels.User{
			Username: name,
			Password: password,
		}

		if u, err := r.Save(user); err != nil {
			log.Println(err.Error())
		} else {
			log.Println(u.ID, u.Username)
		}
	}
}

func TestNewUserReposity_List(t *testing.T) {
	// 1. get db
	db := datasources.GetDb()

	// 2. init user repository
	r := NewUserRepository(db)

	// 4. 测试获取用户列表
	haveNext := true
	offset := 0
	limit := 5
	for haveNext {
		if users, err := r.List(offset, limit); err != nil {
			t.Error(err.Error())
			haveNext = false
		} else {
			if len(users) == limit && limit > 0 {
				haveNext = true
				offset += limit
			} else {
				haveNext = false
			}
			// 输出User
			for _, user := range users {
				log.Println(user.ID, user.Username)
			}
		}
	}
}

func TestPasswordEncrypted(t *testing.T) {
	pwd := "This Is Password"
	if hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pwd), 8); err != nil {
		t.Error(err.Error())
	} else {
		log.Println(string(hashedPassword))
		log.Println("$2a$08$nW5Pmt1.5/qliuCOm0EoJO8IxFgj2aUMuw.M9jkdo8EEyoXC9eZ1K")
	}
}

func TestCompareHashAndPassword(t *testing.T) {
	pwd := "This Is Password"
	hasPassword := "$2a$08$0DmNRgrXCg1hFyMycVzQ7.eAQtm6PWfp5Rh9gEGmyKatGvYZfw8dK"

	if err := bcrypt.CompareHashAndPassword([]byte(hasPassword), []byte(pwd)); err != nil {
		t.Error(err.Error())
		//	crypto/bcrypt: hashedPassword is not the hash of the given password
	} else {
		log.Println("Password Is Ok")
	}
}

func TestUserRepository_GetUserGroups(t *testing.T) {
	// 1. get db
	db := datasources.GetDb()

	// 2. init user repository
	r := NewUserRepository(db)
	var userId int64 = 1
	// 获取用户
	if user, err := r.Get(userId); err != nil {
		t.Error(err.Error())
		return
	} else {
		// 3. 获取用户分组
		if groups, err := r.GetUserGroups(user); err != nil {
			t.Error(err.Error())
		} else {
			// 4. 打印出分组信息
			for _, group := range groups {
				log.Println(group.ID, group.Name)
			}
		}
	}
}

func TestUserRepository_GetUserRoles(t *testing.T) {
	// 1. get db
	db := datasources.GetDb()

	// 2. init user repository
	r := NewUserRepository(db)
	var userId int64 = 1
	// 获取用户
	if user, err := r.Get(userId); err != nil {
		t.Error(err.Error())
		return
	} else {
		// 3. 获取用户角色
		if roles, err := r.GetUserRoles(user); err != nil {
			t.Error(err.Error())
		} else {
			// 4. 打印出分组信息
			for _, role := range roles {
				log.Println(role.ID, role.Name)
			}
		}
	}
}
