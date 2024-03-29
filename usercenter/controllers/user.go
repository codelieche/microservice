package controllers

import (
	"errors"
	"fmt"
	"github.com/codelieche/microservice/codelieche/filters"
	"github.com/codelieche/microservice/usercenter/controllers/forms"
	"github.com/codelieche/microservice/usercenter/core"
	"github.com/codelieche/microservice/usercenter/internal/config"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
	"sync"
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
		//controller.HandleOK(c, claims)
		// 获取用户
		if user, err := controller.service.FindByUsername(c.Request.Context(), claims.Username); err != nil {
			controller.HandleError(c, err)
			return
		} else {
			if user.IsActive {
				controller.HandleOK(c, user)
			} else {
				err = errors.New("用户已被禁用")
				controller.HandleError(c, err)
			}
		}
	}
}

func (controller *UserController) List(c *gin.Context) {
	// 1. 获取相关数据
	// 1-1. 获取分页
	pagination := controller.ParsePagination(c)

	// 1-2：过滤操作选项：过滤和联合表的字段
	filterOptions := []*filters.FilterOption{
		&filters.FilterOption{
			QueryKey: "id__gte",
			Column:   "id",
			Op:       filters.FILTER_GTE,
		},
		&filters.FilterOption{
			QueryKey: "id__lte",
			Column:   "id",
			Op:       filters.FILTER_LTE,
		},
		&filters.FilterOption{
			QueryKey: "username__contains",
			Column:   "username",
			Op:       filters.FILTER_CONTAINS,
		},
		&filters.FilterOption{
			QueryKey: "nickname__contains",
			Column:   "nickname",
			Op:       filters.FILTER_CONTAINS,
		},
		&filters.FilterOption{
			Column: "phone",
		},
	}

	// 1-3: 搜索字段
	searchFields := []string{"username", "email", "phone"}

	// 1-4: 排序字段，默认排序为id
	orderingFields := []string{"id", "username", "phone", "email"}
	defaultOrdering := "id"

	// 1-5：获取过滤操作的列表【来自继承的BaseController】
	filterActions := controller.FilterAction(c, filterOptions, searchFields, orderingFields, defaultOrdering)

	// 2. 开始获取数据
	offset := pagination.PageSize * (pagination.Page - 1)
	ctx := c.Request.Context()

	wg := sync.WaitGroup{}
	wg.Add(2)

	var users []*core.User
	var err error
	var count int64

	// 协程1：获取用户列表
	go func() {
		defer wg.Done()
		users, err = controller.service.List(ctx, offset, pagination.PageSize, filterActions...)
	}()
	// 协程2：获取用户数量
	go func() {
		// 获取用户数
		defer wg.Done()
		count, err = controller.service.Count(ctx, filterActions...)
	}()
	// 等待2个协程完成
	wg.Wait()

	// 3. 处理结果
	if err != nil {
		controller.HandleError(c, err)
		return
	} else {
		r := core.ResponseList{
			CurrentPage: pagination.Page,
			Count:       count,
			Results:     users,
		}
		controller.HandleOK(c, r)
		return
	}
}

func (controller *UserController) ChangePassword(c *gin.Context) {
	// 1. 获取表单
	var form forms.UserChangePasswordForm
	if err := c.ShouldBind(&form); err != nil {
		controller.HandleError(c, err)
		return
	}

	// 2. 开始校验密码
	// 2-1: 校验用户名
	var username string
	if u, exist := c.Get("username"); exist {
		if v, ok := u.(string); ok {
			username = v
		}
	}
	if username == "" {
		err := errors.New("请传入用户名")
		controller.HandleError(c, err)
		return
	} else {
		// 只可修改自己的密码
		if username != form.Username {
			err := errors.New("只可修改自己的密码")
			controller.HandleError(c, err)
			return
		}
	}

	// 2-2: 校验数据，老的密码
	if err := form.Validate(false); err != nil {
		controller.HandleError(c, err)
		return
	}

	if user, err := controller.service.FindByUsername(c.Request.Context(), form.Username); err != nil {
		controller.HandleError(c, err)
		return
	} else {
		if ok, err := user.CheckPassword(form.OldPassword); err != nil || !ok {
			err := errors.New("原来的密码不正确")
			controller.HandleError(c, err)
			return
		}

		//	开始修改密码
		if err := controller.service.SetPassword(c.Request.Context(), user, form.Password); err != nil {
			controller.HandleError(c, err)
			return
		} else {
			// 密码修改成功
			controller.HandleOK(c, "密码修改成功")
			return
		}

	}
}

func (controller *UserController) ResetPassword(c *gin.Context) {
	// 1. 获取表单
	var form forms.UserChangePasswordForm
	if err := c.ShouldBind(&form); err != nil {
		controller.HandleError(c, err)
		return
	}

	// 2. 开始校验密码
	// 2-1: 校验用户名
	username := form.Username
	if username == "" {
		err := errors.New("请传入用户名")
		controller.HandleError(c, err)
		return
	} else {
		// 只有管理员才可设置密码
		// 判断一下当前用户是否是管理员
	}

	// 2-2: 校验数据，老的密码
	if err := form.Validate(true); err != nil {
		controller.HandleError(c, err)
		return
	}

	if user, err := controller.service.FindByUsername(c.Request.Context(), form.Username); err != nil {
		controller.HandleError(c, err)
		return
	} else {
		//	开始修改密码
		if err := controller.service.SetPassword(c.Request.Context(), user, form.Password); err != nil {
			controller.HandleError(c, err)
			return
		} else {
			// 密码修改成功
			msg := fmt.Sprintf("重置%s修改成功", username)
			controller.HandleOK(c, msg)
			return
		}
	}
}
