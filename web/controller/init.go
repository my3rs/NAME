package controller

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func Init(app *iris.Application) {
	mvc.Configure(app.Party("/post"), func(mvcApp *mvc.Application) {
		mvcApp.Handle(new(PostController))
	})

	mvc.Configure(app.Party("/"), func(mvcApp *mvc.Application) {
		mvcApp.Handle(new(IndexController))
	})

	dashboard := app.Party("/admin")
	//dashboard.Use(api.VerifierMiddleware)
	mvc.Configure(dashboard, func(application *mvc.Application) {
		application.Handle(new(DashboardController))
	})

	mvc.Configure(app.Party("/user"), func(mvcApp *mvc.Application) {
		mvcApp.Handle(new(UserController))
	})
}
