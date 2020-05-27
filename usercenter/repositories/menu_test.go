package repositories

import (
	"fmt"
	"log"
	"testing"

	"github.com/codelieche/microservice/usercenter/datamodels"
	"github.com/codelieche/microservice/usercenter/datasources"
)

func TestMenuRepository_Save(t *testing.T) {
	// 1. get db
	db := datasources.GetDb()

	// 2. init menu repository
	r := NewMenuRepository(db)

	// 3. 测试插入菜单
	i := 0
	for i < 10 {
		i++
		title := fmt.Sprintf("菜单%d", i)
		slug := fmt.Sprintf("/menu%d", i)
		menu := &datamodels.Menu{
			Title:   title,
			Slug:    slug,
			Project: "devops",
		}
		if menu, err := r.Save(menu); err != nil {
			t.Error(err)
		} else {
			log.Print(menu.ID, menu.Title, menu.Slug, menu.Level, menu.Permission)
		}
	}
}

func TestMenuRepository_Get(t *testing.T) {
	// 1. get db
	db := datasources.GetDb()

	// 2. init menu repository
	r := NewMenuRepository(db)
	// 3. 测试获取菜单
	i := 0
	for i < 1 {
		i++
		if menu, err := r.Get(int64(i)); err != nil {
			t.Error(err)
		} else {
			t.Log(menu.ID, menu.Title, menu.Slug)
			t.Log(menu.ParentID, menu.Children)
			if menu.Permission != nil {
				t.Log(menu.Permission.ID, menu.Permission.Code, menu.Permission.Project)
			}
		}
	}
}

func TestMenuRepository_List(t *testing.T) {
	// 1. get db
	db := datasources.GetDb()

	// 2. init menu repository
	r := NewMenuRepository(db)
	db.LogMode(true)

	// 3. 测试获取菜单列表
	haveNext := true
	offset := 0
	limit := 5
	for haveNext {
		if menus, count, err := r.List(offset, limit); err != nil {
			t.Error(err)
		} else {
			// 判断是否还有下一页
			if len(menus) == limit && limit > 0 {
				haveNext = true
				offset += limit
			} else {
				haveNext = false
			}

			log.Println("总共有菜单条数为：", count)

			// 输出菜单
			for _, menu := range menus {
				t.Log(menu.ID, menu.Title, menu.Slug)
				t.Log(menu.ParentID, menu.Children)
				children := menu.Children

				if len(children) > 0 {
					for _, m := range children {
						subChildren := m.Children
						t.Log("Children:", m.Title, m.Children, m.Permission)
						for _, m2 := range subChildren {
							t.Log("Children Children:", m2.Title, m2.Children, m2.Permission)
						}
					}
				}

				if menu.Permission != nil {
					t.Log(menu.Permission.ID, menu.Permission.Code, menu.Permission.Project)
				}
			}
		}
	}
}
