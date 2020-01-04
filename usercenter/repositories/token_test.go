package repositories

import (
	"log"
	"testing"
	"time"

	"github.com/codelieche/microservice/usercenter/datamodels"
	"github.com/codelieche/microservice/usercenter/datasources"
)

func TestTokenRepository_Create(t *testing.T) {
	// 1. get db
	db := datasources.GetDb()

	// 2. init token repository
	r := NewTokenRepository(db)

	// 3. create token
	token := &datamodels.Token{
		UserID:    1,
		Token:     "",
		User:      nil,
		ExpiredAt: nil,
		IsActive:  true,
	}
	if token, err := r.Create(token); err != nil {
		t.Error(err)
	} else {
		log.Println(token)
	}
}

func TestTokenRepository_Create02(t *testing.T) {
	// 1. get db
	db := datasources.GetDb()

	// 2. init token repository
	r := NewTokenRepository(db)

	// 3. create token
	expiredAt := time.Now().Add(time.Hour * 24 * 360)
	token := &datamodels.Token{
		UserID:    1,
		Token:     "",
		User:      nil,
		ExpiredAt: &expiredAt,
		IsActive:  true,
	}
	if token, err := r.Create(token); err != nil {
		t.Error(err)
	} else {
		log.Println(token)
	}
}
