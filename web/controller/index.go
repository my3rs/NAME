package controller

import (
	"NAME/service"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type IndexController struct {
	Ctx     iris.Context
	Service service.ContentService
}

func NewIndexController() *IndexController {
	return &IndexController{Service: service.NewContentService()}
}

func (m *IndexController) Get() interface{} {
	posts := m.Service.GetPostsWithOrder(1, 10, "created_at desc")

	return mvc.View{
		Name: "index.html",
		Data: iris.Map{
			"Posts":   posts,
			"isIndex": true,
			"Title":   "网站名",
		},
	}
}
