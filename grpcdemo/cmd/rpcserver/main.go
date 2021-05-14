package main

import (
	"github.com/codelieche/microservice/codelieche/rpctools"
	"github.com/codelieche/microservice/grpcdemo/proto/pb"
	"github.com/codelieche/microservice/grpcdemo/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	// 1. 准备相关数据
	// 1-1：日志logger实例
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Can not create zap logger:%v", err)
	} else {
		defer logger.Sync()
	}

	// 1-2： 监听的grpc地址
	addr := "0.0.0.0:9081"

	// 1-3: grpc server option
	var opts []grpc.ServerOption

	// 2. 注册gRPC Service函数
	registerFunc := func(s *grpc.Server) {
		// 需要注入几个服务就注入几个
		// 1. hello greeter service
		helloService := server.NewHelloService(logger)
		pb.RegisterGreeterServer(s, helloService)

		// 2. stream service
	}

	// 3. 实例化gGRPC server config
	grpcCfg := &rpctools.GrpcServerConfig{
		Name:         "grpcdemo",
		Opts:         opts,
		Address:      addr,
		RegisterFunc: registerFunc,
		Logger:       logger,
	}

	// 启动grpc gateway server
	go startGRPCGateway()

	// 4. 启动grpc服务
	if err := rpctools.RunGRPCServer(grpcCfg); err != nil {
		log.Fatalf("grpc server run error:%v", err)
	}

}
