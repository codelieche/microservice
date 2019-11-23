package repositories

import (
	"fmt"
	"log"
	"testing"

	"github.com/codelieche/microservice/datamodels"

	"github.com/codelieche/microservice/datasources"
)

func TestPermissionRepository_Save(t *testing.T) {
	// 1. get db
	db := datasources.GetDb()

	// 2. init permission repositorry
	r := NewPermissionRepository(db)

	// 3. 测试插入分组
	i := 0
	for i < 10 {
		i++
		name := fmt.Sprintf("Permission:%d", i)
		code := fmt.Sprintf("code_%d", i)
		permission := &datamodels.Permission{
			Name:  name,
			AppID: 1,
			Code:  code,
		}
		if g, err := r.Save(permission); err != nil {
			log.Println(err.Error())
		} else {
			log.Println(g.ID, g.Name)
		}
	}
}
func TestPermissionRepository_List(t *testing.T) {
	// 1. get db
	db := datasources.GetDb()

	// 2. init permission repository
	r := NewPermissionRepository(db)

	// 3. 测试获取分组列表
	haveNext := true
	offset := 0
	limit := 5
	for haveNext {
		if permissions, err := r.List(offset, limit); err != nil {
			t.Error(err.Error())
		} else {
			// 判断是否还有下一页
			if len(permissions) == limit && limit > 0 {
				haveNext = true
				offset += limit
			} else {
				haveNext = false
			}
			// 输出分组
			for _, permission := range permissions {
				log.Println(permission.ID, permission.Name)
			}
		}
	}
}

func TestPermissionRepository_GetByAppIDAndCode(t *testing.T) {
	// 1. get db
	db := datasources.GetDb()

	// 2. init permission repository
	r := NewPermissionRepository(db)

	// 3. get permission
	appID := 1
	code := "code_2"
	if permission, err := r.GetByAppIDAndCode(appID, code); err != nil {
		t.Error(err.Error())
	} else {
		// 4. print permission info
		log.Println(permission.Name, permission.Code, permission.AppID, permission.Application)
	}
}
