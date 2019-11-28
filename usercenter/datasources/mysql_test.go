package datasources

import (
	"log"
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/codelieche/microservice/usercenter/datamodels"
	"github.com/codelieche/microservice/usercenter/web/forms"
)

func TestInitDb(t *testing.T) {
	initDb()
}

func TestGetUserModel(t *testing.T) {
	// 1. get db
	db := GetDb()

	// 2. 查询用户
	userID := 1

	user := datamodels.User{}
	query := db.LogMode(true).Select("id, username, password").Preload("Groups", func(d *gorm.DB) *gorm.DB {
		return d.Select("*")
	}).First(&user, "id = ?", userID)
	if query.Error != nil {
		log.Println(query.Error)
	} else {
		log.Println(user)
		log.Println(user.Groups)
		for _, g := range user.Groups {
			log.Println(g.Name, g.ID)
		}
	}

	log.Println("======||||=====")

	user1 := forms.UserDetailForm{}
	query = db.Select("id, username, password").First(&datamodels.User{}, "id = ?", userID).Scan(&user1)
	if query.Error != nil {
		log.Println(query.Error)
	} else {
		log.Println(user1)
	}

	log.Println("========")
	user2 := forms.UserDetailForm{}
	if err := db.Table("users").
		Select("id,username,email,mobile").
		Where("id = ?", userID+10).Scan(&user2).Error; err != nil {
		log.Println(err)
	} else {
		log.Println(user2)
	}

	log.Println("==== User03 ====")

	user3 := datamodels.User{}
	if err := db.Table("users").Select("id,username,email,mobile").
		Where("id = ?", userID+10).Find(&user3).Error; err != nil {
		log.Println(err)
	} else {
		log.Println(user3)
		log.Println(user3.Username)
	}
}
