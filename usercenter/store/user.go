package store

import (
	"context"
	"github.com/codelieche/microservice/usercenter/core"
	"github.com/codelieche/microservice/usercenter/internal/config"
	"gorm.io/gorm"
)

func NewUserStore(db *gorm.DB) core.UserStore {
	// 自动获取db
	if db == nil {
		db = GetMySQLDB(config.Config)
	}

	return &userStore{
		db: db,
	}
}

type userStore struct {
	db *gorm.DB
}

func (u *userStore) Find(ctx context.Context, i int64) (*core.User, error) {
	var user = &core.User{}
	if err := u.db.Find(user, "id=?", i).Error; err != nil {
		return nil, err
	} else {
		if user.ID > 0 {
			return user, nil
		} else {
			return nil, core.NotFoundErr
		}
	}
}

func (u *userStore) FindByUsername(ctx context.Context, s string) (*core.User, error) {
	var user = &core.User{}
	if err := u.db.Find(user).Where("username=?", s).Error; err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

func (u *userStore) Create(ctx context.Context, user *core.User) (*core.User, error) {
	if err := u.db.Create(user).Error; err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

func (u *userStore) Update(ctx context.Context, user *core.User) error {
	//TODO implement me
	panic("implement me")
}

func (u *userStore) Delete(ctx context.Context, user *core.User) error {
	//TODO implement me
	panic("implement me")
}

func (u *userStore) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := u.db.Model(&core.User{}).Count(&count).Error; err != nil {
		return 0, err
	} else {
		return count, nil
	}
}

func (u *userStore) CountByTeam(ctx context.Context, i int64) (int64, error) {
	//TODO implement me
	panic("implement me")
}
