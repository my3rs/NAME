package route

import (
	"NAME/conf"
	"NAME/controller"
	"NAME/database"
	"NAME/middleware"
	"NAME/model"
	"NAME/service"

	//web "NAME/web/controller"
	"os"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func InitRoute(app *iris.Application) {
	// create directory to store uploaded files
	os.Mkdir(conf.Config().DataPath, 0700)
	err := os.Mkdir(conf.Config().DataPath+"/uploads", 0700)
	if err != nil {
		app.Logger().Info("Failed to create data folder: ", err)
	}

	// 认证
	apiAuth := app.Party("/api/v1/auth")
	mvc.Configure(apiAuth, func(app *mvc.Application) {
		app.Register(
			database.GetDB,
			service.NewUserService,
		)
		app.Handle(new(controller.AuthController))
	})

	// 附件
	attachments := app.Party("/api/v1/attachments")
	if conf.Config().Mode == conf.PROD {
		attachments.Use(middleware.JwtMiddleware())
	}
	mvc.Configure(attachments, func(app *mvc.Application) {
		app.Register(
			database.GetDB,
			service.NewAttachmentService,
		)
		app.Handle(new(controller.AttachmentController))
	})

	// Everyone can GET from /api/API_VERSION/posts
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
	}
	mvc.Configure(posts, func(application *mvc.Application) {
		application.Register(
			database.GetDB,
			service.NewContentService,
			service.NewTagService,
		)
		application.Handle(new(controller.PostController))
	})

	// 评论
	comments := app.Party("/api/v1/comments")
	//if conf.Config().Mode == conf.PROD {
	//	comments.Use(jwtMiddleware)
	//}
	comments.Use(jwtMiddleware)
	mvc.Configure(comments, func(application *mvc.Application) {
		application.Register(
			database.GetDB,
			service.NewCommentService,
		)
		application.Handle(new(controller.CommentController))
	})

	pages := app.Party("/api/v1/pages")
	if conf.Config().Mode == conf.PROD {
		pages.Use(middleware.JwtMiddleware())
	}
	mvc.Configure(pages, func(application *mvc.Application) {
		application.Handle(new(controller.PageController))
	})

	// 用户
	users := app.Party("/api/v1/users")
	if conf.Config().Mode == conf.PROD {
		users.Use(middleware.JwtMiddleware())

	}
	mvc.Configure(users, func(application *mvc.Application) {
		application.Register(
			database.GetDB,
			service.NewUserService,
		)
		application.Handle(new(controller.UserController))
	})

	// 设置
	settings := app.Party("/api/v1/settings")
	if env, found := model.GetSettingsItem("environment"); found && env.Value == model.EnvironmentProd {
		settings.Use(middleware.JwtMiddleware())
	}
	mvc.Configure(settings, func(app *mvc.Application) {
		app.Register(
			database.GetDB,
		)
		app.Handle(new(controller.SettingController))
	})

	// 标签
	tags := app.Party("/api/v1/tags")
	if conf.Config().Mode == conf.PROD {
		tags.Use(jwtMiddleware)
	}
	mvc.Configure(tags, func(app *mvc.Application) {
		app.Register(
			database.GetDB,
			service.NewTagService,
		)

		app.Handle(new(controller.TagController))
	})
}
