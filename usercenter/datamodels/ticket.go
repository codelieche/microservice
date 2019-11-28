package datamodels

import (
	"time"
)

type Ticket struct {
	ID        uint       `gorm:"type:bigint(20) unsigned auto_increment;not null;primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
	//gorm.Model
	UserID      uint      `gorm:"INDEX" json:"user_id""`                      // 用户的ID
	User        *User     `gorm:"ForeignKey:UserID" json:"user"`              // 使用UserID作为外键
	Name        string    `gorm:"type:varchar(100);UNIQUE_INDEX" json:"name"` // Ticket的名字: hash(session id)
	Session     string    `gorm:"size:100" json:"session"`                    // SessionID
	ReturnUrl   string    `gorm:"size:512" json:"return_url" `                // 跳转的URL
	IsActive    bool      `gorm:"type:boolean" json:"is_active"`              // 是否有效：每个ticket只能用一次
	TimeExpired time.Time `json:"time_expired"`                               // 过期时间：推荐为60秒
}
