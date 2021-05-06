package utils

/**
ConfigParse:
默认的环境变量：
1. CONFIGNAME: 配置文件的名称，默认是：config
2. CONFIG_PATH: 配置文件所在路径/目录,默认是: ./
3. 配置文件类型，我们直接使用yaml
4. ENV_CONNECT_SYMBOL: 环境变量连接符，默认是(_)，比如：MYSQL_PASSWORD会解析成mysql.password

配置文件处理流程：
1. 判断CONFIG_PATH传递的是目录还是文件
2. 如果路径是目录，就从目录中读取config.yaml文件
3. 如果是配置文件路径，且需要是.yaml结尾的文件，那么读取文件加载到配置中
4. 从环境变量中加载配置，比如：MYSQL_PASSWORD就是mysql.password

配置文件顺序：
1. 文件中的配置，比如：mysql.password
2. 环境变量的配置，比如：MYSQL_PASSWORD

ConfigSetDefault: 是调用viper.SetDefault(string, interface{})方法
*/

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

func ConfigParse() error {
	log.Println("开始读取配置内容")
	// 1. 从配置文件读取
	// 1.1 读取环境变量，默认是./
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./"
	}

	// 1.2 判断是目录还是文件
	if fileInfo, err := os.Stat(configPath); err != nil {
		log.Println("统计文件路径出错：", err)
		return err
	} else {
		log.Println("配置文件路径为：", configPath)
		if os.IsNotExist(err) {
			log.Println("文件/目录不存在:", err)
			return err
		}

		// 无论是目录还是文件，我们配置文件类型都为yaml
		viper.SetConfigType("yaml")

		// 1.3 是目录的情况
		if fileInfo.IsDir() {
			configName := os.Getenv("CONFIGNAME")
			if configName == "" {
				configName = "config"
			}
			// 传递的是目录
			viper.AddConfigPath(configPath)
			viper.SetConfigName(configName)
			// 读取配置的路径的config.yaml配置文件
			if err := viper.ReadInConfig(); err != nil {
				log.Println("error", err.Error())
				return err
			}

		} else {
			// 1.4 传递的配置是个文件
			if file, err := os.Open(configPath); err != nil {
				log.Println("读取配置文件出错：", err)
				return err
			} else {
				if err := viper.ReadConfig(file); err != nil {
					log.Println("读取配置出错：", err)
					return err
				} else {
					log.Println("读取配置文件成功")
				}
			}
		}
	}

	// 2. 绑定环境变量
	// 环境变量的分隔符: viper默认的是. 我们使用下划线来连接
	// 比如会把MYSQL_PASSWORD解析到mysql.password
	connectSymbol := os.Getenv("ENV_CONNECT_SYMBOL")
	if connectSymbol == "" {
		connectSymbol = "_"
	}
	viper.AutomaticEnv()

	replacer := strings.NewReplacer(".", connectSymbol)
	viper.SetEnvKeyReplacer(replacer)

	// 3. 返回
	return nil
}

func ConfigSetDefault(in map[string]interface{}) {
	// 配置设置默认值
	for key := range in {
		viper.SetDefault(key, in[key])
	}
}
