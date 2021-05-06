package config

type Redis struct {
	Host     string // 主机地址
	Port     int    // 端口号
	Db       int    // DB index
	Password string // 密码
}
