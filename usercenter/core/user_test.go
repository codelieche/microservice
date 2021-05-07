package core

import (
	"github.com/codelieche/microservice/codelieche/utils"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"strings"
	"testing"
	"time"
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

func TestUserChaims(t *testing.T) {
	signingKey := []byte("thisispassword")
	// 2021-05-06 11:00:00
	now := time.Unix(1620270000, 0)
	claims := UserClaims{
		UserID:   123,
		Username: "codelieche",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "codelieche",
			NotBefore: jwt.NewNumericDate(now.Add(-30 * time.Second)),
			ExpiresAt: jwt.NewNumericDate(now.Add(12 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	if ss, err := token.SignedString(signingKey); err != nil {
		t.Error(err)
	} else {
		log.Println(ss)
	}

}

func TestUserChaims2(t *testing.T) {
	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjMsInVzZXJuYW1lIjoiY29kZWxpZWNoZSIsImlzcyI6ImNvZGVsaWVjaGUiLCJleHAiOjE2MjAzMTMyMDAsIm5iZiI6MTYyMDI2OTk3MH0.y8AYJwtu6kZ0FHxfgsrZjkZG_Iinxi2bTmUEy1_yNmU"
	// 2021-05-06 11:00:00 时间戳是：1620270000
	jwt.TimeFunc = func() time.Time {
		return time.Unix(1620270000+1000, 0)
	}
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("thisispassword"), nil
	})
	if err != nil {
		t.Error("parse claims error:", err)
		return
	} else {
		log.Println(token)
		if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
			log.Printf("%v:%v:%v:%v", claims.Username, claims.UserID, claims.RegisteredClaims.Issuer, claims.RegisteredClaims.ExpiresAt)
		}
	}

}
