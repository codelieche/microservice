package config

import "fmt"

type MySQL struct {
	Host     string // 数据库地址
	Port     int    // 数据库端口
	Database string // 数据库
	Charset  string // 字符类型
	User     string // 用户名
	Password string // 用户密码
}

func (m *MySQL) GetDSN() string {
	if m.Database == "" {
		panic("需要配置数据库")
		return ""
	}
	if m.Charset == "" {
		m.Charset = "utf8mb4"
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true", m.User, m.Password, m.Host, m.Port, m.Database, m.Charset)
	return dsn
}
