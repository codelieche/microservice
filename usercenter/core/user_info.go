package core

type UserInfo struct {
	Id          int64  `json:"id" form:"id"`
	Nickname    string `json:"nickname" form:"nickname"`
	Username    string `json:"username" form:"username"`
	Email       string `json:"email" form:"email"`
	Phone       string `json:"phone" form:"phone"`
	IsSuperuser bool   `json:"is_superuser" form:"is_superuser"` // 是否超级用户
	IsActive    bool   `json:"is_active" form:"is_active"`       // 是否启用
	Avatar      string `json:"avatar" form:"avatar"`             // 头像地址
	Signature   string `json:"signature" form:"signature"`       // 个性签名
}
