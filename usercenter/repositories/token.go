package repositories

import (
	"errors"

	"github.com/codelieche/microservice/usercenter/common"

	"github.com/codelieche/microservice/usercenter/datamodels"
	"github.com/jinzhu/gorm"
)

type TokenRepository interface {
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

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{
		fields: []string{"id", "created_at", "user_id", "token", "is_active", "expired_at"},
		db:     db,
	}
}

type tokenRepository struct {
	fields []string
	db     *gorm.DB
}

func (r *tokenRepository) Create(token *datamodels.Token) (*datamodels.Token, error) {
	if token.ID > 0 {
		err := errors.New("不可自行设置ID")
		return nil, err
	}
	// 校验token
	if token.Token == "" {
		i := 0
		for i < 50 {
			tokenStr := common.RandString(32)
			if _, err := r.GetByToken(tokenStr); err == common.NotFountError {
				token.Token = tokenStr
				break
			} else {
				// 当前这个token被占用了，重新生成一个
			}
			i++
		}
		if token.Token == "" {
			err := errors.New("随机自动生成token字符串出错")
			return nil, err
		}
	}
	// 开始创建
	if err := r.db.Create(token).Error; err != nil {
		return nil, err
	} else {
		return token, nil
	}
}

func (r *tokenRepository) Get(id int64) (token *datamodels.Token, err error) {
	token = &datamodels.Token{}
	r.db.Select(r.fields).First(token, "id = ?", id)
	if token.ID > 0 {
		return token, nil
	} else {
		return nil, common.NotFountError
	}
}

func (r *tokenRepository) GetByToken(token string) (*datamodels.Token, error) {
	tokenObj := &datamodels.Token{}
	r.db.Select(r.fields).First(tokenObj, "token = ?", token)
	if tokenObj.ID > 0 {
		return tokenObj, nil
	} else {
		return nil, common.NotFountError
	}
}

func (r *tokenRepository) GetByIdOrToken(idOrToken string) (token *datamodels.Token, err error) {
	token = &datamodels.Token{}
	r.db.Select(r.fields).First(token, "id = ? or token = ?", idOrToken, idOrToken)
	if token.ID > 0 {
		return token, nil
	} else {
		return nil, common.NotFountError
	}
}

func (r *tokenRepository) Delete(token *datamodels.Token) (err error) {
	// 设置IsActive为false
	if token.IsActive {
		updateFields := map[string]interface{}{"is_active": false}
		if _, err = r.UpdateByID(int64(token.ID), updateFields); err != nil {
			return err
		} else {
			return nil
		}
	} else {
		return nil
	}
}

func (r *tokenRepository) Update(token *datamodels.Token, fields map[string]interface{}) (*datamodels.Token, error) {
	return r.UpdateByID(int64(token.ID), fields)
}

func (r *tokenRepository) UpdateByID(id int64, fields map[string]interface{}) (token *datamodels.Token, err error) {
	// 判断ID
	if id <= 0 {
		err = errors.New("传入ID为0,会更新全部数据")
		return nil, err
	}
	// 因为指定了ID了，所以这里可不判断这个ID
	// 丢弃ID/id/Id/iD
	//idKeys := []string{"ID", "id", "Id", "iD"}
	//for _, k := range idKeys {
	//	if _, exist := fields[k]; exist {
	//		delete(fields, k)
	//	}
	//}

	// 更新操作
	if err = r.db.Model(&datamodels.Token{}).Where("id = ?", id).Limit(1).Update(fields).Error; err != nil {
		return nil, err
	} else {
		return r.Get(id)
	}
}

func (r *tokenRepository) List(offset int, limit int) (tokens []*datamodels.Token, err error) {
	query := r.db.Model(&datamodels.Token{}).Select(r.fields).Offset(offset).Limit(limit).Find(&tokens)
	if query.Error != nil {
		return nil, query.Error
	} else {
		return tokens, nil
	}
}
