package controller

import "github.com/kataras/iris/v12"

func Respond(ctx iris.Context, statusCode int, p iris.Map) {
	ctx.StatusCode(statusCode)
	if err := ctx.JSON(p); err != nil {
		ctx.Application().Logger().Error(err.Error())
	}
	return
}
