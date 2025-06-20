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

func (c *TagController) Get() model.ListResponse {
	tags := c.Service.GetAllTags()
	
	// 转换为interface{}切片
	data := make([]interface{}, len(tags))
	for i, tag := range tags {
		data[i] = tag
	}

	return model.NewListResponse(true, "success", data, int64(len(tags)))
}

func (c *TagController) Post(req model.Tag) model.EmptyResponse {
	var tag = model.Tag{
		Slug: req.Slug,
		Text: req.Text,
	}

	err := c.Service.InsertTag(tag)
	if err != nil {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.EmptyResponse{
			Success: false,
			Message: err.Error(),
		}
	}

	return model.EmptyResponse{
		Success: true,
		Message: "新建标签成功！",
	}
}

// DeleteBy handles DELETE /api/v1/tag/1,2,3 批量删除标签
func (c *TagController) DeleteBy(idsReq string) model.EmptyResponse {
	if len(idsReq) == 0 {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.EmptyResponse{
			Success: false,
			Message: "参数错误",
		}
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
			return model.EmptyResponse{
				Success: false,
				Message: err.Error(),
			}
		}
		ids = append(ids, uint(id))
	}

	err := c.Service.DeleteTags(ids)
	if err != nil {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.EmptyResponse{
			Success: false,
			Message: err.Error(),
		}
	}

	return model.EmptyResponse{
		Success: true,
		Message: "删除成功",
	}
}
