module github.com/codelieche/microservice/examples

go 1.13

require (
	github.com/codelieche/microservice/middleware v0.0.0-20200104112237-46c52fc104c0
	github.com/codelieche/microservice/usercenter v0.0.0-20200104112237-46c52fc104c0
	github.com/kataras/iris/v12 v12.1.4
)

replace (
	github.com/codelieche/microservice/middleware => ../middleware
	github.com/codelieche/microservice/usercenter => ../usercenter
)
