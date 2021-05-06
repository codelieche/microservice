package services

import (
	"context"
	"github.com/codelieche/microservice/usercenter/core"
)

func NewUserService(store core.UserStore) core.UserService {
	return &userService{store: store}
}

type userService struct {
	store core.UserStore
}

func (s *userService) Find(ctx context.Context, i int64) (*core.User, error) {
	user, err := s.store.Find(ctx, i)
	return user, err
}

func (s *userService) Create(ctx context.Context, user *core.User) (*core.User, error) {
	if user, err := s.store.Create(ctx, user); err != nil {
		return nil, err
	} else {
		return user, err
	}
}
