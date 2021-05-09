package config

import (
	"github.com/codelieche/microservice/codelieche/utils"
	"github.com/spf13/viper"
	"log"
)

type GlobalConfig struct {
	Web   *Web   // web api相关的配置
	Rpc   *Rpc   // rpc server相关的配置
	MySQL *MySQL // 数据库的配置
	Redis *Redis // Redis相关的配
}

var Config *GlobalConfig

var JwtTokenHeaderPrefix = "Bearer"

var JwtAuthBlackUrlPathMap = map[string]bool{
	"/api/v1/user/login/":  true,
	"/api/v1/user/create/": true,
}

func parseConfig() {
	// 配置默认的值
	defaultConfigValue := map[string]interface{}{
		"web.address":      "0.0.0.0",
		"web.port":         8080,
		"web.jwt.duration": 3600 * 12,
		"web.jwt.key":      "codelieche",
		"web.jwt.issuer":   "codelieche",

		"rpc.address": "0.0.0.0",
		"rpc.port":    "8081",

		"mysql.host":     "127.0.0.1",
		"mysql.port":     3306,
		"mysql.database": "usercenter",
		"charset":        "utf8mb4",

		"redis.host": "127.0.0.1",
		"redis.port": 6379,
		"redis.db":   1,
	}

	// 配置默认的值
	utils.ConfigSetDefault(defaultConfigValue)

	// 调用处理配置的逻辑
	if err := utils.ConfigParse(); err != nil {
		log.Printf("解析配置文件出错：%s", err.Error())
		return
	}

	// 处理配置
	jwt := &JWT{
		Key:      viper.GetString("web.jwt.key"),
		Duration: viper.GetInt("web.jwt.duration"),
		Issuer:   viper.GetString("web.jwt.issuer"),
	}
	if jwt.Duration <= 0 {
		jwt.Duration = 3600 * 12
	}
	if jwt.Key == "" {
		jwt.Key = "codelieche"
	}
	if jwt.Issuer == "" {
		jwt.Issuer = "codelieche"
	}

	web := &Web{
		Address: viper.GetString("web.address"),
		Port:    viper.GetInt("web.port"),
		JWT:     jwt,
	}

	if web.Port <= 0 {
		web.Port = 8080
	}

	// 处理rpc的配置
	rpc := &Rpc{
		Address: viper.GetString("rpc.address"),
		Port:    viper.GetInt("rpc.port"),
	}
	if rpc.Port <= 0 {
		rpc.Port = 8081
	}

	mysql := &MySQL{
		Host:     viper.GetString("mysql.host"),
		Port:     viper.GetInt("mysql.port"),
		Database: viper.GetString("mysql.database"),
		Charset:  viper.GetString("mysql.charset"),
		User:     viper.GetString("mysql.user"),
		Password: viper.GetString("mysql.password"),
	}
	if mysql.Port <= 0 {
		mysql.Port = 3306
	}
	if mysql.Charset == "" {
		mysql.Charset = "utf8mb4"
	}

	redis := &Redis{
		Host:     viper.GetString("redis.host"),
		Port:     viper.GetInt("redis.port"),
		Db:       viper.GetInt("redis.db"),
		Password: viper.GetString("mysql.password"),
	}

	if redis.Port <= 0 {
		redis.Port = 6379
	}

	Config = &GlobalConfig{
		Web:   web,
		Rpc:   rpc,
		MySQL: mysql,
		Redis: redis,
	}
}

func GetConfig() *GlobalConfig {
	if Config != nil {
		return Config
	} else {
		parseConfig()
		if Config != nil {
			return Config
		} else {
			return nil
		}
	}
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	parseConfig()
}
