package forms

type ProjectCreateForm struct {
	Name        string `json:"name" validate:"required"`
	Code        string `json:"code" validate:"min=5"`
	Description string `json:"description" validate:"omitempty"`
}

// 添加项目成员
type ProjectUserAddForm struct {
	Project  string `json:"project" validate:"required"`
	Username string `json:"username" validate:"required"`
	Role     string `json:"role""`
}
