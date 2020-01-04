package main

import (
	"log"

	"github.com/codelieche/microservice/usercenter/app"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("程序开始执行")
	app.Run()
	log.Println("程序执行完毕")
}
