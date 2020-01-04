package services

import (
	"github.com/codelieche/microservice/usercenter/datamodels"
	"github.com/codelieche/microservice/usercenter/repositories"
)

type TokenService interface {
	// 创建
	Create(token *datamodels.Token) (*datamodels.Token, error)
	// 获取
	Get(id int64) (token *datamodels.Token, err error)
	GetByToken(token string) (*datamodels.Token, error)
	GetByIdOrToken(idOrToken string) (token *datamodels.Token, err error)
	// 删除
	Delete(token *datamodels.Token) (err error)
	// 修改
	Update(token *datamodels.Token, fields map[string]interface{}) (*datamodels.Token, error)
	UpdateByID(id int64, fields map[string]interface{}) (*datamodels.Token, error)
	//获取列表
	List(offset int, limit int) (tokens []*datamodels.Token, err error)
}

func NewTokenService(repo repositories.TokenRepository) TokenService {
	return &tokenService{repo: repo}
}

type tokenService struct {
	repo repositories.TokenRepository
}

func (s *tokenService) Create(token *datamodels.Token) (*datamodels.Token, error) {
	return s.repo.Create(token)
}

func (s *tokenService) Get(id int64) (token *datamodels.Token, err error) {
	return s.repo.Get(id)
}

func (s *tokenService) GetByToken(token string) (*datamodels.Token, error) {
	return s.repo.GetByToken(token)
}

func (s *tokenService) GetByIdOrToken(idOrToken string) (token *datamodels.Token, err error) {
	return s.repo.GetByIdOrToken(idOrToken)
}

func (s *tokenService) Delete(token *datamodels.Token) (err error) {
	return s.repo.Delete(token)
}

func (s *tokenService) Update(token *datamodels.Token, fields map[string]interface{}) (*datamodels.Token, error) {
	return s.repo.Update(token, fields)
}

func (s *tokenService) UpdateByID(id int64, fields map[string]interface{}) (*datamodels.Token, error) {
	return s.repo.UpdateByID(id, fields)
}

func (s *tokenService) List(offset int, limit int) (tokens []*datamodels.Token, err error) {
	return s.repo.List(offset, limit)
}
