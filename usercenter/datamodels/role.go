package datamodels

import "github.com/jinzhu/gorm"

type Role struct {
	gorm.Model
	//AppID       uint          `gorm:"UNIQUE_INDEX:uidx_app_id_name"`                           // 应用的角色
	Name        string        `gorm:"type:varchar(40);NOT NULL;UNIQUE_INDEX" json:"name"`      // 角色名称
	Users       []*User       `gorm:"many2many:role_users" json:"users,omitempty"`             // 角色用户
	Permissions []*Permission `gorm:"many2many:role_permissions" json:"permissions,omitempty"` // 角色的权限
}
