package controllers

import (
	"log"
	"time"

	"github.com/go-playground/validator"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

type DemoController struct {
	Session *sessions.Session
	Ctx     iris.Context
}

func (c *DemoController) Get() mvc.Result {

	return mvc.Response{
		Code: iris.StatusOK,
		Object: iris.Map{
			"time": time.Now(),
		},
	}
}

type User01 struct {
	Age        uint   `json:"age" validate:"gte=0,lte=120"`
	Username   string `json:"username" validate:"required"`
	Password   string `json:"password" validate:"required"`
	Repassword string `json:"repassword,omitempty" validate:"omitempty,eqfield=Password"`
	Email      string `json:"email,omitempty" validate:"omitempty,email"`
	Category   string `json:"category" validate:"omitempty,oneof= 1 2 3"`
}

func (c *DemoController) PostValidate() mvc.Result {
	var v *validator.Validate

	v = validator.New()
	//v.RegisterStructValidation(UserStructLevelValidation, &User{})

	var user User01
	if err := c.Ctx.ReadJSON(&user); err != nil {
		log.Println(err.Error())
		return mvc.Response{
			Err: err,
		}
	}

	if err := v.Struct(user); err != nil {
		//if _, ok := err.(*validator.InvalidValidationError); ok {
		//	log.Println(err.Error())
		//}
		for _, e := range err.(validator.ValidationErrors) {
			log.Println(e.Namespace(), e.Field(), e.StructNamespace())
			log.Println("tag\t", e.Tag())
			log.Println(e.Value(), e.Param())
			log.Println(e)
			log.Println("\n")
		}
		return mvc.Response{
			Err: err,
		}
	} else {
		return mvc.Response{
			Object: user,
		}
	}

}

func UserStructLevelValidation(sl validator.StructLevel) {
	user := sl.Current().Interface().(User01)

	if len(user.Username) == 0 && len(user.Password) == 0 {
		sl.ReportError(user.Username, "Username", "username", "User Name", "")
		sl.ReportError(user.Password, "Password", "password", "Password", "")
	}

}
