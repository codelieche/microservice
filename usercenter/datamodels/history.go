package datamodels

// 用户安全日志
// 类型：default、safe、login、other等
type SafeLog struct {
	BaseFields
	Category uint   `gorm:"INDEX" json:"category"`            // 分类
	UserID   uint   `gorm:"INDEX" json:"user_id""`            // 用户的ID
	User     *User  `gorm:"ForeignKey:UserID" json:"user"`    // 使用UserID作为外键
	Content  string `gorm:"type:varchar(256)" json:"content"` // 日志内容
	Success  bool   `gorm:"boolean" json:"success"`           // 是否登录成功
	Address  string `gorm:"size:20" json:"address"`           // 登录地址
	Device   string `gorm:"size:256" json:"device"`           // 登录的客户端：设备
}
