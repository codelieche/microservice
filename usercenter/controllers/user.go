package controllers

import (
	"errors"
	"fmt"
	"github.com/codelieche/microservice/usercenter/controllers/forms"
	"github.com/codelieche/microservice/usercenter/core"
	"github.com/codelieche/microservice/usercenter/internal/config"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
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

	if user, err := controller.service.FindByIdOrUsername(c.Request.Context(), id); err != nil {
		controller.HandleError(c, err)
		return
	} else {
		controller.HandleOK(c, user)
		return
	}

	//if isDigit, err := regexp.Match("^\\d+$", []byte(id)); err != nil {
	//	controller.HandleError(c, err)
	//	return
	//} else {
	//	if isDigit {
	//		// 是数字类型
	//		if userId, err := strconv.Atoi(id); err != nil {
	//			controller.HandleError(c, err)
	//			return
	//		} else {
	//			if user, err := controller.service.Find(c.Request.Context(), int64(userId)); err != nil {
	//				controller.HandleError(c, err)
	//				return
	//			} else {
	//				controller.HandleOK(c, user)
	//			}
	//		}
	//	} else {
	//		// 字符类型
	//		if user, err := controller.service.FindByUsername(c.Request.Context(), id); err != nil {
	//			controller.HandleError(c, err)
	//			return
	//		} else {
	//			controller.HandleOK(c, user)
	//		}
	//	}
	//}
}

func (controller *UserController) Login(c *gin.Context) {
	// 1. 处理表单
	var form forms.UserLoginForm

	if err := c.ShouldBind(&form); err != nil {
		controller.HandleError(c, err)
		return
	}
	// 2. 准备登录
	// 2-1: 检查登录方式
	if form.Category != "" && form.Category != "username" {
		err := fmt.Errorf("暂时我们只支持通过用户名登录")
		controller.HandleError(c, err)
		return
	}
	//	2-2：获取用户
	user, err := controller.service.FindByUsername(c.Request.Context(), form.Username)
	if err != nil {
		if err == core.ErrNotFound {
			err = fmt.Errorf("用户不存在")
		}
		controller.HandleError(c, err)
		return
	}
	// 2-3：操作登录
	if ok, err := user.CheckPassword(form.Password); err != nil || !ok {
		err = fmt.Errorf("用户名或者密码错误")
		controller.HandleError(c, err)
		return
	}

	//	3. 返回JWT Token
	if token, err := controller.service.SigningToken(c.Request.Context(), user); err != nil {
		controller.HandleError(c, err)
		return
	} else {
		// 组合Data
		data := map[string]interface{}{
			"token": token,
			"user":  user,
		}
		controller.HandleOK(c, data)
		return
	}

}

func (controller *UserController) Auth(c *gin.Context) {
	// 1. 获取用户传递的Token
	authorizationHeader := c.GetHeader("Authorization")
	if authorizationHeader == "" {
		err := errors.New("请传入Token")
		controller.HandleError(c, err)
		return
	}
	// 获取tokenStr
	var tokenStr string
	if config.JwtTokenHeaderPrefix != "" {
		tokenStr = strings.TrimPrefix(authorizationHeader, fmt.Sprintf("%s ", config.JwtTokenHeaderPrefix))
	} else {
		tokenStr = authorizationHeader
	}

	// 解析token
	if claims, err := controller.service.ParseToken(c.Request.Context(), tokenStr); err != nil {
		controller.HandleError(c, err)
		return
	} else {
		controller.HandleOK(c, claims)
		return
	}
}
