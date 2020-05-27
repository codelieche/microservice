package repositories

import (
	"fmt"
	"log"
	"testing"

	"github.com/codelieche/microservice/usercenter/datamodels"
	"github.com/codelieche/microservice/usercenter/datasources"
)

func TestProjectRepository_Save(t *testing.T) {
	// 1. get db
	db := datasources.GetDb()

	// 2. init project repository
	r := NewProjectRepository(db)

	// 测试插入数据
	i := 0
	for i < 10 {
		i++
		name := fmt.Sprintf("Project:%d", i)
		code := fmt.Sprintf("project_code_%d", i)
		project := &datamodels.Project{
			Name: name,
			Code: code,
		}

		if p, err := r.Save(project); err != nil {
			log.Println(err.Error())
		} else {
			log.Println(p.ID, p.Name, p.Code, p.Description)
		}
	}
}

func TestProjectRepository_List(t *testing.T) {
	// 1. get db
	db := datasources.GetDb()

	// 2. init project repository
	r := NewProjectRepository(db)

	// 测试胡哦哦去列表数据
	haveNext := true
	offset := 0
	limit := 5
	//db.LogMode(true)
	for haveNext {
		if projects, count, err := r.List(offset, limit); err != nil {
			t.Error(err.Error())
		} else {
			log.Println("count = ", count)
			// 判断是否还有下一页
			if len(projects) == limit && limit > 0 {
				haveNext = true
				offset += limit
			} else {
				//log.Println(projects)
				log.Println("没有下一页了")
				haveNext = false
			}

			// 输出Project
			for _, project := range projects {
				log.Println(project.ID, project.Name, project.Code, project.Description, project.Users)
			}
		}
	}
}

func TestProjectRepository_Get(t *testing.T) {
	// 1. get db
	db := datasources.GetDb()

	// 2. init project repository
	r := NewProjectRepository(db)

	// 测试插入数据
	i := 0
	for i < 1 {
		i++
		if p, err := r.Get(int64(i)); err != nil {
			log.Println(err.Error())
		} else {
			log.Println(p.ID, p.Name, p.Code, p.Description, p.Users)
		}
	}
}

func TestProjectRepository_Update(t *testing.T) {
	// 1. get db
	db := datasources.GetDb()

	// 2. init project repository
	r := NewProjectRepository(db)

	// 测试插入数据
	i := 0
	for i < 1 {
		i++
		if p, err := r.Get(int64(i)); err != nil {
			log.Println(err.Error())
		} else {
			log.Println(p.ID, p.Name, p.Code, p.Description)
			p.Users = []*datamodels.ProjectUser{
				&datamodels.ProjectUser{
					Role:    "admin",
					Project: "devops",
					User:    "admin",
				},
			}

			log.Println(r.Save(p))
		}
	}
}

func TestProjectRepository_GetPermissionList(t *testing.T) {
	// 1. get db
	db := datasources.GetDb()

	// 2. init project repository
	r := NewProjectRepository(db)
	// 测试插入数据
	i := 0
	for i < 1 {
		i++
		if p, err := r.GetByIdOrCode("devops"); err != nil {
			log.Println(err.Error())
		} else {
			log.Println(p.ID, p.Name, p.Code, p.Description, p.Users)

			// 获取所有的权限
			if permissions, err := r.GetPermissionList(p, 0, 100); err != nil {
				log.Println(err.Error())
			} else {
				for _, permission := range permissions {
					log.Println(permission.ID, permission.Name, permission.Project, permission.Code)
				}
			}
		}
	}
}

func TestProjectRepository_AddProjectUser(t *testing.T) {
	// 1. get db
	db := datasources.GetDb()

	// 2. init project repository
	r := NewProjectRepository(db)

	// 3. 测试添加项目成员
	if project, err := r.GetByIdOrCode("devops"); err != nil {
		t.Error(err)
	} else {
		// 插入用户
		i := 0
		for i < 5 {
			i++
			username := fmt.Sprintf("user%d", i)
			role := "develop"
			if i%2 == 0 {
				role = "test"
			}
			if pUser, err := r.AddProjectUser(project.Code, username, role); err != nil {
				t.Error(err)
			} else {
				log.Println(pUser)
			}
		}
	}
}

func TestProjectRepository_DeleteProjectUser(t *testing.T) {
	// 1. get db
	db := datasources.GetDb()

	// 2. init project repository
	r := NewProjectRepository(db)

	// 3. 测试删除用户
	if project, err := r.GetByIdOrCode("devops"); err != nil {
		t.Error(err)
	} else {
		// 删除用户
		i := 0
		for i < 5 {
			i++
			username := fmt.Sprintf("user%d", i)
			if success, err := r.DeleteProjectUser(project.Code, username); err != nil {
				t.Error(err)
			} else {
				if success {
					log.Println(project.Code, "删除用户成功", username)
				} else {
					log.Println(project.Code, "删除用户没有成功！！！", username)
				}
			}
		}
	}
}
