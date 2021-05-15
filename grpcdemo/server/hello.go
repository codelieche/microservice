package server

import (
	"context"
	"fmt"
	"github.com/codelieche/microservice/grpcdemo/proto/pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

// HelloService hello service
type HelloService struct {
	logger *zap.Logger // 日志
}

func NewHelloService(logger *zap.Logger) pb.GreeterServer {
	return &HelloService{logger: logger}
}

func (h HelloService) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloResponse, error) {
	msg := fmt.Sprintf("-->:%s", request.Message)
	h.logger.Info(msg)
	response := &pb.HelloResponse{
		Message: msg,
	}

	return response, nil
}

func (h HelloService) Ping(ctx context.Context, empty *emptypb.Empty) (*pb.Pong, error) {
	pong := &pb.Pong{
		Status:  true,
		Message: "pong",
	}
	return pong, nil
}
