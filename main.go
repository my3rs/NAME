package main

import (
	"NAME/conf"
	"NAME/route"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/cors"
)

func main() {
	app := iris.New()

	conf.Init()

	// Register CORS (allow any origin to pass through) middleware.
	app.UseRouter(cors.New().
		ExtractOriginFunc(cors.DefaultOriginExtractor).
		ReferrerPolicy(cors.NoReferrerWhenDowngrade).
		AllowOriginFunc(cors.AllowAnyOrigin).
		Handler())

	app.Use(logMiddleware)
	//app.OnErrorCode(iris.StatusUnauthorized, handleUnauthorized)

	// 初始化一些设置
	route.InitRoute(app)

	tmpl := iris.HTML("./web/view", ".tmpl")
	tmpl.Reload(true)
	app.RegisterView(tmpl)
	app.HandleDir("/assets", "./web/assets")

	// Listens and serves incoming http requests
	// on http://localhost:8000.
	app.Listen(":" + strconv.Itoa(conf.Config().Port))
}

func logMiddleware(ctx iris.Context) {
	ctx.Application().Logger().Infof("%s %s", ctx.Request().Method, ctx.Path())
	ctx.Next()
}
