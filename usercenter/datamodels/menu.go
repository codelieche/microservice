package datamodels

/**
项目的菜单
前端会调用这个接口，渲染左边的导航
每一级菜单，都可设置相应的权限，有相关权限的人，才会显示当前菜单
*/
type Menu struct {
	BaseFields
	Project      string      `gorm:"type:varchar(40);NOT NULL;UNIQUE_INDEX:idx_project_slug" json:"project"` // 项目Code
	Title        string      `gorm:"type:varchar(20);NOT NULL;" json:"title"`                                // 菜单的标题
	Slug         string      `gorm:"type:varchar(100);NOT NULL;UNIQUE_INDEX:idx_project_slug" json:"slug"`   // 菜单网址
	Icon         string      `gorm:"type:varchar(40);default:'angle-right'" json:"icon"`                     // 菜单图标
	ParentID     uint        `gorm:"INDEX" json:"parent_id"`                                                 // 父菜单的ID
	Level        uint        `gorm:"default:'1'"`                                                            // 菜单级别
	PermissionID uint        `json:"permission_id"`                                                          // 权限ID
	Permission   *Permission `gorm:"foreignkey:PermissionID" json:"permission"`                              // 权限
	Target       string      `gorm:"type:varchar(10);default:'_self'" json:"target"`                         // 连接点击跳转方式
	IsLink       bool        `json:"is_link"`                                                                // 是否是站外链接
	Link         string      `gorm:"type:varchar(128)" json:"link"`                                          // 连接                                                                // 是否是站外链接
	Order        uint        `gorm:"default:'1'" json:"order"`                                               // 排序
	Children     []*Menu     `gorm:"foreignkey:ParentID" json:"children"`
	//Parent       *Menu       `gorm:"foreignkey:ParentID" json:"parent"`                                      // 父菜单
}
