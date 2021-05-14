package server

import (
	"github.com/codelieche/microservice/grpcdemo/proto/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"testing"
)

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
