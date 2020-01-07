package forms

// 注册用户表单
type UserCreateForm struct {
	Username   string `json:"username" validate:"required,min=6,max=40"`
	Password   string `json:"password" validate:"required,min=8,max=40"`       // 必填，长度为：8 <= length <= 40
	Repassword string `json:"repassword" validate:"required,eqfield=Password"` // 必填，且需要和Password字段相同
	Mobile     string `json:"mobile" validate:"omitempty,max=100"`
	Email      string `json:"email" validate:"omitempty,email"`
}

// 编辑用户表单
type UserUpdateform struct {
	Mobile      string  `json:"mobile" validate:"omitempty,max=100"`
	Email       string  `json:"email" validate:"omitempty,email"`
	Groups      []int64 `json:"groups" validate:"int"`
	Roles       []int64 `json:"roles" validate:"int"`
	Permissions []int64 `json:"permissions" validate:"int"`
}

// 用户登录表单
type UserLoginForm struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"` // 必填
	Mobile   string `json:"mobile" validate:"omitempty,max=100"`
	Email    string `json:"email" validate:"omitempty,email"`
}

// 用户修改密码
type UserChangePasswrodForm struct {
	Password   string `json:"password" validate:"required,min=8,max=40"`       // 必填，长度为：8 <= length <= 40
	Repassword string `json:"repassword" validate:"required,eqfield=Password"` // 必填，且需要和Password字段相同
}

// 用户基本信息表单
type UserDetailForm struct {
	ID       int    `json:"id"`       // 用户ID
	Username string `json:"username"` // 用户名称
	Mobile   string `json:"mobile"`   // 用户手机号
	Email    string `json:"email"`    // 用户邮箱
}
