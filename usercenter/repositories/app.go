package repositories

import (
	"github.com/codelieche/microservice/common"
	"github.com/codelieche/microservice/datamodels"
	"github.com/jinzhu/gorm"
)

type ApplicationRepository interface {
	// 保存App
	Save(app *datamodels.Application) (*datamodels.Application, error)
	// 获取App的列表
	List(offset int, limit int) ([]*datamodels.Application, error)
	// 获取App信息
	Get(id int64) (*datamodels.Application, error)
	// 根据ID或者Code获取App信息
	GetByIdOrCode(idOrCode string) (app *datamodels.Application, err error)
	// 获取App的Permission列表
	GetPermissionList(app *datamodels.Application, offset int, limit int) ([]*datamodels.Permission, error)
}

func NewApplicationRepository(db *gorm.DB) ApplicationRepository {
	return &appRepository{db: db}
}

type appRepository struct {
	db *gorm.DB
}

// 保存Application
func (r *appRepository) Save(app *datamodels.Application) (*datamodels.Application, error) {
	if app.ID > 0 {
		// 是更新操作
		if err := r.db.Model(&datamodels.Application{}).Update(app).Error; err != nil {
			return nil, err
		} else {
			return app, nil
		}
	} else {
		// 是创建操作
		// 随机生成个Token
		token := common.RandString(32)
		app.Token = token
		if err := r.db.Create(app).Error; err != nil {
			return nil, err
		} else {
			return app, nil
		}

	}
}

// 获取App的列表
func (r *appRepository) List(offset int, limit int) (apps []*datamodels.Application, err error) {
	query := r.db.Model(&datamodels.Application{}).Offset(offset).Limit(limit).Find(&apps)
	if query.Error != nil {
		return nil, query.Error
	} else {
		return apps, nil
	}
	return
}

// 根据ID获取App
func (r *appRepository) Get(id int64) (app *datamodels.Application, err error) {

	app = &datamodels.Application{}
	r.db.First(app, "id = ?", id)
	if app.ID > 0 {
		return app, nil
	} else {
		return nil, common.NotFountError
	}
}

// 根据ID/Code获取App
func (r *appRepository) GetByIdOrCode(idOrCode string) (app *datamodels.Application, err error) {

	app = &datamodels.Application{}
	r.db.First(app, "id = ? or code = ?", idOrCode, idOrCode)
	if app.ID > 0 {
		return app, nil
	} else {
		return nil, common.NotFountError
	}
}

// 获取App的用户
func (r *appRepository) GetPermissionList(
	app *datamodels.Application, offset int, limit int) (permissions []*datamodels.Permission, err error) {
	query := r.db.Model(app).Offset(offset).Limit(limit).Related(&permissions, "Permissions")
	if query.Error != nil {
		return nil, query.Error
	} else {
		return permissions, nil
	}
}
