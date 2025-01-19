package controller

import (
	"NAME/dict"
	"NAME/model"
	"NAME/service"
	"github.com/kataras/iris/v12"
	"strconv"
	"strings"
)

type CategoryController struct {
	Ctx             iris.Context
	CategoryService service.CategoryService
	UserService     service.UserService
}

func (c *CategoryController) Get(req model.QueryRequest) model.PageResponse {
	// 检查请求参数
	if req.PageSize <= 0 || req.PageIndex <= 0 {
		c.Ctx.Application().Logger().Error("Bad request: pageIndex=", req.PageIndex, ",pageSize=", req.PageSize)
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.NewPageResponse(false, dict.ErrInvalidParameters.Error()+": pageSize or pageIndex", nil, req.PageIndex, req.PageSize, 0)
	}

	total := c.CategoryService.GetCategoriesCount()
	if int64(req.PageIndex) > (total/int64(req.PageSize) + 1) {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.NewPageResponse(false, "pageIndex too large", nil, req.PageIndex, req.PageSize, total)
	}

	if req.Order == "" {
		req.Order = "created_at desc"
	}

	// 读取数据
	items := c.CategoryService.GetCategories(req.PageIndex-1, req.PageSize, req.Order)

	return model.NewPageResponse(true, "success", items, req.PageIndex, req.PageSize, total)
}

func (c *CategoryController) Post(req model.Category) model.EmptyResponse {
	if req.Text == "" {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.EmptyResponse{Success: false, Message: "text cannot be empty"}
	}

	if req.Slug == "" {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.EmptyResponse{Success: false, Message: "no cannot be empty"}
	}

	var category = model.Category{
		Text: req.Text,
		Slug: req.Slug,
	}

	err := c.CategoryService.InsertCategory(category)
	if err != nil {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.EmptyResponse{Success: false, Message: err.Error()}
	}

	return model.EmptyResponse{Success: true, Message: "success"}
}

func (c *CategoryController) PutBy(id int) model.EmptyResponse {
	if id <= 0 {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.EmptyResponse{Success: false, Message: "无效ID"}
	}

	var req model.Category
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.EmptyResponse{Success: false, Message: err.Error()}
	}

	if id != int(req.ID) {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.EmptyResponse{Success: false, Message: "无效ID"}
	}

	if req.Text == "" {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.EmptyResponse{Success: false, Message: "text cannot be empty"}
	}

	if req.Slug == "" {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.EmptyResponse{Success: false, Message: "no cannot be empty"}
	}

	var category = model.Category{
		ID:   req.ID,
		Text: req.Text,
		Slug: req.Slug,
	}

	err = c.CategoryService.UpdateCategory(category)
	if err != nil {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.EmptyResponse{Success: false, Message: err.Error()}
	}

	return model.EmptyResponse{Success: true, Message: "success"}
}

func (c *CategoryController) DeleteBy(idsReq string) model.BatchResponse {
	if len(idsReq) == 0 {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.BatchResponse{Success: false, Message: "bad params in url"}
	}

	// 去掉末尾的逗号
	if idsReq[len(idsReq)-1] == ',' {
		idsReq = idsReq[:len(idsReq)-1]
	}

	// 将字符中的ID转换为数组
	var ids []uint
	str := strings.Split(idsReq, ",")
	for _, item := range str {
		if len(item) == 0 {
			continue
		}
		id, err := strconv.Atoi(item)
		if err != nil {
			c.Ctx.StatusCode(iris.StatusBadRequest)
			return model.BatchResponse{Success: false, Message: err.Error()}
		}
		ids = append(ids, uint(id))
	}

	// 调用 Service 进行删除
	err := c.CategoryService.DeleteCategories(ids)
	if err != nil {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.BatchResponse{Success: false, Message: err.Error()}
	}

	c.Ctx.StatusCode(iris.StatusOK)
	return model.NewBatchResponse(true, "success", ids, nil)
}
