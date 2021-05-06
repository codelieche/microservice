package services

import (
	"context"
	"github.com/codelieche/microservice/usercenter/core"
	"github.com/codelieche/microservice/usercenter/internal/config"
	"github.com/codelieche/microservice/usercenter/store"
	"reflect"
	"testing"
)

func Test_userService_Create(t *testing.T) {

	db := store.GetMySQLDB(config.Config)
	userStore := store.NewUserStore(db)

	type args struct {
		ctx  context.Context
		user *core.User
	}

	tests := []struct {
		name    string
		store   core.UserStore
		args    args
		want    *core.User
		wantErr bool
	}{
		{
			name:  "创建用户",
			store: userStore,
			args: args{
				ctx: context.Background(),
				user: &core.User{
					Username: "hello001",
					Nickname: "hello001",
					Password: "password",
				},
			},
			want: &core.User{
				Username: "hello001",
				Nickname: "hello001",
				Password: "password",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &userService{
				store: tt.store,
			}
			got, err := s.Create(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Username, tt.want.Username) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_Find(t *testing.T) {
	type fields struct {
		store core.UserStore
	}
	type args struct {
		ctx context.Context
		i   int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *core.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &userService{
				store: tt.fields.store,
			}
			got, err := s.Find(tt.args.ctx, tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Find() got = %v, want %v", got, tt.want)
			}
		})
	}
}
