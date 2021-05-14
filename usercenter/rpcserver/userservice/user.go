package userservice

import (
	"context"
	"errors"
	"fmt"
	"github.com/codelieche/microservice/usercenter/core"
	"github.com/codelieche/microservice/usercenter/proto/userpb"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"log"
)

// UserService 用户服务
type UserService struct {
	service core.UserService
	//userpb.UnimplementedUserServiceServer
	logger *zap.Logger
}

// NewUserService 创建UserService
func NewUserService(service core.UserService, logger *zap.Logger) userpb.UserServiceServer {
	return &UserService{service: service, logger: logger}
}

func (s *UserService) Login(ctx context.Context, request *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	s.logger.Info(fmt.Sprintf("我收到登录请求：%s", request))
	// 1. 根据传递的用户名和密码执行登录
	var user *userpb.User
	token, err := s.service.Login(ctx, request.Username, request.Password, request.Category)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	} else {
		// 获取用户
		u, err := s.service.FindByUsername(ctx, request.Username)
		if err != nil {
			return nil, status.Errorf(codes.NotFound, err.Error())
		} else {
			user = &userpb.User{
				Id:          int64(u.ID),
				Username:    u.Username,
				Nickname:    u.Nickname,
				Email:       u.Email,
				Photo:       u.Phone,
				IsSuperuser: u.IsSuperuser,
				IsActive:    u.IsActive,
			}
		}
	}

	// 2. 返回token
	return &userpb.LoginResponse{
		Token:    token,
		Userinfo: user,
	}, nil
}

// GetUser 获取用户信息
func (s *UserService) GetUser(ctx context.Context, in *userpb.GetUserRequest) (*userpb.User, error) {
	var user *core.User
	var err error
	var u *userpb.User
	if in.Id > 0 {
		user, err = s.service.Find(ctx, in.Id)
	} else if in.Username != "" {
		user, err = s.service.FindByUsername(ctx, in.Username)
	} else {
		err = errors.New("传递的用户id或者username有误")
		return nil, err
	}

	log.Println("user info ====>:", user)

	if user != nil && user.ID > 0 && user.Username != "" {
		u = &userpb.User{
			Id:          int64(user.ID),
			Username:    user.Username,
			Nickname:    user.Nickname,
			Email:       user.Email,
			Photo:       user.Phone,
			IsSuperuser: user.IsSuperuser,
			IsActive:    user.IsActive,
		}
	} else {
		log.Printf("获取到的用户有误：%v", user)
	}

	return u, err
}

func (s *UserService) ListUser(ctx context.Context, request *userpb.ListRequest) (*userpb.ListResponse, error) {

	//users, err := s.service.List(ctx, int(request.Page), int(request.PageSize))
	if request.PageSize <= 0 {
		request.PageSize = 10
	}
	users, err := s.service.List(ctx, int(request.Page), 10)
	count, err := s.service.Count(ctx)

	if err != nil {
		return nil, err
	}

	var results []*anypb.Any
	if users != nil && len(users) > 0 {
		for _, u := range users {
			a, _ := anypb.New(&userpb.User{
				Id:          int64(u.ID),
				Username:    u.Username,
				Nickname:    u.Nickname,
				Email:       u.Email,
				Photo:       u.Phone,
				IsSuperuser: u.IsSuperuser,
				IsActive:    u.IsActive,
			})
			results = append(results, a)
		}
	}

	return &userpb.ListResponse{
		Count:   count,
		Results: results,
	}, nil
}
