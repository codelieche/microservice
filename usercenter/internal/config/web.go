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
