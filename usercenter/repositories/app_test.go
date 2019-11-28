package repositories

import (
	"fmt"
	"log"
	"testing"

	"github.com/codelieche/microservice/usercenter/datamodels"

	"github.com/codelieche/microservice/usercenter/datasources"
)

func TestApplicationRepository_Save(t *testing.T) {
	// 1. get db
	db := datasources.GetDb()

	// 2. init app repositorry
	r := NewApplicationRepository(db)

	// 3. 测试插入分组
	i := 0
	for i < 10 {
		i++
		name := fmt.Sprintf("Application:%d", i)
		code := fmt.Sprintf("app_code_%d", i)
		app := &datamodels.Application{
			Name: name,
			Code: code,
		}
		if g, err := r.Save(app); err != nil {
			log.Println(err.Error())
		} else {
			log.Println(g.ID, g.Name)
		}
	}
}
func TestApplicationRepository_List(t *testing.T) {
	// 1. get db
	db := datasources.GetDb()

	// 2. init app repository
	r := NewApplicationRepository(db)

	// 3. 测试获取分组列表
	haveNext := true
	offset := 0
	limit := 5
	for haveNext {
		if apps, err := r.List(offset, limit); err != nil {
			t.Error(err.Error())
		} else {
			// 判断是否还有下一页
			if len(apps) == limit && limit > 0 {
				haveNext = true
				offset += limit
			} else {
				haveNext = false
			}
			// 输出分组
			for _, app := range apps {
				log.Println(app.ID, app.Name, app.Code)
			}
		}
	}
}

func TestAppRepository_GetPermissionList(t *testing.T) {
	// 1. get db
	db := datasources.GetDb()

	// 2. init app repository
	r := NewApplicationRepository(db)

	// 3. 获取app的权限
	appId := 1
	if app, err := r.Get(int64(appId)); err != nil {
		t.Error(err.Error())
		return
	} else {
		// 获取app的权限
		offset := 0
		limit := 5
		haveNext := true
		for haveNext {
			if permissions, err := r.GetPermissionList(app, offset, limit); err != nil {
				haveNext = false
				t.Error(err.Error())
			} else {
				// 判断是否有下一页
				if len(permissions) == limit && limit > 0 {
					haveNext = true
					offset += limit

				} else {
					haveNext = false
				}

				// 4. 打印出权限
				for _, permission := range permissions {
					log.Println(app.Name, app.Code, permission.Code, permission.Name)
				}
			}
		}
	}
}
