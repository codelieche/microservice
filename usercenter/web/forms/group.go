package forms

// 分组创建表单
type GroupCreateFrom struct {
	Name        string  `form:"name" validate:"required" json:"name"`
	Users       []int64 `form:"users" validate:"omitempty" json:"users"`
	Permissions []int64 `form:"permissions" validate:"omitempty"`
}
