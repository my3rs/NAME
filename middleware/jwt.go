package middleware

import (
	"NAME/controller"
	"NAME/service"
	"github.com/kataras/iris/v12"
)

func JWT(ctx iris.Context) {
	err := service.GetJWTService().VerifyAccessToken(ctx)

	if err != nil {
		ctx.Application().Logger().Error(err.Error())
		controller.Respond(ctx, iris.StatusUnauthorized, iris.Map{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.Next()
}
