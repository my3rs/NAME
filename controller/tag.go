package controller

import (
	"NAME/model"
	"NAME/service"
	"github.com/kataras/iris/v12"
	"strconv"
	"strings"
)

type getTagsRequest struct {
	Path string `url:"path" json:"path"`
}

type newTagRequest struct {
	No   string `json:"no"`
	Text string `json:"text"`
}

type TagController struct {
	Ctx     iris.Context
	Service service.TagService
}

func (c *TagController) Get(request getTagsRequest) iris.Map {
	var tags []model.Tag

	tags = c.Service.GetAllTags()
	c.Ctx.StatusCode(iris.StatusOK)

	return iris.Map{
		"success": true,
		"items":   tags,
	}
}

func (c *TagController) Post(req newTagRequest) {
	var tag = model.Tag{
		No:   req.No,
		Text: req.Text,
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

// DeleteBy handles DELETE /api/v1/tag/1,2,3 批量删除标签
func (c *TagController) DeleteBy(idsReq string) {
	if len(idsReq) == 0 {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		c.Ctx.JSON(iris.Map{
			"Success": false,
			"Message": "参数错误",
		})
		return
	}

	// 去掉末尾的逗号
	if idsReq[len(idsReq)-1] == ',' {
		idsReq = idsReq[:len(idsReq)-1]
	}

	var ids []uint
	str := strings.Split(idsReq, ",")
	for _, item := range str {
		if len(item) == 0 {
			continue
		}
		id, err := strconv.Atoi(item)
		if err != nil {
			c.Ctx.StatusCode(iris.StatusBadRequest)
			c.Ctx.JSON(iris.Map{
				"Success": false,
				"Message": err.Error(),
			})
			return
		}
		ids = append(ids, uint(id))
	}

	err := c.Service.DeleteTags(ids)
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
		"Message": "删除标签成功！",
	})
}
