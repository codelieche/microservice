package store

import (
	"context"
	"errors"
	"github.com/codelieche/microservice/codelieche/filters"
	"github.com/codelieche/microservice/codelieche/utils"
	"github.com/codelieche/microservice/usercenter/core"
	"github.com/codelieche/microservice/usercenter/internal/config"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"log"
	"time"
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
			return nil, core.ErrNotFound
		}
	}
}

func (u *userStore) FindByUsername(ctx context.Context, username string) (*core.User, error) {
	var user = &core.User{}
	if err := u.db.Find(user, "username=?", username).Error; err != nil {
		return nil, err
	} else {
		if user.ID > 0 {
			return user, nil
		} else {
			return nil, core.ErrNotFound
		}
	}
}

func (u *userStore) FindByIdOrUsername(ctx context.Context, idOrUsername string) (*core.User, error) {
	var user = &core.User{}
	if err := u.db.Find(user, "id=? or username=?", idOrUsername, idOrUsername).Error; err != nil {
		return nil, err
	} else {
		if user.ID > 0 {
			return user, nil
		} else {
			return nil, core.ErrNotFound
		}
	}
}

func (u *userStore) Create(ctx context.Context, user *core.User) (*core.User, error) {
	if err := u.db.Create(user).Error; err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

func (u *userStore) SigningToken(ctx context.Context, user *core.User) (signingStr string, err error) {
	// 1. 先做判断
	if user.ID <= 0 || user.Username == "" {
		err := errors.New("传入的用户有误")
		return "", err
	}
	if !user.IsActive {
		err := errors.New("用户已经被禁用")
		return "", err
	}

	// 2. 开始签发
	now := time.Now()
	cfg := config.Config.Web.JWT
	claims := &core.UserClaims{
		UserID:   int64(user.ID),
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.Issuer,
			NotBefore: jwt.NewNumericDate(now.Add(-30 * time.Second)),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(cfg.Duration) * time.Second)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 3. 返回签署的字符、错误
	return token.SignedString([]byte(cfg.Key))
}

func (u *userStore) List(ctx context.Context, offset int, limit int, filterActions ...filters.Filter) (users []*core.User, err error) {
	query := u.db.Model(&core.User{}).
		Select("id, username, nickname, email,phone").
		Offset(offset).Limit(limit)

	if filterActions != nil && len(filterActions) > 0 {
		for _, action := range filterActions {
			log.Println("action, action == nil, action != nil:", action, action == nil, action != nil)
			if action == nil {
				continue
			}
			if action != nil {
				query = action.Filter(query)
			}
		}
	}

	query.Find(&users)
	if query.Error != nil {
		return nil, query.Error
	} else {
		return users, nil
	}
}

func (u *userStore) SetPassword(ctx context.Context, user *core.User, password string) error {
	if user.ID <= 0 || user.Username == "" {
		err := errors.New("传入的用户有误")
		return err
	}
	if !user.IsActive {
		err := errors.New("用户已经被禁用")
		return err
	}
	// 开始设置密码
	if hashedPassword, err := utils.HashPassword(password); err != nil {
		return err
	} else {
		user.Password = hashedPassword
		//	保存密码
		if err := u.db.Model(user).Where("id=?", user.ID).Limit(1).
			Update("password", user.Password).Error; err != nil {
			return err
		} else {
			// 修改密码之后应该把老的token过期一下
			return nil
		}
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

func (u *userStore) Count(ctx context.Context, filterActions ...filters.Filter) (int64, error) {
	var count int64

	query := u.db.Model(&core.User{})

	if filterActions != nil && len(filterActions) > 0 {
		for _, action := range filterActions {
			//log.Println("action, action == nil, action != nil:", action, action == nil, action != nil)
			if action != nil {
				query = action.Filter(query)
			}
		}
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	} else {
		return count, nil
	}
}

func (u *userStore) CountByTeam(ctx context.Context, i int64) (int64, error) {
	//TODO implement me
	panic("implement me")
}
