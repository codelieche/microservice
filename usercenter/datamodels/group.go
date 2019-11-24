package datamodels

import "github.com/jinzhu/gorm"

type Group struct {
	gorm.Model
	Name        string        `gorm:"type:varchar(100);UNIQUE_INDEX"`
	Users       []*User       `gorm:"many2many:group_users" json:"users,omitempty"`
	Permissions []*Permission `gorm:"many2many:group_permissions" json:"permissions,omitempty"`
}
