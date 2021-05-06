package core

import (
	"github.com/codelieche/microservice/codelieche/utils"
	"log"
	"strings"
	"testing"
)

func TestUser(t *testing.T) {

	u := &User{
		Username: strings.ToLower(utils.RandomString(8)),
	}
	t.Run("ValidateUsername", testUserUsernameValidate(u))
}

func testUserUsernameValidate(u *User) func(t *testing.T) {

	return func(t *testing.T) {
		if result, err := u.ValidateUsername(); err != nil {
			t.Error(err)
		} else {
			if result {
				log.Println("验证用户成功:", u.Username)
			} else {
				t.Error("验证用户名失败：", u.Username)
			}
		}
	}
}
