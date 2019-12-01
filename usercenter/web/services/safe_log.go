package services

import (
	"github.com/codelieche/microservice/usercenter/datamodels"
	"github.com/codelieche/microservice/usercenter/repositories"
)

//SafeLog Service Interface
type SafeLogService interface {
	Save(safeLog *datamodels.SafeLog) (*datamodels.SafeLog, error)
	GetById(id int64) (safeLog *datamodels.SafeLog, err error)
	List(offset int, limit int) (safeLogs []*datamodels.SafeLog, err error)
	ListByUser(userID int, offset int, limit int) (safeLogs []*datamodels.SafeLog, err error)
}

// 实例化SafeLog Service
func NewSafeLogService(repo repositories.SafeLogRepository) SafeLogService {
	return &safeLogService{repo: repo}
}

//safeLog Service
type safeLogService struct {
	repo repositories.SafeLogRepository
}

func (s *safeLogService) Save(safeLog *datamodels.SafeLog) (*datamodels.SafeLog, error) {
	return s.repo.Save(safeLog)
}

func (s *safeLogService) GetById(id int64) (safeLog *datamodels.SafeLog, err error) {
	return s.repo.Get(id)
}

// 获取SafeLog列表
func (s *safeLogService) List(offset int, limit int) (safeLogs []*datamodels.SafeLog, err error) {
	return s.repo.List(offset, limit)
}

// 获取用户SafeLog列表
func (s *safeLogService) ListByUser(userID int, offset int, limit int) (safeLogs []*datamodels.SafeLog, err error) {
	return s.repo.ListByUser(userID, offset, limit)
}
