package controller

import (
	"NAME/model"
	"NAME/service"
	"github.com/kataras/iris/v12"
	"log"
)

type getTagsRequest struct {
	Meta string `url:"meta" json:"meta"`
	Path string `url:"path" json:"path"`
}

type newTagRequest struct {
	No       string `json:"no"`
	Text     string `json:"text"`
	ParentID uint   `json:"parentID"`
}

type TagController struct {
	Ctx     iris.Context
	Service service.TagService
}

type GetTagsResponse struct {
	Success bool        `json:"success"`
	Data    []model.Tag `json:"data"`
	Meta    interface{} `json:"meta,omitempty"`
}

func (c *TagController) Get(request getTagsRequest) GetTagsResponse {
	var tags []model.Tag

	if c.Service == nil {
		log.Printf("tagController addr: %p\n", c)
		log.Panic("c.Service is nil")
	}

	if request.Path != "" {
		tags = c.Service.GetAllTagsWithPath()
	} else {
		tags = c.Service.GetAllTags()
	}

	var rsp GetTagsResponse
	rsp.Success = true
	rsp.Data = tags

	if request.Meta != "" {
		meta := c.Service.GetMetadata()
		rsp.Meta = meta
	}

	c.Ctx.StatusCode(iris.StatusOK)
	return rsp
}

func (c *TagController) Post(req newTagRequest) {
	var tag = model.Tag{
		No:       req.No,
		Text:     req.Text,
		ParentID: req.ParentID,
	}

	err := c.Service.InsertTag(tag)
	if err != nil {
		c.Ctx.StatusCode(400)
		c.Ctx.JSON(iris.Map{
			"Success": false,
			"Message": err.Error(),
		})
		return
	}

	c.Ctx.StatusCode(iris.StatusOK)
	c.Ctx.JSON(iris.Map{
		"Success": true,
		"Message": "新建标签成功！",
	})
}
