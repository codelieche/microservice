package forms

type PermissionCreateForm struct {
	Project string `json:"project" form:"project" validate:"required"`
	Name    string `json:"name" validate:"required"`
	Code    string `json:"code" validate:"required"`
}

type PermissionCheckForm struct {
	App   string   `json:"app" validate:"required"`
	Codes []string `json:"code" validate:"required"`
}
