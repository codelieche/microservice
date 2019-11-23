package repositories

import (
	"fmt"
	"log"
	"testing"

	"github.com/codelieche/microservice/datamodels"
	"github.com/codelieche/microservice/datasources"
)

func TestRoleRepository_Save(t *testing.T) {
	// 1. get db
	db := datasources.GetDb()

	// 2. init role repository
	r := NewRoleRepository(db)

	// 3. 测试插入Role
	i := 0
	for i < 10 {
		i++
		name := fmt.Sprintf("Role:%d", i)
		role := &datamodels.Role{
			Name:        name,
			Permissions: nil,
		}
		if g, err := r.Save(role); err != nil {
			log.Println(err.Error())
		} else {
			log.Println(g.ID, g.Name)
		}
	}
}

func TestRoleRepository_List(t *testing.T) {
	// 1. get db
	db := datasources.GetDb()

	// 2. init role repository
	r := NewRoleRepository(db)

	// 3. 测试获取Role列表
	haveNext := true
	offset := 0
	limit := 5
	for haveNext {
		if roles, err := r.List(offset, limit); err != nil {
			t.Error(err.Error())
			haveNext = false
		} else {
			if len(roles) == limit && limit > 0 {
				haveNext = true
				offset += limit
			} else {
				haveNext = false
			}
			// 输出分组
			for _, role := range roles {
				log.Println(role.ID, role.Name)
			}
		}
	}

}
