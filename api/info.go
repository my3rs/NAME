package api

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type InfoController struct {
	Ctx iris.Context
}

type UserInfoResponse struct {
	ID   int
	Name string
	Role string
}

// GetUser handles http://localhost:8080/api/v1/info/user
func (c *InfoController) GetUser() {
	claims := GetClaims(c.Ctx)

	if claims != nil {
		response := UserInfoResponse{}

		response.ID = claims.ID
		response.Name = claims.Username
		response.Role = claims.Role

		c.Ctx.JSON(response)
	}

}

func GetUserHandler() iris.Handler {
	return func(context *context.Context) {
		claims := GetClaims(context)

		response := UserInfoResponse{claims.ID, claims.Username, claims.Role}

		context.JSON(response)
	}
}
