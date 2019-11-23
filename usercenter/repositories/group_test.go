package repositories

import (
	"fmt"
	"log"
	"testing"

	"github.com/codelieche/microservice/datamodels"

	"github.com/codelieche/microservice/datasources"
)

func TestGroupRepository_Save(t *testing.T) {
	// 1. get db
	db := datasources.GetDb()

	// 2. init group repositorry
	r := NewGroupRepository(db)

	// 3. 测试插入分组
	i := 0
	for i < 10 {
		i++
		name := fmt.Sprintf("Group:%d", i)
		group := &datamodels.Group{
			Name:        name,
			Permissions: nil,
		}
		if g, err := r.Save(group); err != nil {
			log.Println(err.Error())
		} else {
			log.Println(g.ID, g.Name)
		}
	}
}
func TestGroupRepository_List(t *testing.T) {
	// 1. get db
	db := datasources.GetDb()

	// 2. init group repository
	groupRepository := NewGroupRepository(db)

	// 3. 测试获取分组列表
	haveNext := true
	offset := 0
	limit := 5
	for haveNext {
		if groups, err := groupRepository.List(offset, limit); err != nil {
			t.Error(err.Error())
		} else {
			// 判断是否还有下一页
			if len(groups) == limit && limit > 0 {
				haveNext = true
				offset += limit
			} else {
				haveNext = false
			}
			// 输出分组
			for _, group := range groups {
				log.Println(group.ID, group.Name)
			}
		}
	}

}
