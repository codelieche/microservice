package utils

import (
	"fmt"
	"log"
	"testing"
)

func TestRandomString(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Println(RandomString(10))
	}
}

func TestEncodePassword(t *testing.T) {
	source := RandomString(16)

	if hashedPassword, err := HashPassword(source); err != nil {
		t.Error(err)
	} else {
		log.Printf("原始密码是：%s, 加密后密码为:%s", source, hashedPassword)
	}
}

func TestCheckEncodePassword(t *testing.T) {
	hashedPassword := "$2a$08$YrDpk0mnaHGhVSml4tBUuOnqHTHRa00PvB8dnIyHxIcwZkbZFvJuq"
	source := "codelieche"

	if success, err := CheckHashedPassword(hashedPassword, source); err != nil {
		t.Error(err)
	} else {
		if success {
			log.Printf("password is ok: %s", source)
		}
	}
}
