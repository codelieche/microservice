package userservice

import (
	"context"
	"fmt"
	"github.com/codelieche/microservice/usercenter/internal/config"
	"github.com/codelieche/microservice/usercenter/proto/userpb"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"log"
	"testing"
	"time"
)

func TestUserService_GetUser(t *testing.T) {
	// 1. connect
	cfg := config.Config.Rpc
	addr := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)
	//conn, err := grpc.Dial(addr) // grpc: no transport security set
	conn, err := grpc.Dial(addr, grpc.WithInsecure())

	if err != nil {
		t.Error(err)
		return
	}

	// 2. get user rpc client
	userClient := userpb.NewUserServiceClient(conn)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	r := &userpb.GetUserRequest{
		Id:       1,
		Username: "test001",
	}
	if user, err := userClient.GetUser(ctx, r); err != nil {
		t.Error(err)
	} else {
		log.Println("user info ====>:", user)
		if data, err := proto.Marshal(user); err != nil {
			t.Error(err)
		} else {
			log.Printf("%s\n", data)
			log.Printf("%v\n", data)
		}
	}
}
