package datamodels

import "github.com/jinzhu/gorm"

type Permission struct {
	gorm.Model
	Name        string      `gorm:"type:varchar(100)" json:"name"`                              // 权限的名称
	Code        string      `gorm:"size:100;UNIQUE_INDEX:idx_app_id_code;NOT NULL" json:"code"` // 权限的简称
	AppID       uint        `gorm:"UNIQUE_INDEX:idx_app_id_code" json:"app_id""`                // 应用的ID
	Application Application `gorm:"ForeignKey:AppID"`                                           // 使用AppID作为外键
}
