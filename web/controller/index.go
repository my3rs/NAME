package controller

import (
	"change/model"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type IndexController struct {
	Ctx iris.Context
}

func (m *IndexController) Get() interface{} {
	content := model.GetContentByPage(1, 10)

	return mvc.View{
		Name: "index.html",
		Data: iris.Map{
			"Posts":   content,
			"isIndex": true,
			"Title":   "网站名",
		},
	}
}
