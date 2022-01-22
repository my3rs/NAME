package api

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/kataras/iris/v12/mvc"
)

type ResponseJSON struct {
	Code    int
	Message string
	Path    string
}

func Init(app *iris.Application) {
	apiAuth := app.Party("/api/v1/auth")
	mvc.Configure(apiAuth, func(mvcApp *mvc.Application) {
		mvcApp.Handle(new(Controller))
	})
	Verifier.Extractors = []jwt.TokenExtractor{jwt.FromHeader}
	//Verifier.Extractors = append(Verifier.Extractors, func(ctx iris.Context) string {
	//	return ctx.GetCookie("Authorization")
	//})

	info := app.Party("/api/v1/info")
	info.Use(VerifierMiddleware)
	mvc.Configure(info, func(application *mvc.Application) {
		application.Handle(new(InfoController))
	})

}
