package controller

import (
	"NAME/model"
	"NAME/service"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	Ctx         iris.Context
	UserService service.UserService
}

// 移除字符串两端的引号
func trimQuotes(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}

// PostLoginBy handles POST: /auth/login/:username
func (c *AuthController) PostLoginBy(username string) model.BaseResponse {
	// Read json from request body
	var json model.User
	if err := c.Ctx.ReadJSON(&json); err != nil {
		c.Ctx.Application().Logger().Infof("Failed to read json from request: %s", err.Error())
		c.Ctx.StatusCode(iris.StatusBadRequest)

		return model.NewResponse(false, err.Error())
	}

	// Check `username`s in URL and request body
	if json.Username != username {
		c.Ctx.Application().Logger().Infof("username in URL and body doesn't match: %s %s",
			username, json.Username)
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.NewResponse(false, "username in URL and body doesn't match")
	}

	user, err := c.UserService.GetUserByName(username)
	if err != nil {
		return model.NewResponse(false, "user not found")

	}

	c.Ctx.Application().Logger().Infof("user %s login", username)

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(json.Password)); err != nil {
		c.Ctx.Application().Logger().Infof("user %s login with password %s", username, json.Password)
		return model.NewResponse(false, "invalid username or password")
	}

	// 生成token对
	tokenPair, err := service.GetJWTService().GenerateTokenPair(user)
	if err != nil {
		c.Ctx.Application().Logger().Infof("Failed to generate token pair: %s", err.Error())

		return model.NewResponse(false, "failed to generate token pair")
	}

	c.Ctx.Header("Authorization", "Bearer "+trimQuotes(string(tokenPair.AccessToken)))
	c.Ctx.Header("X-Refresh-Token", trimQuotes(string(tokenPair.RefreshToken)))

	return model.BaseResponse{
		Success: true,
		Message: "success",
	}
}

// PostRefresh handles POST: /auth/refresh
func (c *AuthController) PostRefresh() model.BaseResponse {
	claims, err := service.GetJWTService().VerifyRefreshToken(c.Ctx)
	if err != nil {
		c.Ctx.StatusCode(iris.StatusUnauthorized)
		return model.NewResponse(false, "Invalid refresh token: "+err.Error())
	}

	user, err := c.UserService.GetUserByName(claims.Subject)
	if err != nil {
		c.Ctx.StatusCode(iris.StatusUnauthorized)
		return model.NewResponse(false, "User not found")
	}

	tokenPair, err := service.GetJWTService().GenerateTokenPair(user)
	if err != nil {
		c.Ctx.StatusCode(iris.StatusInternalServerError)
		return model.NewResponse(false, "Failed to generate token pair")
	}

	c.Ctx.Header("Authorization", "Bearer "+string(tokenPair.AccessToken))
	c.Ctx.Header("X-Refresh-Token", string(tokenPair.RefreshToken))
	
	return model.NewResponse(true, "Token refreshed successfully")
}

// HashAndSalt : Generate hashed password
// @password: plain password
func HashAndSalt(password []byte) string {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {

	}
	return string(hash)
}

// GetClaims returns the current authorized client claims.
func GetClaims(ctx iris.Context) *model.Claims {
	claims := jwt.Get(ctx).(*model.Claims)
	return claims
}
