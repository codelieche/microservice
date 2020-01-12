package datasources

import (
	"log"
	"os"

	"github.com/codelieche/microservice/usercenter/common"

	"github.com/go-redis/redis/v7"
)

var redisClient *redis.Client

func connectRedis() {
	config := common.GetConfig()
	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.Database.Redis.Host,
		Password: config.Database.Redis.Password, // no password set
		DB:       config.Database.Redis.DB,       // use default DB
	})

	if pong, err := redisClient.Ping().Result(); err != nil {
		log.Println(err)
		os.Exit(1)
	} else {
		// PONG
		log.Println(pong)
	}
}

func GetRedisClient() *redis.Client {
	if redisClient != nil {
		//if _, err := redisClient.Ping().Result(); err != nil {
		//	log.Println("连接redis出错：", err)
		//}
		return redisClient
	} else {
		connectRedis()
		return redisClient
	}
}
