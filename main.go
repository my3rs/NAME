package main

import (
	"NAME/conf"
	_ "NAME/docs"
	"NAME/route"
	"github.com/iris-contrib/swagger/v12"
	"github.com/iris-contrib/swagger/v12/swaggerFiles"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/cors"
	"strconv"
)

// @title NAME API
// @version 1.0
// @description This is a distributed blog server.

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email me@seahi.me

// @host localhost:8000
func main() {
	app := iris.New()
	swaggerConfig := &swagger.Config{
		URL:         "http://localhost:8000/swagger/doc.json",
		DeepLinking: true,
	}
	app.Get("/swagger/{any:path}", swagger.CustomWrapHandler(swaggerConfig, swaggerFiles.Handler))

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
	ctx.Application().Logger().Infof("Runs before %s %s", ctx.Request().Method, ctx.Path())
	ctx.Next()
}
