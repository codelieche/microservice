package controllers

import (
	"github.com/codelieche/microservice/codelieche/filters"
	"github.com/codelieche/microservice/usercenter/core"
	"github.com/codelieche/microservice/usercenter/internal/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

// Handle401 响应401错误
func (controller *BaseController) Handle401(c *gin.Context, err error) {
	r := core.Response{
		Code:    http.StatusUnauthorized,
		Message: err.Error(),
	}
	c.JSON(http.StatusUnauthorized, r)
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

// ParsePagination 解析分页
func (controller *BaseController) ParsePagination(c *gin.Context) *core.Pagination {
	// 分页我们一般是根据?page=1&page_size=10 类分割
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	// 如果配置了最大页数，超出返回我们就重置page
	if config.MaxPage > 0 && page > config.MaxPage {
		page = config.MaxPage
	}

	// 获取pageSize
	pageSizeStr := c.DefaultQuery(config.PageSizeQueryParam, "10")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 10
	}
	// 如果pageSize大于了设置的范围，那么我们限制一下
	if pageSize > config.MaxPageSize {
		pageSize = config.MaxPageSize
	}

	// 返回分页对象
	return &core.Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

// FilterAction 处理过滤操作
func (controller *BaseController) FilterAction(
	c *gin.Context, filterOptions []*filters.FilterOption,
	searchFields []string, orderingFields []string, defaultOrdering string) (filterActions []filters.Filter) {

	// 1-1: 过滤的操作
	filterAction := filters.FromQueryGetFilterAction(c, filterOptions)
	// 判断过滤操作是否为空
	if filterAction != nil {
		filterActions = append(filterActions, filterAction)
	}

	// 1-2: 搜索的操作
	searchAction := filters.FromQueryGetSearchAction(c, []string{"username", "email", "phone"})
	// 判断搜索操作是否为空
	if searchAction != nil {
		filterActions = append(filterActions, searchAction)
	}

	// 1-3: 排序
	var orderingAction filters.Filter
	if orderingFields != nil && defaultOrdering != "" {
		orderingAction = filters.FromQueryGetOrderingActionWithDefault(c, orderingFields, defaultOrdering)
	} else {
		orderingAction = filters.FromQueryGetOrderingAction(c, orderingFields)
	}
	// 判断排序操作是否为空
	if orderingAction != nil {
		filterActions = append(filterActions, orderingAction)
	}

	// 2. 最后返回过滤操作的列表
	return filterActions
}
