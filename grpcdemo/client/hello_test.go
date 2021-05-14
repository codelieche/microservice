package client

import (
	"context"
	"github.com/codelieche/microservice/grpcdemo/proto/pb"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"log"
	"testing"
	"time"
)

func TestGreeter_SayHello(t *testing.T) {
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
