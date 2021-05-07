package controllers

import (
	"github.com/codelieche/microservice/usercenter/core"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseController struct{}

// HandleOK 响应200正确的结果
func (controller *BaseController) HandleOK(c *gin.Context, data interface{}) {
	r := core.Response{
		Code:    0,
		Data:    data,
		Message: "ok",
	}
	c.JSON(http.StatusOK, r)
}

// HandleCreated 响应201正确的结果
func (controller *BaseController) HandleCreated(c *gin.Context, data interface{}) {
	r := core.Response{
		Code:    0,
		Data:    data,
		Message: "ok",
	}
	c.JSON(http.StatusCreated, r)
}

// HandleError 响应400错误
func (controller *BaseController) HandleError(c *gin.Context, err error) {
	if err == core.ErrNotFound {
		controller.Handle404(c, err)
		return
	}

	r := core.Response{
		Code:    http.StatusBadRequest,
		Message: err.Error(),
	}

	c.JSON(http.StatusBadRequest, r)
	//c.AbortWithStatusJSON(400, err.Error())
}

// Handle404 响应404错误
func (controller *BaseController) Handle404(c *gin.Context, err error) {
	r := core.Response{
		Code:    http.StatusNotFound,
		Message: err.Error(),
	}
	c.JSON(http.StatusNotFound, r)
	//c.AbortWithStatusJSON(400, err.Error())
}

// HandleServerError 响应500错误
func (controller *BaseController) HandleServerError(c *gin.Context, err error) {
	r := core.Response{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	}
	c.JSON(http.StatusInternalServerError, r)
}
