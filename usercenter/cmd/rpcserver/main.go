package main

import (
	"context"
	"fmt"
	"github.com/codelieche/microservice/usercenter/internal/config"
	"github.com/codelieche/microservice/usercenter/internal/datasources"
	"github.com/codelieche/microservice/usercenter/proto/userpb"
	"github.com/codelieche/microservice/usercenter/rpcserver/interceptors"
	"github.com/codelieche/microservice/usercenter/rpcserver/userservice"
	"github.com/codelieche/microservice/usercenter/services"
	"github.com/codelieche/microservice/usercenter/store"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("rpc server")

	logger, err := zap.NewProduction()
	defer logger.Sync()

	if err != nil {
		log.Fatalf("Can not create logger: %s", err.Error())
	}

	db := datasources.GetMySQLDB()
	userStore := store.NewUserStore(db)

	usv := services.NewUserService(userStore)

	userRpcService := userservice.NewUserService(usv, logger)
	log.Println(userRpcService)

	cfg := config.Config.Rpc
	addr := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("Failed to listen: %v\n", err)
		return
	}

	// 注册rpc服务:
	// 实例化拦截器
	in := interceptors.NewAuthInterceptor()
	s := grpc.NewServer(grpc.UnaryInterceptor(in))
	userpb.RegisterUserServiceServer(s, userRpcService)

	// 启动grpc gateway server
	go startGRPCGateway()

	log.Println("start grpc server")
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Done")
	}
}

// startGRPCGateway 启动grpc gateway
func startGRPCGateway() {
	// recover error
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recover error: %v\n", r)
		}
	}()

	// 1. 准备ctx
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// 2. mux
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(
		runtime.MIMEWildcard, &runtime.JSONPb{},
	))

	// 3. RegisterUserServiceHandlerFromEndpoint
	cfg := config.Config.Rpc
	addr := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)
	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := userpb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, addr, opts); err != nil {
		log.Fatalf("cannot start grpc gateway: %v", err)
	}

	// 4. start http server
	addr2 := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port+1)
	if err := http.ListenAndServe(addr2, mux); err != nil {
		log.Fatalf("Cannot listen and server: %v", err)
	}
}
