package main

import (
	"NAME/conf"
	_ "NAME/docs"
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

	app.RegisterView(iris.HTML("./web/view", ".html"))
	app.HandleDir("/public", "./web/public")

	// Listens and serves incoming http requests
	// on http://localhost:8000.
	app.Listen(":" + strconv.Itoa(conf.Config().Port))
}

func logMiddleware(ctx iris.Context) {
	ctx.Application().Logger().Infof("%s %s", ctx.Request().Method, ctx.Path())
	ctx.Next()
}
