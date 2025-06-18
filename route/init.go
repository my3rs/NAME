package route

import (
	"NAME/conf"
	"NAME/controller"
	"NAME/database"
	"NAME/middleware"
	"NAME/service"

	//web "NAME/web/controller"
	"os"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func InitRoute(app *iris.Application) {
	// 创建保存附件的目录
	os.Mkdir(conf.GetConfig().DataPath, 0700)
	err := os.Mkdir(conf.GetConfig().DataPath+"/uploads", 0700)
	if err != nil {
		app.Logger().Info("Failed to create data folder: ", err)
	}

	// 博客状态
	meta := app.Party("/api/v1/meta")
	mvc.Configure(meta, func(app *mvc.Application) {
		app.Register(
			database.GetDB(),
			service.NewContentService(),
		)
		app.Handle(new(controller.MetaController))
	})

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
	attachments.Use(middleware.JWT)

	mvc.Configure(attachments, func(app *mvc.Application) {
		app.Register(
			database.GetDB,
			service.NewAttachmentService,
		)
		app.Handle(new(controller.AttachmentController))
	})

	// 文章

	// 所有人都可以向 /api/API_VERSION/posts 发送GET请求
	jwtFilter := func(ctx iris.Context) bool {
		if method := ctx.Method(); method == iris.MethodGet {
			return false
		}
		return true
	}
	jwtMiddleware := iris.NewConditionalHandler(jwtFilter, middleware.JWT)

	posts := app.Party("/api/v1/posts", jwtMiddleware)

	mvc.Configure(posts, func(application *mvc.Application) {
		application.Register(
			database.GetDB,
			service.NewContentService,
			service.NewTagService,
		)
		application.Handle(new(controller.PostController))
	})

	// 文章分类
	categories := app.Party("/api/v1/categories")
	mvc.Configure(categories, func(app *mvc.Application) {
		app.Register(
			database.GetDB(),
			service.NewCategoryService(),
		)
		app.Handle(new(controller.CategoryController))
	})

	// 评论
	comments := app.Party("/api/v1/comments")
	mvc.Configure(comments, func(application *mvc.Application) {
		application.Register(
			database.GetDB,
			service.NewCommentService,
		)
		application.Handle(new(controller.CommentController))
	})

	pages := app.Party("/api/v1/pages")
	pages.Use(middleware.JWT)

	mvc.Configure(pages, func(application *mvc.Application) {
		application.Handle(new(controller.PageController))
	})

	// 用户
	filterForUser := func(ctx iris.Context) bool {
		// 用户注册时跳过JWT验证
		if method := ctx.Method(); method == iris.MethodPost {
			return false
		}
		return true
	}
	jwtMiddlewareForUser := iris.NewConditionalHandler(filterForUser, middleware.JWT)

	users := app.Party("/api/v1/users", jwtMiddlewareForUser)
	mvc.Configure(users, func(application *mvc.Application) {
		application.Register(
			database.GetDB,
			service.NewUserService,
		)
		application.Handle(new(controller.UserController))
	})

	// 设置
	settings := app.Party("/api/v1/settings")
	if conf.GetConfig().Mode == conf.PROD {
		settings.Use(middleware.JWT)
	}
	mvc.Configure(settings, func(app *mvc.Application) {
		app.Register(
			database.GetDB,
		)
		app.Handle(new(controller.SettingController))
	})

	// 标签
	tags := app.Party("/api/v1/tags")
	if conf.GetConfig().Mode == conf.PROD {
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
