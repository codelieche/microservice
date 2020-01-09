package forms

type PermissionCreateForm struct {
	AppID int64  `json:"app_id" form:"app_id" validate:"required"`
	Name  string `json:"name" validate:"required"`
	Code  string `json:"code" validate:"required"`
}
