package repositories

import (
	"errors"

	"github.com/codelieche/microservice/usercenter/common"

	"github.com/codelieche/microservice/usercenter/datamodels"
	"github.com/jinzhu/gorm"
)

type MenuRepository interface {
	// 保存
	Save(menu *datamodels.Menu) (*datamodels.Menu, error)
	// 获取
	Get(id int64) (menu *datamodels.Menu, err error)
	GetByProjectAndSlug(project string, slug string) (menu *datamodels.Menu, err error)

	// 列表
	List(offset int, limit int) (menus []*datamodels.Menu, count int, err error)

	// 更新
	Update(menu *datamodels.Menu, updateFields map[string]interface{}) (*datamodels.Menu, error)
	UpdateByID(id int64, updateFields map[string]interface{}) (*datamodels.Menu, error)

	// 删除
	Delete(menu *datamodels.Menu) (success bool, err error)
}

func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{
		db: db,
		infoFields: []string{"id", "project", "title", "slug", "icon",
			"parent_id", "level", "permission_id",
			"target", "is_link", "link",
			"`order`",
			"is_deleted",
		},
		permissionFields: []string{"id", "name", "project", "code"},
	}
}

type menuRepository struct {
	db               *gorm.DB
	infoFields       []string
	permissionFields []string
}

func (r *menuRepository) Save(menu *datamodels.Menu) (*datamodels.Menu, error) {
	if menu.Title == "" {
		err := errors.New("title字段不可为空")
		return nil, err
	}
	if menu.Slug == "" {
		err := errors.New("slug字段不可为空")
		return nil, err
	}

	if menu.ID > 0 {
		// 是更新操作
		tx := r.db.Begin()
		// 修改menu的信息
		if err := tx.Model(&datamodels.Menu{}).Update(menu).Error; err != nil {
			tx.Rollback()
			return nil, err
		} else {
			// 提交更新事务
			tx.Commit()
			return r.Get(int64(menu.ID))
		}
	} else {
		// 是创建操作
		if err := r.db.Create(menu).Error; err != nil {
			return nil, err
		} else {
			return menu, nil
		}
	}
}

func (r *menuRepository) Get(id int64) (menu *datamodels.Menu, err error) {
	menu = &datamodels.Menu{}
	err = r.db.
		Preload("Permission", func(db *gorm.DB) *gorm.DB {
			return db.Select(r.permissionFields)
		}).
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Select(r.infoFields)
		}).
		Select(r.infoFields).
		First(menu, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	if menu.ID > 0 {
		return menu, nil
	} else {
		return nil, common.NotFountError
	}
}

func (r *menuRepository) GetByProjectAndSlug(project string, slug string) (menu *datamodels.Menu, err error) {
	menu = &datamodels.Menu{}
	query := r.db.Select(r.infoFields).
		Preload("Permission", func(db *gorm.DB) *gorm.DB {
			return r.db.Select(r.permissionFields)
		}).Preload("Children", func(db *gorm.DB) *gorm.DB {
		return r.db.Select(r.infoFields)
	}).First(menu, "project = ? and slug = ?", project, slug)

	if err = query.Error; err != nil {
		return nil, err
	} else {
		if menu.ID > 0 {
			return menu, nil
		} else {
			return nil, common.NotFountError
		}
	}
}

func (r *menuRepository) List(offset int, limit int) (menus []*datamodels.Menu, count int, err error) {
	query := r.db.Model(&datamodels.Menu{}).
		Select(r.infoFields).
		Preload("Permission", func(db *gorm.DB) *gorm.DB {
			return r.db.Select(r.permissionFields)
		}).Preload("Children", func(db *gorm.DB) *gorm.DB {
		// 继续Preload一层就可以查询到三层Children了
		return r.db.Select(r.infoFields).
			Preload("Permission", func(db *gorm.DB) *gorm.DB {
				return r.db.Select(r.permissionFields)
			}).
			Preload("Children", func(db *gorm.DB) *gorm.DB {
				return r.db.Select(r.infoFields).Preload("Permission", func(db *gorm.DB) *gorm.DB {
					return r.db.Select(r.permissionFields)
				})
			})
	}).Count(&count).Offset(offset).Limit(limit).
		Find(&menus)
	if err = query.Error; err != nil {
		return nil, 0, err
	} else {
		return menus, count, nil
	}
}

func (r *menuRepository) Update(menu *datamodels.Menu, updateFields map[string]interface{}) (*datamodels.Menu, error) {
	panic("implement me")
}

func (r *menuRepository) UpdateByID(id int64, updateFields map[string]interface{}) (*datamodels.Menu, error) {
	panic("implement me")
}

func (r *menuRepository) Delete(menu *datamodels.Menu) (success bool, err error) {
	panic("implement me")
}
