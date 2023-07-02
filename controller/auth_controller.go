package controller

import (
	"NAME/middleware"
	"NAME/model"
	"NAME/service"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserClaims struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type RegisterJSON struct {
	Username string `json:"username"`
	Mail     string `json:"mail"`
	Password string `json:"password"`
}

type LoginJSON struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Success           bool   `json:"success"`
	Message           string `json:"message"`
	IsUsernameInvalid bool   `json:"isUsernameInvalid,omitempty"`
	IsPasswordInvalid bool   `json:"isPasswordInvalid,omitempty"`
}

type AuthController struct {
	Ctx         iris.Context
	UserService service.UserService
}

func NewAuthController() *AuthController {
	return &AuthController{UserService: service.GetUserService()}
}

// PostRegisterBy handles POST: https://localhost/api/v1/auth/register/:username
func (c *AuthController) PostRegisterBy(username string) {
	// Read json from request body
	var json RegisterJSON
	if err := c.Ctx.ReadJSON(&json); err != nil {
		c.Ctx.Application().Logger().Infof("Failed to read json from request: ", err.Error())

		c.Ctx.StatusCode(400)
		c.Ctx.JSON(iris.Map{"message": err.Error()})
		return
	}

	// Check `username`s in URL and request body
	if json.Username != username {
		c.Ctx.Application().Logger().Infof("username in URL and body doesn't match: %s %s",
			username, json.Username)

		c.Ctx.StatusCode(400)
		c.Ctx.JSON(iris.Map{"message": "username in URL and body doesn't match"})
		return
	}

	// Create user in database
	err := c.UserService.InsertUser(model.User{Name: json.Username, Mail: json.Mail,
		HashedPassword: HashAndSalt([]byte(json.Password))})

	if err != nil {
		c.Ctx.StopWithJSON(400, iris.Map{"message": err.Error()})
		return
	}

	c.Ctx.StatusCode(200)
	c.Ctx.JSON(iris.Map{"message": "success"})

}

func trimQuotes(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}

// PostLoginBy handles POST: https://localhost/api/v1/auth/login/:username
func (c *AuthController) PostLoginBy(username string) loginResponse {
	// Read json from request body
	var json LoginJSON
	if err := c.Ctx.ReadJSON(&json); err != nil {
		c.Ctx.Application().Logger().Infof("Failed to read json from request: ", err.Error())
		c.Ctx.StatusCode(400)
		return loginResponse{
			Success: false,
			Message: err.Error(),
		}
	}

	// Check `username` in URL and request body
	if json.Username != username {
		c.Ctx.Application().Logger().Infof("usernames in URL and body doesn't match: %s %s",
			username, json.Username)
		c.Ctx.StatusCode(400)
		c.Ctx.JSON(iris.Map{"message": "usernames in URL and body doesn't match"})
		return loginResponse{
			Success:           false,
			Message:           "usernames in URL and body don't match",
			IsUsernameInvalid: true,
		}
	}

	// read `User` from database
	user, err := c.UserService.GetUserByName(json.Username)
	if err != nil {
		c.Ctx.StatusCode(iris.StatusBadRequest)

		return loginResponse{
			Success: false,
			Message: err.Error(),
		}
	}

	err = c.UserService.VerifyPassword(user, json.Password)
	if err != nil {
		return loginResponse{
			Success: false,
			Message: err.Error(),
		}
	}

	var pair jwt.TokenPair
	pair, err = middleware.GenerateTokenPair(user)
	if err != nil {
		return loginResponse{
			Success: false,
			Message: err.Error(),
		}
	}

	c.Ctx.StatusCode(200)
	c.Ctx.Header("authorization", trimQuotes(string(pair.AccessToken)))
	c.Ctx.Header("refresh-token", trimQuotes(string(pair.RefreshToken)))

	return loginResponse{
		Success:           true,
		IsUsernameInvalid: false,
		IsPasswordInvalid: false,
		Message:           "login success: check tokens in the HTTP header",
	}
}

// PostRefresh handles POST: https://localhost/api/v1/auth/refresh
// @func: verify the refresh token and then generate a new token pair,
// both access token and refresh token
func (c *AuthController) PostRefresh() {
	middleware.RefreshToken(c.Ctx)
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
func GetClaims(ctx iris.Context) *UserClaims {
	claims := jwt.Get(ctx).(*UserClaims)
	return claims
}
