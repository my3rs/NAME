package route

import (
	"NAME/conf"
	"NAME/controller"
	"NAME/middleware"
	web "NAME/web/controller"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"os"
)

func InitRoute(app *iris.Application) {
	// create directory to store uploaded files
	os.Mkdir(conf.Config().DataPath, 0700)
	err := os.Mkdir(conf.Config().DataPath+"/uploads", 0700)
	if err != nil {
		app.Logger().Error("Failed to create data folder: ", err)
	}

	// backend
	apiAuth := app.Party("/api/v1/auth")
	mvc.Configure(apiAuth, func(mvcApp *mvc.Application) {
		mvcApp.Handle(controller.NewAuthController())
	})

	attachments := app.Party("/api/v1/attachments")
	if conf.Config().Mode == conf.PROD {
		attachments.Use(middleware.JwtMiddleware())
	}
	mvc.Configure(attachments, func(application *mvc.Application) {
		application.Handle(controller.NewAttachmentController())
	})

	jwtFilter := func(ctx iris.Context) bool {
		if method := ctx.Method(); method == iris.MethodGet {
			ctx.Next()
		}
		return true
	}
	jwtMiddleware := iris.NewConditionalHandler(jwtFilter, middleware.JwtMiddleware())

	posts := app.Party("/api/v1/posts")
	if conf.Config().Mode == conf.PROD {
		posts.Use(jwtMiddleware)
		//posts.Get("/").RemoveHandler()
		//posts.Get("/").RemoveHandler("middleware.JwtMiddleware.func2")
		//for _, h := range posts.Get("/", func(context2 iris.Context) {}).Handlers {
		//	log.Print("xxx ", context.HandlerName(h))
		//}
		//log.Printf("xxxx %+v", posts.Get("/", func(context2 iris.Context) {}).Handlers)
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

	mvc.Configure(app.Party("/posts"), func(mvcApp *mvc.Application) {
		mvcApp.Handle(web.NewPostController())
	})

}
