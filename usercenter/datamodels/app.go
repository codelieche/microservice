package datamodels

import (
	"encoding/json"
)

type Application struct {
	//gorm.Model
	BaseFields
	Name        string          `gorm:"size:100;NOT NULL;INDEX" json:"name"`        // 应用的名称
	Code        string          `gorm:"size:100;NOT NULL;UNIQUE_INDEX" json:"code"` // 应用的code 唯一值
	Token       string          `gorm:"size:100" json:"token"`                      // 每个应用给个TOKEN
	Info        json.RawMessage `gorm:"type:text" json:"info,omitempty"`            // 应用的信息
	Permissions []*Permission   `gorm:"ForeignKey:AppID"`                           // 应用的权限
}

//func (app *Application) TableName() string {
//	return "usercenter_applications"
//}
