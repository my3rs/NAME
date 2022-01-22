package controller

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type DashboardController struct {
	Ctx iris.Context
}

func (c *DashboardController) Get() mvc.Result {
	return c.GetIndex()
}

func (c *DashboardController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/new_post", "GetNewPost")
}

// GetIndex handles http://localhost:8080/admin/index
func (c *DashboardController) GetIndex() mvc.Result {
	//tokenFromHeader := c.Ctx.GetHeader("Authorization")
	//token := strings.Split(tokenFromHeader, " ")[1]
	//_, err := jwt.Verify(jwt.HS256, conf.GetSecretKey(), []byte(token))
	//if err != nil {
	//	c.Ctx.Application().Logger().Info("Unauthorized ", token)
	//	c.Ctx.Application().Logger().Info(err)
	//	return mvc.Response{
	//		Path: "/user/login",
	//	}
	//}

	//c.Ctx.Application().Logger().Info("Succeed to authorize")
	indexView := mvc.View{
		Name: "user/index.html",
		Data: iris.Map{},
	}

	return indexView
}

func (c *DashboardController) GetNewPost() mvc.Result {
	newPostView := mvc.View{
		Name: "/user/new_post.html",
		Data: iris.Map{},
	}

	return newPostView
}
