package controller

import (
	"NAME/model"
	"NAME/service"
	"github.com/kataras/iris/v12"
	"strconv"
	"strings"
)

type TagController struct {
	Ctx     iris.Context
	Service service.TagService
}

func (c *TagController) Get() {
	var tags []model.Tag

	tags = c.Service.GetAllTags()

	Respond(c.Ctx, iris.StatusOK, iris.Map{
		"success": true,
		"items":   tags,
	})
}

func (c *TagController) Post(req model.Tag) {
	var tag = model.Tag{
		Slug: req.Slug,
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

	Respond(c.Ctx, iris.StatusOK, iris.Map{
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
			Respond(c.Ctx, iris.StatusBadRequest, iris.Map{
				"Success": false,
				"Message": err.Error(),
			})

			return
		}
		ids = append(ids, uint(id))
	}

	err := c.Service.DeleteTags(ids)
	if err != nil {
		Respond(c.Ctx, iris.StatusBadRequest, iris.Map{
			"Success": false,
			"Message": err.Error(),
		})

		return
	}

	Respond(c.Ctx, iris.StatusOK, iris.Map{
		"Success": true,
		"Message": "success",
	})
}
