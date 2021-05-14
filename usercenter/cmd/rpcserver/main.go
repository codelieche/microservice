package main

import (
	"fmt"
	"github.com/codelieche/microservice/codelieche/rpctools"
	"github.com/codelieche/microservice/usercenter/core"
	"github.com/codelieche/microservice/usercenter/internal/config"
	"github.com/codelieche/microservice/usercenter/internal/datasources"
	"github.com/codelieche/microservice/usercenter/proto/userpb"
	"github.com/codelieche/microservice/usercenter/rpcserver/interceptors"
	"github.com/codelieche/microservice/usercenter/rpcserver/userservice"
	"github.com/codelieche/microservice/usercenter/services"
	"github.com/codelieche/microservice/usercenter/store"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main01() {
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
	jwtCfg := config.Config.Web.JWT
	j := core.NewJwt(jwtCfg.Issuer, []byte(jwtCfg.Key), jwtCfg.Duration)
	in := interceptors.NewAuthInterceptor(j)
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

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	// 1. 准备grpc server config需要的参数
	// 1-1： 日志实例
	logger, err := zap.NewProduction()
	defer logger.Sync()
	if err != nil {
		log.Fatalf("Can not create logger: %s", err.Error())
	}

	// 1-2：准备grpc server Option
	var opts []grpc.ServerOption
	// 1-2-1：实例化auth拦截器
	jwtCfg := config.Config.Web.JWT
	j := core.NewJwt(jwtCfg.Issuer, []byte(jwtCfg.Key), jwtCfg.Duration)
	in := interceptors.NewAuthInterceptor(j)
	// append拦截器
	opts = append(opts, grpc.UnaryInterceptor(in))

	// 1-3：注册grpc服务的函数
	// 注册grpc服务函数
	registerFunc := func(s *grpc.Server) {
		// 需要注入几个服务就注入几个
		// 1. user grpc service
		// 1-1：准备store
		db := datasources.GetMySQLDB()
		userStore := store.NewUserStore(db)
		// 1-2：准备service
		usv := services.NewUserService(userStore)
		// 实例化user grpc service
		// 1-3：准备grpc service
		userRpcService := userservice.NewUserService(usv, logger)
		// 1-4：执行register操作
		userpb.RegisterUserServiceServer(s, userRpcService)

		//	2. team grpc service
	}

	// 2：实例化grpc server config
	grpcCfg := &rpctools.GrpcServerConfig{
		Name:         "user",
		Opts:         opts,
		Address:      fmt.Sprintf("%s:%d", config.Config.Rpc.Address, config.Config.Rpc.Port),
		RegisterFunc: registerFunc,
		Logger:       logger,
	}

	// 启动grpc gateway server
	go startGRPCGateway()

	// 3. 启动grpc服务
	if err := rpctools.RunGRPCServer(grpcCfg); err != nil {
		logger.Fatal("grpc server error", zap.Error(err))
	}
}
