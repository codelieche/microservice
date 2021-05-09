package userservice

import (
	"context"
	"errors"
	"github.com/codelieche/microservice/usercenter/core"
	"github.com/codelieche/microservice/usercenter/proto/userpb"
	"log"
)

// UserService 用户服务
type UserService struct {
	service core.UserService
	userpb.UnimplementedUserServiceServer
}

// NewUserService 创建UserService
func NewUserService(service core.UserService) *UserService {
	return &UserService{service: service}
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
			Id:       int64(user.ID),
			Username: user.Username,
			Nickname: user.Nickname,
			Email:    user.Email,
			Photo:    user.Phone,
		}
	} else {
		log.Printf("获取到的用户有误：%v", user)
	}

	return u, err
}
