package main

import (
	"github.com/codelieche/microservice/usercenter/apiserver/app"
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {

	// 实例化 Web App
	application := app.NewApp()

	application.Run()
}
