package common

import (
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

var config *Config

type AppInfo struct {
	ID    int    `json:"id" yaml:"id"`
	Code  string `json:"code" yaml:"code"`
	Token string `json:"token" yaml:"token"`
}

// DingDing 项目配置
type Config struct {
	Http     *HttpConfig
	Database *Database
	App      *AppInfo
	Debug    bool
}

// Http Config
type HttpConfig struct {
	Host      string            `json:"host" yaml:"host"`
	Port      int               `json:"port" yaml:"port"`
	Timeout   int               `json:"timeout" yaml:"timeout"`
	BasicAuth map[string]string `json:"basicauth" yaml:"basicauth"`
	Domains   []string          `json:"domains" yaml:"domains"` // 可以跳转的域名
}

type Database struct {
	MySQL *MySQLDatabase `json:"mysql" yaml:"mysql"`
	Redis *RedisDatabase `json:"redis" yaml:"redis"`
}

// MySQL数据库相关配置
type MySQLDatabase struct {
	Host     string `json:"host" yaml:"host"`          // 数据库地址
	Port     int    `json:"port" yaml:"port"`          // 端口号
	User     string `json:"user" yaml:"user"`          // 用户
	Password string `json:"password" yaml:"password"`  // 用户密码
	Database string `json:"database" yaml: "database"` // 数据库
}

// Redis配置
type RedisDatabase struct {
	Host     string   `json:"host" yaml:"host"`         // redis主机，不填会是默认的127.0.0.1：6739
	Clusters []string `json:"clusters" yaml:"clusters"` // Redis集群地址
	Password string   `json:"password" yaml:"password"` // redis的密码
	DB       int      `json:"db" yaml:db`               // 哪个库
}

// 解析项目的相关配置
func ParseConfig() (err error) {
	var (
		fileName   string
		content    []byte
		contentStr string
	)
	// log.Println(os.Getenv("PWD"))

	// 获取配置文件: 每次要调试，执行的时候工作路径不同，所以设置成用环境变量来处理
	// 如果传递的最后一个参数是.yaml那么它是配置文件
	if strings.HasSuffix(os.Args[len(os.Args)-1], ".yaml") {
		fileName = os.Args[len(os.Args)-1]
	} else {
		if os.Getenv("CONFIG_FILENAME") != "" {
			fileName = os.Getenv("CONFIG_FILENAME")
		} else {
			fileName = "./config.yaml"
		}
	}

	// 判断文件是否存在
	if _, err = os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			if fileName == "./config.yaml" {
				fileName = "../config.yaml"
			} else {
				log.Println("配置文件不存在：", fileName)
				return
			}
		}
	}

	if content, err = ioutil.ReadFile(fileName); err != nil {
		return
	} else {
		contentStr = string(content)
		//log.Println(contentStr)

		// 正则替换环境变量
		r := regexp.MustCompile(`\$\{(.*?)\}`)
		results := r.FindAllStringSubmatch(contentStr, -1)

		for _, envStr := range results {
			var envName, envValue, envDefault string
			if envStr[1] != "" {
				envNameAndDefaultArry := strings.Split(envStr[1], ":")
				envName = envNameAndDefaultArry[0]
				envValue = os.Getenv(envName)
				if len(envNameAndDefaultArry) == 2 {
					envDefault = envNameAndDefaultArry[1]

				}
				if envValue == "" && envDefault != "" {
					envValue = envDefault
				}
			}
			// 对环境变量进行替换
			contentStr = strings.ReplaceAll(contentStr, envStr[0], envValue)
		}

		// 替换完了置换，修改content
		content = []byte(contentStr)

		//log.Println(string(content))
	}

	// 解析配置
	config = &Config{
		Http: &HttpConfig{
			Host:    "0.0.0.0",
			Port:    3306,
			Timeout: 30,
		},
		Database: &Database{
			MySQL: &MySQLDatabase{
				Host:     "127.0.0.1",
				Port:     3306,
				User:     "root",
				Password: "",
				Database: "usercenter",
			},
			Redis: &RedisDatabase{
				Host:     "127.0.0.1:6379",
				Clusters: []string{"127.0.0.1:6379"},
				Password: "",
				DB:       0,
			},
		},
		App: &AppInfo{
			ID:    1,
			Code:  "default",
			Token: "default",
		},
		Debug: false,
	}

	if err = yaml.Unmarshal(content, config); err != nil {
		log.Println("解析配置文件出错：", err.Error())
		return err
	} else {
		// 解析配置成功
		//log.Println(*config.Http)
		//log.Println(*config.DingDing)
		//log.Println(*config.Database)

		// 对app进行校验

		if config.Http.Timeout < 10 {
			config.Http.Timeout = 30
		}

	}
	return
}

func GetConfig() *Config {
	if config != nil {
		return config
	} else {
		if err := ParseConfig(); err != nil {
			log.Println(err.Error())
			return nil
		} else {
			return config
		}
	}
}
