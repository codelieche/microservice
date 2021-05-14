package main

import (
	"context"
	"fmt"
	"github.com/codelieche/microservice/usercenter/internal/config"
	"github.com/codelieche/microservice/usercenter/proto/userpb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"log"
	"net/http"
)

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

	// 3-1: User Service
	cfg := config.Config.Rpc
	addr := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)
	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := userpb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, addr, opts); err != nil {
		log.Fatalf("cannot start grpc gateway: %v", err)
	}

	// 3-2: hello service

	// 4. start http server
	addr2 := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port+1)
	if err := http.ListenAndServe(addr2, mux); err != nil {
		log.Fatalf("Cannot listen and server: %v", err)
	}
}
