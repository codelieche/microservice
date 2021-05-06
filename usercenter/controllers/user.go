package controllers

import (
	"github.com/codelieche/microservice/usercenter/controllers/forms"
	"github.com/codelieche/microservice/usercenter/core"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

type UserController struct {
	BaseController
	service core.UserService
}

func NewUserController(service core.UserService) *UserController {
	return &UserController{service: service}
}

// Create 通过Post创建用户
func (controller *UserController) Create(c *gin.Context) {
	// 1. 处理表单
	var form forms.UserCreateForm
	if err := c.ShouldBind(&form); err != nil {
		log.Println("shouldBind error:", err.Error())
		controller.HandleError(c, err)
		return // 记得返回
	}

	// 2. 对表单进行校验
	if err := form.Validate(); err != nil {
		controller.HandleError(c, err)
		return
	}

	// 3. 准备创建用户
	user := &core.User{
		Username: form.Username,
		Password: form.Password,
		Email:    form.Email,
		Phone:    form.Phone,
	}

	if user, err := controller.service.Create(c.Request.Context(), user); err != nil {
		controller.HandleError(c, err)
		return
	} else {
		controller.HandleCreated(c, user)
	}
}

// Find 通过GET获取用户信息
func (controller *UserController) Find(c *gin.Context) {
	// 1. 获取用户的id
	id := c.Param("id")
	if userId, err := strconv.Atoi(id); err != nil {
		controller.HandleError(c, err)
		return
	} else {
		if user, err := controller.service.Find(c.Request.Context(), int64(userId)); err != nil {
			controller.HandleError(c, err)
			return
		} else {
			controller.HandleOK(c, user)
		}
	}
}
