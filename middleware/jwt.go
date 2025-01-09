package middleware

import (
	"NAME/service"
	"github.com/kataras/iris/v12"
)

func JWT(ctx iris.Context) {
	err := service.GetJWTService().VerifyAccessToken(ctx)
	if err != nil {
		ctx.Application().Logger().Error(err.Error())
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.Next()
}
