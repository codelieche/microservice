package datamodels

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Ticket struct {
	gorm.Model
	Name        string    `json:"name" gorm:"size:100,UNIQUE_INDEX"` // Ticket的名字: hash(session id)
	Session     string    `json:"session" gorm:"size:100"`           // SessionID
	ReturnUrl   string    `json:"return_url" gorm:"size:512"`        // 跳转的URL
	IsActive    bool      `json:"is_active"`                         // 是否有效：每个ticket只能用一次
	TimeExpired time.Time `json:"time_expired"`                      // 过期时间：推荐为60秒
}
