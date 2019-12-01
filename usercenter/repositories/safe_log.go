package repositories

import (
	"errors"

	"github.com/codelieche/microservice/usercenter/common"
	"github.com/codelieche/microservice/usercenter/datamodels"
	"github.com/jinzhu/gorm"
)

type SafeLogRepository interface {
	// 保存SafeLog
	Save(safeLog *datamodels.SafeLog) (*datamodels.SafeLog, error)
	// 获取SafeLog信息
	Get(id int64) (*datamodels.SafeLog, error)
	// 获取SafeLog的列表
	List(offset int, limit int) ([]*datamodels.SafeLog, error)
}

func NewSafeLogRepository(db *gorm.DB) SafeLogRepository {
	return &safeLogRepository{db: db}
}

type safeLogRepository struct {
	db *gorm.DB
}

// 保存SafeLog
func (r *safeLogRepository) Save(safeLog *datamodels.SafeLog) (*datamodels.SafeLog, error) {
	if safeLog.ID > 0 {
		// 是更新操作
		// 安全日志不支持修改查找
		err := errors.New("ID>0SafeLog的，不支持更新操作")
		return nil, err
	} else {
		// 判断是否指定了UserID
		if safeLog.UserID < 1 {
			err := errors.New("未指定用户")
			return nil, err
		}
		// 是创建操作
		if err := r.db.Create(safeLog).Error; err != nil {
			return nil, err
		} else {
			return safeLog, nil
		}

	}
}

// 获取SafeLog的列表
func (r *safeLogRepository) List(offset int, limit int) (safeLogs []*datamodels.SafeLog, err error) {
	query := r.db.Model(&datamodels.SafeLog{}).Offset(offset).Limit(limit).Find(&safeLogs)
	if query.Error != nil {
		return nil, query.Error
	} else {
		return safeLogs, nil
	}
	return
}

// 根据ID获取SafeLog
func (r *safeLogRepository) Get(id int64) (safeLog *datamodels.SafeLog, err error) {

	safeLog = &datamodels.SafeLog{}
	r.db.First(safeLog, "id = ?", id)
	if safeLog.ID > 0 {
		return safeLog, nil
	} else {
		return nil, common.NotFountError
	}
}
