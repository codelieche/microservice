package forms

type RoleCreateForm struct {
	Name        string  `json:"name" validate:"required;min=4"`
	Users       []int64 `json:"users" validate:"omitempty"`
	Permissions []int64 `json:"permissions" validate:"omitempty"`
}
