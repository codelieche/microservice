package client

import (
	"context"
	"github.com/codelieche/microservice/grpcdemo/intercepetor"
	"github.com/codelieche/microservice/grpcdemo/proto/pb"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"testing"
	"time"
)

func getGreeterServerClient(opts []grpc.DialOption) (pb.GreeterClient, error) {
	// 1. connect
	addr := "127.0.0.1:9081"
	//conn, err := grpc.Dial(addr) // grpc: no transport security set
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(addr, opts...)

	if err != nil {
		return nil, err
	}

	// 2. get grpc client
	greeterClient := pb.NewGreeterClient(conn)

	// 3. return
	return greeterClient, nil
}

func TestGreeter_Ping(t *testing.T) {
	// 1. connect
	addr := "127.0.0.1:9081"
	//conn, err := grpc.Dial(addr) // grpc: no transport security set
	conn, err := grpc.Dial(addr, grpc.WithInsecure())

	if err != nil {
		t.Error(err)
		return
	}

	// 2. get user rpc client
	userClient := pb.NewGreeterClient(conn)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	r := &emptypb.Empty{}
	if response, err := userClient.Ping(ctx, r); err != nil {
		t.Error(err)
		return
	} else {
		log.Println(response)
	}
}

func TestGreeter_SayHello(t *testing.T) {
	// 1. connect grpc server
	userClient, err := getGreeterServerClient(nil)
	if err != nil {
		t.Error(err)
		return
	}

	// 2. ready context
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	// 3. request
	r := &pb.HelloRequest{
		Message: "this is test from golang",
	}
	if response, err := userClient.SayHello(ctx, r); err != nil {
		t.Error("say hello error：", err)
	} else {
		log.Println("response is ====>:", response)
		if data, err := proto.Marshal(response); err != nil {
			t.Error(err)
		} else {
			log.Printf("%s\n", data)
			log.Printf("%v\n", data)
		}
	}
}

func TestGreeter_Ping_With_Interceptor(t *testing.T) {
	// 1. get grpc client
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithUnaryInterceptor(intercepetor.RerequestBeforePrint))
	client, err := getGreeterServerClient(opts)

	if err != nil {
		t.Error(err)
		return
	}

	// 2. ready context
	// 2-1: 准备个可取消的context
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	// 2-2：准备metadata
	md := metadata.New(map[string]string{
		"username": "root",
		"password": "password",
		"token":    "token value",
	})
	// 2-3：往context中注入metadata
	ctx = metadata.NewOutgoingContext(ctx, md)

	r := &emptypb.Empty{}
	if response, err := client.Ping(ctx, r); err != nil {
		t.Error(err)
		return
	} else {
		log.Println(response)
	}
}
