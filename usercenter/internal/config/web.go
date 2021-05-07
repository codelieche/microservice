package config

type Web struct {
	Address string
	Port    int
	JWT     *JWT // JWT相关的配置
}

type JWT struct {
	Key      string
	Duration int
	Issuer   string
}

// PageSizeQueryParam page size query参数值默认是page_size
const PageSizeQueryParam = "page_size"

// MaxPageSize 分页最大的数量
const MaxPageSize = 100

// MaxPage 最大的页数
const MaxPage = 0
