package core

import (
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Team{})
	db.AutoMigrate(&TeamUser{})

	//user := User{
	//	Username: "codelieche",
	//}
	//if err := db.Create(&user).Error; err != nil {
	//	log.Printf("出现错误：%s", err.Error())
	//}
}
