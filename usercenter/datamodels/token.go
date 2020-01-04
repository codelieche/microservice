package datamodels

import "time"

type Token struct {
	BaseFields
	UserID    uint       `json:"user_id"`                                    // 用户ID
	Token     string     `gorm:"NOT NULL;UNIQUE_INDEX;size:60" json:"token"` // Token值，不可为空
	User      *User      `gorm:"ForeignKey:UserID" json:"user"`              // Token对应的用户
	ExpiredAt *time.Time `json:"expired_at"`                                 // 过期时间
	IsActive  bool       `gorm:"type:boolean" json:"is_active"`              // 是否有效
}
