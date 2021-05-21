package core

import (
	"context"
	"gorm.io/gorm"
)

type (
	Team struct {
		gorm.Model
		Name        string `gorm:"size:64;index" json:"name"`   // 团队名称，不唯一
		Description string `gorm:"size:512" json:"description"` // 团队描述
	}

	TeamUser struct {
		gorm.Model
		//Team     Team   `gorm:"foreignKey:TeamID"`                                     // 团队
		//User     User   `gorm:"foreignKey:UserID"`                                     // 用户

		TeamID   int    `gorm:"uniqueIndex:idx_team_user;uniqueIndex:idx_team_username" json:"team_id"` // 团队ID
		UserID   int    `gorm:"uniqueIndex:idx_team_user" json:"user_id"`                               // 用户ID
		Username string `gorm:"size:60;uniqueIndex:idx_team_username" json:"username"`                  // 在团队里面的username
		Nickname string `gorm:"size:60" json:"nickname"`                                                // 用户昵称
		IsActive bool   `gorm:"type:boolean;default:true" json:"is_active"`                             // 状态
		IsAdmin  bool   `gorm:"type:boolean;default:false" json:"is_admin"`                             // 是否是这个团队的管理员
	}

	// TeamStore 团队操作数据库的接口
	TeamStore interface {
		// Find 根据ID获取团队
		Find(context.Context, int64) (*Team, error)

		// Create 创建团队
		Create(ctx context.Context, team *Team) error

		// Update 更新团队
		Update(ctx context.Context, team *Team) error

		// Delete 删除团队
		Delete(ctx context.Context, team *Team) error

		// DeleteByID 根据id删除团队
		DeleteByID(ctx context.Context, id int64) error
	}
)
