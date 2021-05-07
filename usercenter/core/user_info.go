package core

type UserInfo struct {
	Id          int64  `json:"id" form:"id"`                     // 用户ID
	Nickname    string `json:"nickname" form:"nickname"`         // 昵称
	Username    string `json:"username" form:"username"`         // 用户名
	Email       string `json:"email" form:"email"`               // 邮箱
	Phone       string `json:"phone" form:"phone"`               // 电话号码
	IsSuperuser bool   `json:"is_superuser" form:"is_superuser"` // 是否超级用户
	IsActive    bool   `json:"is_active" form:"is_active"`       // 是否启用
	Avatar      string `json:"avatar" form:"avatar"`             // 头像地址
	Signature   string `json:"signature" form:"signature"`       // 个性签名
}
