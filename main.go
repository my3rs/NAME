package main

import (
	"NAME/conf"
	"NAME/route"
	"flag"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/cors"
)

var configFilePath = flag.String("c", "./name.conf", "配置文件路径")

func main() {
	flag.Parse()

	app := iris.New()
	//app.Logger().SetLevel("debug")
	app.Logger().Info("Config file: ", *configFilePath)

	conf.SetConfigPath(*configFilePath)

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

	tmpl := iris.HTML("./web/view", ".gohtml")
	tmpl.Reload(true)
	app.RegisterView(tmpl)
	app.HandleDir("/assets", "./web/assets")
	app.HandleDir("/uploads", conf.Config().DataPath+"/uploads")

	// Listens and serves incoming http requests
	// on http://localhost:8000.
	app.Listen(":" + strconv.Itoa(conf.Config().Port))
}

func logMiddleware(ctx iris.Context) {
	ctx.Application().Logger().Infof("%s %s", ctx.Request().Method, ctx.Path())
	ctx.Next()
}
