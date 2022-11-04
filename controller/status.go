package controller

import (
	"NAME/service"
	"github.com/kataras/iris/v12"
)

type StatusController struct {
	Ctx            iris.Context
	ContentService service.ContentService
	TagService     service.TagService
}

func NewStatusController() *StatusController {
	return &StatusController{
		ContentService: service.NewContentService(),
		TagService:     service.NewTagService(),
	}
}

func (c *StatusController) Get() iris.Map {
	postsCount := c.ContentService.GetPostsCount()
	pagesCount := c.ContentService.GetPageCount()
	tagsCount := c.TagService.GetTagsCount()

	c.Ctx.StatusCode(iris.StatusOK)
	return iris.Map{"success": true,
		"data": iris.Map{
			"postsNum": postsCount,
			"pagesNum": pagesCount,
			"tagsNum":  tagsCount,
		},
	}

}
