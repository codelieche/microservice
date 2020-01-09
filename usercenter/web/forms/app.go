package forms

type ApplicationCreateForm struct {
	Name string `json:"name" validate:"required;"`
	Code string `json:"code" validate:"min=5"`
	Info string `json:"info"`
}

type ApplicationUpdateForm struct {
	Name string `json:"name" validate:"required;"`
	Code string `json:"code" validate:"min=5"`
	Info string `json:"info"`
	//Permissions []int64 `json:"permissions" validate:"omitempty"`
}
