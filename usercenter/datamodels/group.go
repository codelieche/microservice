package datamodels

import "github.com/jinzhu/gorm"

type Group struct {
	gorm.Model
	Name        string        `gorm:"type:varchar(100);UNIQUE_INDEX"`
	Permissions []*Permission `gorm:"many2many:group_permissions" json:"permissions,omitempty"`
}
