package controller

import (
	"NAME/model"
	"NAME/service"
	"github.com/kataras/iris/v12"
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

func (c *TagController) Get(request getTagsRequest) {
	var meta iris.Map
	if request.Meta != "" {
		meta = c.Service.GetMetadata()
	}

	if request.Path != "" {
		tags := c.Service.GetAllTagsWithPath()

		c.Ctx.StatusCode(iris.StatusOK)
		c.Ctx.JSON(iris.Map{
			"success": true,
			"data":    tags,
			"meta":    meta,
		})

	} else {
		tags := c.Service.GetAllTags()
		c.Ctx.StatusCode(iris.StatusOK)
		c.Ctx.JSON(iris.Map{
			"success": true,
			"data":    tags,
		})
	}

	return
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
