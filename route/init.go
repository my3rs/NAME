package route

import (
	"NAME/conf"
	"NAME/controller"
	"NAME/middleware"
	web "NAME/web/controller"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func InitRoute(app *iris.Application) {
	// backend
	apiAuth := app.Party("/api/v1/auth")
	mvc.Configure(apiAuth, func(mvcApp *mvc.Application) {
		mvcApp.Handle(controller.NewAuthController())
	})

	posts := app.Party("/api/v1/posts")
	if conf.Config().Mode == conf.PROD {
		posts.Use(middleware.JwtMiddleware())
	}
	mvc.Configure(posts, func(application *mvc.Application) {
		application.Handle(controller.NewPostController())
	})

	pages := app.Party("/api/v1/pages")
	if conf.Config().Mode == conf.PROD {
		pages.Use(middleware.JwtMiddleware())
	}
	mvc.Configure(pages, func(application *mvc.Application) {
		application.Handle(new(controller.PageController))
	})

	users := app.Party("/api/v1/users")
	if conf.Config().Mode == conf.PROD {
		users.Use(middleware.JwtMiddleware())
	}
	mvc.Configure(users, func(application *mvc.Application) {
		application.Handle(controller.NewUserController())
	})

	status := app.Party("/api/v1/status")
	if conf.Config().Mode == conf.PROD {
		status.Use(middleware.JwtMiddleware())
	}
	mvc.Configure(status, func(application *mvc.Application) {
		application.Handle(controller.NewStatusController())
	})

	tags := app.Party("/api/v1/tags")
	if conf.Config().Mode == conf.PROD {
		tags.Use(middleware.JwtMiddleware())
	}
	mvc.Configure(tags, func(application *mvc.Application) {
		application.Handle(controller.NewTagController())
	})

	// front
	mvc.Configure(app.Party("/"), func(mvcApp *mvc.Application) {
		mvcApp.Handle(web.NewIndexController())
	})

	mvc.Configure(app.Party("/post"), func(mvcApp *mvc.Application) {
		mvcApp.Handle(web.NewPostController())
	})

}
