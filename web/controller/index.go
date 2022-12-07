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
		Name: "layout.tmpl",
		Data: iris.Map{
			"Posts":   posts,
			"isIndex": true,
			"isPost":  false,
			"Title":   "听海",
		},
	}
}
