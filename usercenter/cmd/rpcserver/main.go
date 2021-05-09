package main

import (
	"fmt"
	"github.com/codelieche/microservice/usercenter/internal/config"
	"github.com/codelieche/microservice/usercenter/internal/datasources"
	"github.com/codelieche/microservice/usercenter/proto/userpb"
	"github.com/codelieche/microservice/usercenter/rpcserver/userservice"
	"github.com/codelieche/microservice/usercenter/services"
	"github.com/codelieche/microservice/usercenter/store"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("rpc server")

	db := datasources.GetMySQLDB()
	userStore := store.NewUserStore(db)

	usv := services.NewUserService(userStore)

	userRpcService := userservice.NewUserService(usv)
	log.Println(userRpcService)

	cfg := config.Config.Rpc
	addr := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("Failed to listen: %v\n", err)
		return
	}

	s := grpc.NewServer()
	userpb.RegisterUserServiceServer(s, userRpcService)

	log.Println("start grpc server")
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Done")
	}
}
