package datamodels

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Ticket struct {
	gorm.Model
	Name        string    `gorm:"type:varchar(100);UNIQUE_INDEX" json:"name"` // Ticket的名字: hash(session id)
	Session     string    `gorm:"size:100" json:"session"`                    // SessionID
	ReturnUrl   string    `gorm:"size:512" json:"return_url" `                // 跳转的URL
	IsActive    bool      `gorm:"type:boolean" json:"is_active"`              // 是否有效：每个ticket只能用一次
	TimeExpired time.Time `json:"time_expired"`                               // 过期时间：推荐为60秒
}