package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/codelieche/microservice/usercenter/core"
	"github.com/codelieche/microservice/usercenter/internal/config"
	"github.com/golang-jwt/jwt/v4"
)

func NewUserService(store core.UserStore) core.UserService {
	return &userService{store: store}
}

type userService struct {
	store core.UserStore
}

func (s *userService) Login(ctx context.Context, username string, password string, category string) (token string, err error) {
	//	1：获取用户
	if username == "" || password == "" {
		err = errors.New("输入的用户名或密码不合法")
		return "", err
	}
	if category != "" && category != "username" {
		err = errors.New("暂时只支持通过用户名登录")
		return "", err
	}

	user, err := s.FindByUsername(ctx, username)
	if err != nil {
		if err == core.ErrNotFound {
			err = fmt.Errorf("用户不存在")
		}
		return "", err
	}

	// 2：检查用户密码是否正确
	if ok, err := user.CheckPassword(password); err != nil || !ok {
		err = fmt.Errorf("用户名或者密码错误")
		return "", err
	}

	//	3. 返回JWT Token
	if token, err := s.SigningToken(ctx, user); err != nil {
		return "", err
	} else {
		return token, err
	}
}

func (s *userService) Find(ctx context.Context, i int64) (*core.User, error) {
	user, err := s.store.Find(ctx, i)
	return user, err
}

func (s *userService) FindByUsername(ctx context.Context, username string) (*core.User, error) {
	user, err := s.store.FindByUsername(ctx, username)
	return user, err
}

func (s *userService) FindByIdOrUsername(ctx context.Context, idOrUsername string) (*core.User, error) {
	user, err := s.store.FindByIdOrUsername(ctx, idOrUsername)
	return user, err
}

func (s *userService) Create(ctx context.Context, user *core.User) (*core.User, error) {
	if user, err := s.store.Create(ctx, user); err != nil {
		return nil, err
	} else {
		return user, err
	}
}

func (s *userService) SigningToken(ctx context.Context, user *core.User) (signingStr string, err error) {
	return s.store.SigningToken(ctx, user)
}

func (s *userService) ParseToken(ctx context.Context, tokenStr string) (claims *core.UserClaims, err error) {
	cfg := config.Config.Web.JWT

	if token, err := jwt.ParseWithClaims(tokenStr, &core.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Key), nil
	}); err != nil {
		return nil, err
	} else {
		if claims, ok := token.Claims.(*core.UserClaims); ok && token.Valid {
			// 成功
			return claims, nil
		} else {
			err := core.ErrUnauthorized
			return nil, err
		}
	}
}

func (s *userService) List(ctx context.Context, offset int, limit int) (users []*core.User, err error) {
	return s.store.List(ctx, offset, limit)
}

func (s *userService) Count(ctx context.Context) (int64, error) {
	return s.store.Count(ctx)
}

func (s *userService) SetPassword(ctx context.Context, user *core.User, password string) error {
	return s.store.SetPassword(ctx, user, password)
}
