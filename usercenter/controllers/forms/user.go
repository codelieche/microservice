package forms

import "fmt"

type UserCreateForm struct {
	Username   string `json:"username" form:"username"`
	Password   string `json:"password" form:"password"`
	RePassword string `json:"re_password" form:"re_password"`
	Nickname   string `json:"nickname" form:"nickname"`
	Email      string `json:"email" form:"email"`
	Phone      string `json:"phone" form:"phone"`
}

func (form *UserCreateForm) Validate() error {
	var err error
	// 1. 判断密码和确认密码是否相等
	if form.Password != form.RePassword {
		err = fmt.Errorf("密码和确认密码不相等")
	}

	// 2： 密码是否为空
	if form.Password == "" {
		err = fmt.Errorf("密码为空")
	}

	// 3. 如果未配置Nickname那么就让其等于username
	return err
}
