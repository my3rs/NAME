package controller

import (
	"NAME/auth"
	"NAME/dict"
	"NAME/model"
	"NAME/service"

	"github.com/kataras/iris/v12"
)

type UserController struct {
	Ctx         iris.Context
	UserService service.UserService
}

func (c *UserController) Get(req model.QueryRequest) model.PageResponse {
	if req.PageSize <= 0 || req.PageIndex <= 0 {
		c.Ctx.Application().Logger().Info("request: pageIndex=", req.PageIndex, ",pageSize=", req.PageSize)
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.NewPageResponse(
			false,
			dict.ErrInvalidParameters.Error()+": pageSize or pageIndex",
			nil,
			req.PageIndex,
			req.PageSize,
			0,
		)
	}

	// 获取总数
	total := c.UserService.GetUserNum()
	if total == 0 {
		return model.NewPageResponse(
			true,
			"success",
			[]model.User{},
			req.PageIndex,
			req.PageSize,
			0,
		)
	}

	// 检查页码
	totalPages := total/int64(req.PageSize) + 1
	if int64(req.PageIndex) > totalPages {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.NewPageResponse(
			false,
			"pageIndex too large",
			nil,
			req.PageIndex,
			req.PageSize,
			total,
		)
	}

	// 设置默认排序
	if len(req.Order) == 0 {
		req.Order = "id asc"
	}

	// 获取用户列表
	users, err := c.UserService.GetUsersWithOrder(req.PageIndex-1, req.PageSize, req.Order)
	if err != nil {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.NewPageResponse(
			false,
			err.Error(),
			nil,
			req.PageIndex,
			req.PageSize,
			total,
		)
	}

	return model.NewPageResponse(
		true,
		"success",
		users,
		req.PageIndex,
		req.PageSize,
		total,
	)
}

// GetMe 获取当前登录用户信息
// handle GET /api/v1/users/me
func (c *UserController) GetMe() model.DetailResponse {
	// 从 JWT Claims 获取当前用户信息
	claims := auth.GetJWTService().GetClaimsFromContext(c.Ctx)
	if claims == nil {
		c.Ctx.StatusCode(iris.StatusUnauthorized)
		return model.DetailResponse{
			Success: false,
			Message: "未登录",
		}
	}

	// 获取用户信息
	user, err := c.UserService.GetUserByID(int(claims.ID))
	if err != nil {
		c.Ctx.Application().Logger().Error("Failed to get user info: ", err.Error())
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.DetailResponse{
			Success: false,
			Message: err.Error(),
		}
	}

	return model.DetailResponse{
		Success: true,
		Message: "success",
		Data:    user,
	}
}
