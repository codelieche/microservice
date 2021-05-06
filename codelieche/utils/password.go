package utils

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strings"
	"time"
)

var r *rand.Rand

func RandomString(len int) string {
	// 1. 随机数种子
	if r == nil {
		r = rand.New(rand.NewSource(time.Now().Unix()))
	}

	// 2. 开始随机取字符
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26)

		switch i % 3 {
		case 0:
			b = b + 97
		case 1:
			b = b + 65
		case 2:
			b = b%10 + 48
		default:
			b += 65
		}
		bytes[i] = byte(b)
	}

	// 3. 返回字符
	return string(bytes)
}

func HashPassword(password string) (string, error) {
	// 1. 取出两边的空格
	password = strings.TrimSpace(password)

	// 2. 判断密码是否为空
	if password == "" {
		return "", errors.New("password is null")
	}

	// 3. 加密密码
	if hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8); err != nil {
		return "", err
	} else {
		return string(hashedPassword), nil
	}
}

func CheckHashedPassword(hashedPassword, password string) (bool, error) {
	// 1. 密码都需要取出两边的空格
	hashedPassword = strings.TrimSpace(hashedPassword)
	password = strings.TrimSpace(password)

	// 2. 检查是否为空
	if hashedPassword == "" || password == "" {
		return false, errors.New("密码为空，不可核对")
	}

	// 3. 开始检查密码
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		// 如果匹配不成功，那么错误是：crypto/bcrypt: hashedPassword is not the hash of the given password
		return false, err
	} else {
		return true, nil
	}
}
