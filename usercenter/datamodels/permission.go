package datamodels

type Permission struct {
	//gorm.Model
	BaseFields
	Name    string `gorm:"type:varchar(100)" json:"name"`                                 // 权限的名称
	Code    string `gorm:"size:100;UNIQUE_INDEX:idx_project_code;NOT NULL" json:"code"`   // 权限的简称
	Project string `gorm:"type:varchar(40);UNIQUE_INDEX:idx_project_code" json:"project"` // 项目的Code
	//AppID       uint         `gorm:"UNIQUE_INDEX:idx_app_id_code" json:"app_id""`       // 应用的ID
	//Application *Application `gorm:"ForeignKey:AppID" json:"application"`               // 使用AppID作为外键
}
