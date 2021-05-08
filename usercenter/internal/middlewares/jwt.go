package middlewares

import (
	"github.com/codelieche/microservice/usercenter/core"
	"github.com/codelieche/microservice/usercenter/internal/config"
	"github.com/gin-gonic/gin"
	"time"
)

func JwtAuth(cfg *config.JWT) gin.HandlerFunc {
	j := core.NewJwt(cfg.Issuer, []byte(cfg.Key), cfg.Duration)
	return func(c *gin.Context) {
		// 1. 获取token并判断token是否为空
		token := core.GetRequestToken(c, config.JwtTokenHeaderPrefix)
		if token == "" {
			// 判断不需要jwt auth的接口地址
			if v, exist := config.JwtAuthBlackUrlPathMap[c.Request.URL.Path]; !v || !exist {
				c.AbortWithStatusJSON(401, gin.H{
					"code":    401,
					"message": "请传入Token",
				})
				return
			}
			// 白名单中，无需检测
		} else {
			// 2. 解析token
			if claims, err := j.ParseToken(token); err != nil {
				c.AbortWithStatusJSON(401, gin.H{
					"code":    401,
					"message": err.Error(),
				})
				return
			} else {
				if claims.ExpiresAt.Sub(time.Now()).Seconds() < 0 {
					c.AbortWithStatusJSON(401, gin.H{
						"code":    401,
						"message": "Token已经过期",
					})
					return
				}

				//	3. 根据username获取Token
				// 设置用户名和用户ID
				c.Set("userID", claims.UserID)
				c.Set("username", claims.Username)
			}
		}

		c.Next()
	}
}
