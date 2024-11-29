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

type postRequest struct {
	ID   uint   `json:"id"`
	Text string `json:"text"`
	No   string `json:"no"`
}

func (c *CategoryController) Get(req model.QueryRequest) iris.Map {
	// 检查请求参数
	if req.PageSize <= 0 || req.PageIndex <= 0 {
		c.Ctx.Application().Logger().Error("Bad request: pageIndex=", req.PageIndex, ",pageSize=", req.PageSize)

		c.Ctx.StatusCode(iris.StatusBadRequest)
		return iris.Map{
			"success": false,
			"message": dict.ErrInvalidParameters.Error() + ": pageSize or pageIndex",
		}
	}

	// 构造分页
	var page model.Page
	page.PageIndex = req.PageIndex
	page.PageSize = req.PageSize
	page.ItemsCount = c.CategoryService.GetCategoriesCount()
	page.PageCount = page.ItemsCount/int64(req.PageSize) + 1

	if int64(req.PageIndex) > page.PageCount {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return iris.Map{
			"success": false,
			"message": "pageIndex too large",
		}
	}

	if req.Order == "" {
		req.Order = "created_at desc"
	}

	page.Order = req.Order

	// 读取数据并返回
	items := c.CategoryService.GetCategories(req.PageIndex-1, req.PageSize, req.Order)

	return iris.Map{
		"success":    true,
		"items":      items,
		"pagination": &page,
	}
}

// Post handles POST /api/v1/categories
func (c *CategoryController) Post(req postRequest) iris.Map {
	if req.Text == "" || len(req.Text) == 0 {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return iris.Map{
			"success": false,
			"message": "text cannot be empty",
		}
	}

	// 后端不允许空的 slug 值
	// TODO: 在前端检测到中文标题时，自动生成其拼音作为 slug 值
	if req.No == "" || len(req.No) == 0 {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return iris.Map{
			"success": false,
			"message": "no cannot be empty",
		}
	}

	var category = model.Category{
		Text: req.Text,
		No:   req.No,
	}

	err := c.CategoryService.InsertCategory(category)
	if err != nil {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return iris.Map{
			"success": false,
			"message": err,
		}
	}

	return iris.Map{
		"success": true,
	}

}

// PutBy handles POST /api/v1/categories/{id:int}
func (c *CategoryController) PutBy(id int) iris.Map {
	if id <= 0 {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return iris.Map{
			"success": false,
			"message": "无效ID",
		}
	}

	var req postRequest
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return iris.Map{
			"success": false,
			"message": err,
		}
	}

	if id != int(req.ID) {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return iris.Map{
			"success": false,
			"message": "无效ID",
		}
	}

	if req.Text == "" || len(req.Text) == 0 {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return iris.Map{
			"success": false,
			"message": "text cannot be empty",
		}
	}

	// 后端不允许空的 slug 值
	// TODO: 在前端检测到中文标题时，自动生成其拼音作为 slug 值
	if req.No == "" || len(req.No) == 0 {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return iris.Map{
			"success": false,
			"message": "no cannot be empty",
		}
	}

	var category = model.Category{
		ID:   req.ID,
		Text: req.Text,
		No:   req.No,
	}

	err = c.CategoryService.UpdateCategory(category)
	if err != nil {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return iris.Map{
			"success": false,
			"message": err,
		}
	}

	return iris.Map{
		"success": true,
	}
}

// DeleteBy handles DELETE /api/v1/tag/1,2,3 批量删除标签
func (c *CategoryController) DeleteBy(idsReq string) iris.Map {
	if len(idsReq) == 0 {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return iris.Map{
			"success": false,
			"message": "bad params in url",
		}
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
			return iris.Map{
				"success": false,
				"message": err.Error(),
			}
		}
		ids = append(ids, uint(id))
	}

	// 调用 Service 进行删除
	err := c.CategoryService.DeleteCategories(ids)
	if err != nil {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return iris.Map{
			"success": false,
			"message": err,
		}
	}

	c.Ctx.StatusCode(iris.StatusOK)
	return iris.Map{
		"success": true,
	}
}
