package store

import (
	"context"
	"encoding/json"
	"github.com/codelieche/microservice/codelieche/utils"
	"github.com/codelieche/microservice/usercenter/core"
	"github.com/codelieche/microservice/usercenter/internal/datasources"
	"log"
	"strings"
	"testing"
)

var noContext = context.TODO()

func TestUserStore(t *testing.T) {
	db := datasources.GetMySQLDB()
	if db == nil {
		panic("获取数据库出错")
	}
	// 清理数据

	store := NewUserStore(db).(*userStore)
	user := &core.User{
		// 测试字符去随机数
		//Team:     strings.ToLower(utils.RandomString(16)),
		Username: strings.ToLower(utils.RandomString(12)),
		Password: utils.RandomString(256),
	}
	//	执行测试
	t.Run("create", testUserCreate(store, user))
	t.Run("find", testUserFind(store, 1))
	t.Run("count", testUserCount(store))

}

func testUserCreate(store *userStore, user *core.User) func(t *testing.T) {
	return func(t *testing.T) {
		if user, err := store.Create(noContext, user); err != nil {
			t.Error(err)
		} else {
			log.Println(user)
		}
	}
}

func testUserFind(store *userStore, i int64) func(t *testing.T) {
	return func(t *testing.T) {
		if user, err := store.Find(noContext, i); err != nil {
			t.Error(err)
		} else {
			data, _ := json.Marshal(user)
			log.Println("获取到用户", string(data))

		}
	}
}

func testUserCount(store *userStore) func(t *testing.T) {
	return func(t *testing.T) {
		if count, err := store.Count(noContext); err != nil {
			t.Error(err)
		} else {
			log.Println("获取到用户总数为：", count)

		}
	}
}
