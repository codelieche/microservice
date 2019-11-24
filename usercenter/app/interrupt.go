package app

import (
	"log"

	"github.com/codelieche/microservice/datasources"
)

// 处理control/cmd + c关闭的时候
func handleAppInterupt() {
	log.Println("程序即将退出！")
	// 关闭数据库连接等
	db := datasources.GetDb()
	db.Close()
}
