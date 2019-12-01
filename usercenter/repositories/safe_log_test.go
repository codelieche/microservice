package repositories

import (
	"log"
	"testing"

	"github.com/codelieche/microservice/usercenter/datasources"

	"github.com/codelieche/microservice/usercenter/datamodels"
)

func TestSafeLogRepository_Save(t *testing.T) {
	// 1. get db
	db := datasources.GetDb()

	db.AutoMigrate(&datamodels.SafeLog{})

	// 2. init safeLog repository
	r := NewSafeLogRepository(db)

	// 创建10条safeLog
	i := 0
	for i < 10 {
		i++
		safeLog := &datamodels.SafeLog{
			Category: uint(i),
			UserID:   uint(i),
			Content:  "消息内容",
			Success:  true,
			Address:  "192.168.1.101",
			Device:   "Chrome",
		}
		if safeLog, err := r.Save(safeLog); err != nil {
			t.Error(err.Error())
		} else {
			log.Println(safeLog.ID, safeLog.Category, safeLog.CreatedAt, safeLog.Content)
		}
	}
}

func TestSafeLogRepository_List(t *testing.T) {
	// 1. get db
	db := datasources.GetDb()

	// 2. init safeLog repository
	r := NewSafeLogRepository(db)

	// 3. list safeLog
	haveNext := true
	offset := 0
	limit := 5
	for haveNext {
		if safeLogs, err := r.List(offset, limit); err != nil {
			t.Error(err.Error())
			haveNext = false
		} else {
			if len(safeLogs) == limit && limit > 0 {
				haveNext = true
				offset += limit
			} else {
				haveNext = false
			}
			// 输出
			for _, safeLog := range safeLogs {

				log.Println(safeLog.ID, safeLog.Category, safeLog.Content)
			}
		}
	}
}
