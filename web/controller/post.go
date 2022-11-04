package controller

import (
	"github.com/kataras/iris/v12"
)

type PostController struct {
	Ctx iris.Context
}

func NewPostController() *PostController {
	return &PostController{}
}

func (m *PostController) Get() string {
	return "hey"
}

//func (m *PostController) GetBy(id int64) interface{} {
//	post, exist := model.GetContentById(id)
//	if exist {
//		html := model.Markdown2Html(post.Text)
//		return mvc.View{
//			Name: "post.html",
//			Data: iris.Map{
//				"PostHtml": html,
//				"isPost":   true,
//			},
//		}
//	} else {
//		return 404
//	}
//}

func (m *PostController) GetHello() interface{} {
	return 404
}
