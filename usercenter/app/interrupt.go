package app

import (
	"log"

	"github.com/codelieche/microservice/usercenter/datasources"
)

// 处理control/cmd + c关闭的时候
func handleAppInterupt() {
	log.Println("程序即将退出！")
	// 关闭数据库连接等
	db := datasources.GetDb()
	db.Close()

	// 关闭session的redis数据库
	redisDB.Close()

	// 断开redis的连接
	redisClient := datasources.GetRedisClient()
	redisClient.Close()
}
