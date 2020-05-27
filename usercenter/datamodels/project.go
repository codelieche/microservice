package datamodels

/**
项目
Code是唯一的，从项目中心同步数据过来
*/
type Project struct {
	BaseFields
	Name        string         `gorm:"size:100;NOT NULL;INDEX" json:"name"`            // 项目名称
	Code        string         `gorm:"size:100;NOT NULL;UNIQUE_INDEX" json:"code"`     // 项目Code
	Description string         `gorm:"type:text" json:"description,omitempty"`         // 项目描述
	Users       []*ProjectUser `gorm:"foreignkey:Project;association_foreignkey:Code"` // 项目成员
	Permissions []*Permission  `gorm:"foreignkey:Project;association_foreignkey:Code"` // 项目权限
}

// 项目成员
type ProjectUser struct {
	BaseFields
	Role    string `gorm:"type:varchar(40);" json:"role"`                                 // 项目成员角色
	Project string `gorm:"type:varchar(40);UNIQUE_INDEX:idx_project_user" json:"project"` // 项目Code
	User    string `gorm:"type:varchar(40);UNIQUE_INDEX:idx_project_user" json:"user"`    // 用户的username
}
