package core

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"strings"
	"time"
)

type Jwt struct {
	Issuer     string
	SigningKey []byte
	Duration   int
}

func NewJwt(issuer string, signingKey []byte, duration int) *Jwt {
	return &Jwt{
		Issuer:     issuer,
		SigningKey: signingKey,
		Duration:   duration,
	}
}

// CreateToken 创建Token
func (j *Jwt) CreateToken(user *User) (string, error) {
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
	claims := &UserClaims{
		UserID:   int64(user.ID),
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.Issuer,
			NotBefore: jwt.NewNumericDate(now.Add(-30 * time.Second)),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(j.Duration) * time.Second)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 3. 返回签署的字符、错误
	return token.SignedString(j.SigningKey)
}

// ParseToken 解析Token
func (j *Jwt) ParseToken(tokenStr string) (*UserClaims, error) {
	if token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	}); err != nil {
		return nil, err
	} else {
		if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
			// 成功
			return claims, nil
		} else {
			err := ErrUnauthorized
			return nil, err
		}
	}
}

// GetRequestToken 从请求头中获取Token
func GetRequestToken(c *gin.Context, tokenHeaderPrefix string) (tokenStr string) {
	// 1. 获取用户传递的Token
	authorizationHeader := c.GetHeader("Authorization")
	if authorizationHeader == "" {
		return ""
	}
	// 获取tokenStr
	if tokenHeaderPrefix != "" {
		tokenStr = strings.TrimPrefix(authorizationHeader, fmt.Sprintf("%s ", tokenHeaderPrefix))
	} else {
		tokenStr = authorizationHeader
	}
	// 返回
	return tokenStr
}
