package controller

import (
	"NAME/model"
	"NAME/service"
	"github.com/kataras/iris/v12"
	"log"
)

type getTagsRequest struct {
	Format string `url:"format" json:"format"`
}

type TagController struct {
	Ctx     iris.Context
	Service service.TagService
}

type GetTagsResponse struct {
	Success bool        `json:"success"`
	Data    []model.Tag `json:"data"`
}

func (c *TagController) Get(request getTagsRequest) GetTagsResponse {
	var tags []model.Tag

	if c.Service == nil {
		log.Printf("tagController addr: %p\n", c)
		log.Panic("c.Service is nil")
	}

	switch request.Format {
	case "withPath":
		tags = c.Service.GetAllTagsWithPath()

	default:
		tags = c.Service.GetAllTags()
	}

	c.Ctx.StatusCode(iris.StatusOK)
	return GetTagsResponse{
		Success: true,
		Data:    tags,
	}
}
