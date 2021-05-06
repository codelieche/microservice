package utils

import (
	"github.com/spf13/viper"
	"log"
	"testing"
)

func TestConfigSetDefault(t *testing.T) {
	var defaultValue = map[string]interface{}{
		"mysql.host": "127.0.0.1",
		"mysql.port": 3306,
		"mysql.user": "root",
	}
	ConfigSetDefault(defaultValue)
	// 获取值
	for key := range defaultValue {
		log.Println(key, viper.GetString(key))
	}

	if viper.GetString("mysql.user") == "root" {
		log.Println("设置默认值测试成功")
	} else {
		t.Error("设置默认值失败")
	}
}
