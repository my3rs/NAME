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
	ID       uint   `json:id`
	Username string `json:username`
	Role     string `json:role`
}

type RegisterJSON struct {
	Username string `json:username`
	Mail     string `json:mail`
	Password string `json:password`
}

type LoginJSON struct {
	Username string `json:username`
	Password string `json:password`
}

type AuthController struct {
	Ctx         iris.Context
	UserService service.UserService
}

func NewAuthController() *AuthController {
	return &AuthController{UserService: service.NewUserService()}
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

// PostLoginBy handles POST: https://localhost/api/v1/auth/login/:username
func (c *AuthController) PostLoginBy(username string) {
	// Read json from request body
	var json LoginJSON
	if err := c.Ctx.ReadJSON(&json); err != nil {
		c.Ctx.Application().Logger().Infof("Failed to read json from request: ", err.Error())
		c.Ctx.StatusCode(400)
		c.Ctx.JSON(iris.Map{"message": err.Error()})
		return
	}

	// Check `username` in URL and request body
	if json.Username != username {
		c.Ctx.Application().Logger().Infof("usernames in URL and body doesn't match: %s %s",
			username, json.Username)
		c.Ctx.StatusCode(400)
		c.Ctx.JSON(iris.Map{"message": "usernames in URL and body doesn't match"})
		return
	}

	pair, err := c.UserService.Login(json.Username, json.Password)
	if err != nil {
		c.Ctx.Application().Logger().Infof("Failed to login: %v", err.Error())
		c.Ctx.StopWithJSON(401, iris.Map{"message": err.Error()})

		return
	}

	c.Ctx.StatusCode(200)
	c.Ctx.Header("authorization", string(pair.AccessToken))
	c.Ctx.Header("refresh-token", string(pair.RefreshToken))
	c.Ctx.JSON(iris.Map{"message": "login success: check tokens in the HTTP header"})
	return
}

// PostRefresh handles POST: https://localhost/api/v1/auth/refresh
// @func: verify the refresh token and then generate a new token pair, both access token and refresh token
func (c *AuthController) PostRefresh() {
	// Verify refresh token and get username from the token
	username := middleware.VerifyRefreshTokenAndGetUserName(c.Ctx)
	if len(username) == 0 {
		c.Ctx.StopWithJSON(401, iris.Map{"message": "invalid refresh token"})
		return
	}

	user, err := c.UserService.GetUserByName(username)
	if err != nil {
		c.Ctx.StopWithJSON(401, iris.Map{"message": err.Error()})
		return
	}

	// Generate new token pair
	tokens, err := middleware.GenerateTokenPair(user.ID, user.Name)
	if err != nil {
		c.Ctx.StopWithJSON(401, iris.Map{"message": "failed to refresh token"})
		return
	}
	c.Ctx.StatusCode(200)
	c.Ctx.Header("authorization", string(tokens.AccessToken))
	c.Ctx.Header("refresh-token", string(tokens.RefreshToken))
	c.Ctx.JSON(iris.Map{"message": "refresh success: check new tokens in the HTTP header"})

	return
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
