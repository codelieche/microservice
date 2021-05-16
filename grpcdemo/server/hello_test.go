package server

import (
	"github.com/codelieche/microservice/grpcdemo/intercepetor"
	"github.com/codelieche/microservice/grpcdemo/proto/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"testing"
)

func getGrpcServerListener() (net.Listener, error) {
	// 测试启动Hello Service
	addr := "0.0.0.0:9081"

	// 1. net listen
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		// 直接抛出错误
		log.Fatal(err)
		return nil, err
	}
	return lis, nil
}

func TestNewHelloService(t *testing.T) {
	// 测试启动Hello Service
	addr := "0.0.0.0:9081"

	// 1. net listen
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		t.Error(err)
		return
	}

	// 2. grpc server
	s := grpc.NewServer()

	// 3. register grpc server
	// 3-1: 实例化server
	logger, _ := zap.NewProduction()
	helloServer := NewHelloService(logger)

	// 3-2: register grpc server
	pb.RegisterGreeterServer(s, helloServer)

	// 4. start grpc server
	if err := s.Serve(lis); err != nil {
		t.Error(err)
	}
}

func TestHelloService_Interceptor(t *testing.T) {
	// 1. 准备 lis
	lis, err := getGrpcServerListener()
	if err != nil {
		t.Error(err)
		return
	}

	// 2. 准备ServerOption
	var opts []grpc.ServerOption
	opts = append(opts, grpc.UnaryInterceptor(intercepetor.RequestInfoPrint))

	// 3. 实例化grpc server
	s := grpc.NewServer(opts...)

	// 4. 注入grpc服务
	logger, _ := zap.NewProduction()
	greeterServer := NewHelloService(logger)

	pb.RegisterGreeterServer(s, greeterServer)

	// 5. 开始监听grpc server
	if err := s.Serve(lis); err != nil {
		t.Error(err)
		return
	}
}
