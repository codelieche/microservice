package userservice

import (
	"github.com/codelieche/microservice/usercenter/internal/datasources"
	"github.com/codelieche/microservice/usercenter/proto/userpb"
	"github.com/codelieche/microservice/usercenter/services"
	"github.com/codelieche/microservice/usercenter/store"
	"google.golang.org/grpc"
	"log"
	"net"
	"testing"
)

func TestUserService_RunServer(t *testing.T) {
	db := datasources.GetMySQLDB()
	userStore := store.NewUserStore(db)

	userService := services.NewUserService(userStore)

	userRpcService := NewUserService(userService)
	log.Println(userRpcService)

	lis, err := net.Listen("tcp", "0:8081")
	if err != nil {
		log.Printf("Failed to listen: %v\n", err)
		t.Error(err)
		return
	}

	s := grpc.NewServer()
	userpb.RegisterUserServiceServer(s, userRpcService)

	log.Println("start grpc server")
	if err := s.Serve(lis); err != nil {
		t.Error(err)
		log.Fatal(err)
	} else {
		log.Printf("Done")
	}
}
