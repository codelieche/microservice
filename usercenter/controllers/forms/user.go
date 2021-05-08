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
		err = fmt.Errorf("密码和确认密码不相同")
		return err
	}

	// 2： 密码是否为空
	if form.Password == "" {
		err = fmt.Errorf("密码为空")
		return err
	}

	// 3. 如果未配置Nickname那么就让其等于username
	return err
}

type UserLoginForm struct {
	Username string `json:"username" form:"username"` // 用户名
	Password string `json:"password" form:"password"` // 密码
	Category string `json:"category" form:"category"` // 登录方式
}

type UserChangePasswordForm struct {
	Username    string `json:"username" form:"username"`
	OldPassword string `json:"old_password" form:"old_password"`
	Password    string `json:"password" form:"password"`
	RePassword  string `json:"re_password" form:"re_password"`
}

func (form *UserChangePasswordForm) Validate(isReset bool) error {
	var err error
	// 1. 判断密码和确认密码是否相等
	if form.Password != form.RePassword {
		err = fmt.Errorf("密码和确认密码不相同")
		return err
	}

	if !isReset {
		if form.Password == form.OldPassword {
			err = fmt.Errorf("新密码和老密码相同")
			return err
		}
	}

	// 2： 密码是否为空
	if form.Password == "" {
		err = fmt.Errorf("密码为空")
		return err
	}

	// 3. 如果未配置Nickname那么就让其等于username
	return err
}
