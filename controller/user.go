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
	claims, err := service.GetJWTService().GetClaimsFromContext(c.Ctx)
	if err != nil {
		c.Ctx.StatusCode(iris.StatusUnauthorized)
		c.Ctx.JSON(iris.Map{"message": err.Error()})
		return model.DetailResponse{
			Success: false,
			Message: err.Error(),
		}
	}

	// 获取用户信息
	user, err := c.UserService.GetUserByName(claims.Subject)
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

// Post 处理用户注册
// handle POST /api/v1/users
func (c *UserController) Post(user model.User) model.DetailResponse {
	// 验证必要的请求参数
	if user.Username == "" || user.Password == "" || user.Mail == "" {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.DetailResponse{
			Success: false,
			Message: dict.ErrInvalidParameters.Error() + ": username, password and mail are required",
		}
	}

	// 加密密码
	user.HashedPassword = HashAndSalt([]byte(user.Password))

	// 设置默认值
	user.Role = model.UserRoleReader
	user.Activated = false

	// 调用服务层创建用户
	err := c.UserService.InsertUser(user)
	if err != nil {
		c.Ctx.StatusCode(iris.StatusInternalServerError)
		return model.DetailResponse{
			Success: false,
			Message: "Failed to create user: " + err.Error(),
		}
	}

	newUser, err := c.UserService.GetUserByName(user.Username)

	// 返回成功响应
	return model.DetailResponse{
		Success: true,
		Message: "User registered successfully",
		Data:    newUser,
	}
}
