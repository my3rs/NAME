package main

import (
	"NAME/conf"
	"NAME/route"
	"log"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/cors"
)

func init() {
	log.Println("初始化，Package main")

}

func main() {

	app := iris.New()
	//app.Logger().SetLevel("debug")

	// Register CORS (allow any origin to pass through) middleware.
	app.UseRouter(cors.New().
		ExtractOriginFunc(cors.DefaultOriginExtractor).
		ReferrerPolicy(cors.NoReferrerWhenDowngrade).
		AllowOriginFunc(cors.AllowAnyOrigin).
		Handler())

	app.Use(logMiddleware)
	//app.OnErrorCode(iris.StatusUnauthorized, handleUnauthorized)

	// Init some configures
	route.InitRoute(app)

	app.HandleDir("/assets", conf.GetConfig().AssetsPath)
	app.HandleDir("/uploads", conf.GetConfig().DataPath+"/uploads")

	// Listens and serves incoming http requests
	// on http://localhost:8000.
	app.Listen(conf.GetConfig().Host + ":" + strconv.Itoa(conf.GetConfig().Port))
}

func logMiddleware(ctx iris.Context) {
	ctx.Application().Logger().Infof("%s %s", ctx.Request().Method, ctx.Path())
	ctx.Next()
}
