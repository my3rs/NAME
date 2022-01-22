package main

import (
	"change/api"
	"change/model"
	"change/web/controller"
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()
	app.Use(logMiddleware)

	// 初始化一些设置
	model.Init()
	controller.Init(app)
	api.Init(app)

	app.RegisterView(iris.HTML("./web/view", ".html"))
	app.HandleDir("/public", "./web/public")

	// Listens and serves incoming http requests
	// on http://localhost:8080.
	app.Listen(":8080")
}

func myMiddleware() iris.Handler {
	return func(ctx iris.Context) {
		ctx.Application().Logger().Infof("Runs before %s, cookie: %s, token: %s", ctx.Path(), ctx.GetCookie("Authorization"), ctx.GetHeader("Authorization"))
		ctx.Next()
	}
}

func logMiddleware(ctx iris.Context) {
	ctx.Application().Logger().Infof("Runs before %s, cookie: %s, token: %s", ctx.Path(), ctx.GetCookie("Authorization"), ctx.GetHeader("Authorization"))
	ctx.Next()
}
