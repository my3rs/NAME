package controller

import (
	"NAME/service"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type PostController struct {
	Ctx            iris.Context
	ContentService service.ContentService
}

func NewPostController() *PostController {
	return &PostController{ContentService: service.NewContentService()}
}

func (c *PostController) GetBy(id int) mvc.View {
	post := c.ContentService.GetFormattedPostByID(id)
	return mvc.View{
		Name: "post.gohtml",
		Data: iris.Map{
			"Post":   post,
			"Title":  post.Title,
			"isPost": true,
		},
	}
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
