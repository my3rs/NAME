package controller

import (
	"NAME/dict"
	"NAME/model"
	"NAME/service"
	"github.com/kataras/iris/v12"
)

type UserController struct {
	Ctx         iris.Context
	UserService service.UserService
}

func (c *UserController) Get(req model.QueryRequest) model.TestResponse {
	if req.PageSize <= 0 || req.PageIndex <= 0 {
		c.Ctx.Application().Logger().Info("request: pageIndex=", req.PageIndex, ",pageSize=", req.PageSize)
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.TestResponse{Success: false, Message: dict.ErrInvalidParameters.Error() + ": pageSize or pageIndex"}
	}

	var rsp model.TestResponse
	var page model.Page

	page.ItemsCount = c.UserService.GetUserNum()
	page.PageIndex = req.PageIndex
	page.PageSize = req.PageSize
	page.PageCount = page.ItemsCount/int64(req.PageSize) + 1

	if int64(req.PageIndex) > page.PageCount {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.TestResponse{Success: false, Message: "pageIndex too large"}
	}

	if len(req.Order) == 0 {
		req.Order = "id asc"
	}
	page.Order = req.Order

	users, err := c.UserService.GetUsersWithOrder(req.PageIndex-1, req.PageSize, req.Order)
	if err != nil {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.TestResponse{Success: false, Message: err.Error()}
	}

	rsp.Data = users
	rsp.Page = &page

	rsp.Success = true

	return rsp

}
