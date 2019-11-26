package controllers

import (
	"fmt"
	"time"

	"github.com/kataras/iris/v12/sessions"

	"github.com/kataras/iris/v12"
)

type IndexController struct {
	Ctx       iris.Context
	StartTime time.Time
}

func (c *IndexController) Get(ctx iris.Context) {
	// 通过session获取visits
	sess := sessions.Get(ctx)
	visits := sess.Increment("visits", 1)
	//visits := sess.GetIntDefault("visits", 1)
	sinces := time.Now().Sub(c.StartTime).Seconds()

	//log.Println(ctx.Path())
	//username, password, _ := ctx.Request().BasicAuth()
	//log.Println(ctx.Path(), username, password)
	//ctx.Writef("%s %s %s", ctx.Path(), username, password)
	//msg := fmt.Sprintf("%s session id: %s", ctx.Path(), sess.ID())

	ctx.JSON(iris.Map{
		"path":    ctx.Path(),
		"session": sess.ID(),
		"time":    time.Now(),
		"other":   "Index Controller",
		"visits":  visits,
		"sinces":  sinces,
	})
}

func (c *IndexController) GetPing(ctx iris.Context) {

	session := sessions.Get(ctx)
	session.Set("ping", "pong")

	result := session.Get("ping")

	if result != nil {

	} else {
		// session出问题了
		ctx.StatusCode(500)
	}

	ctx.JSON(
		iris.Map{
			"session": session.ID(),
			"message": result,
		})

}

// 用户登录页面
func (c *IndexController) GetUserLogin() {
	session := sessions.Get(c.Ctx)
	userID := session.GetIntDefault("userID", 0)
	if userID > 0 {
		msg := fmt.Sprintf("已经登录，用户名id是：%d", userID)
		c.Ctx.ViewData("msg", msg)
	} else {

	}

	//c.Ctx.ViewData("msg", "提示消息")
	c.Ctx.View("user/login.html")
}
