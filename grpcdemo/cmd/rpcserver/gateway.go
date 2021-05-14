package main

import (
	"context"
	"github.com/codelieche/microservice/codelieche/rpctools"
	"github.com/codelieche/microservice/grpcdemo/proto/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func startGRPCGateway() {
	// 1. 准备相关参数
	// 1-1. 准备mux
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(
		runtime.MIMEWildcard, &runtime.JSONPb{},
	))

	//	1-2. 准备地址
	address := "0.0.0.0:9082"

	// 1-3. 准备logger实例
	logger, _ := zap.NewProduction()

	//	2. 重点，注册函数
	registerHandleFunc := func(ctx context.Context, mux *runtime.ServeMux) error {
		// 2-1: 注册hello的服务
		addrHello := "0.0.0.0:9081"
		opts := []grpc.DialOption{grpc.WithInsecure()}
		// 注意参数
		if err := pb.RegisterGreeterHandlerFromEndpoint(ctx, mux, addrHello, opts); err != nil {
			return err
		}

		// 没有错误就表示正常
		return nil
	}

	//	3. 实例化
	grpcServer := rpctools.GrpcGatewayServer{
		Address:      address,
		Mux:          mux,
		RegisterFunc: registerHandleFunc,
		Logger:       logger,
	}
	// start
	grpcServer.Run()
}
