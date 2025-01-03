package controller

import (
	"NAME/auth"
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
	var json LoginJSON
	if err := c.Ctx.ReadJSON(&json); err != nil {
		c.Ctx.Application().Logger().Infof("Failed to read json from request: %s", err.Error())
		return loginResponse{
			Success: false,
			Message: err.Error(),
		}
	}

	user, err := c.UserService.GetUserByName(username)
	if err != nil {
		return loginResponse{
			Success:           false,
			Message:           "Invalid username or password",
			IsUsernameInvalid: true,
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(json.Password)); err != nil {
		return loginResponse{
			Success:           false,
			Message:           "Invalid username or password",
			IsPasswordInvalid: true,
		}
	}

	tokenPair, err := auth.GetJWTService().GenerateTokenPair(user)
	if err != nil {
		return loginResponse{
			Success: false,
			Message: "Failed to generate token",
		}
	}

	c.Ctx.Header("Authorization", "Bearer "+string(tokenPair.AccessToken))
	c.Ctx.Header("X-Refresh-Token", string(tokenPair.RefreshToken))

	return loginResponse{
		Success: true,
		Message: "Login successful",
	}
}

// PostRefresh handles POST: https://localhost/api/v1/auth/refresh
func (c *AuthController) PostRefresh() {
	refreshToken := auth.GetJWTService().GetTokenFromHeader(c.Ctx, auth.TypeRefreshToken)
	claims, err := auth.GetJWTService().VerifyToken(refreshToken)
	if err != nil {
		c.Ctx.StopWithJSON(401, iris.Map{
			"message": "Invalid refresh token",
		})
		return
	}

	user, err := c.UserService.GetUserByName(claims.Username)
	if err != nil {
		c.Ctx.StopWithJSON(401, iris.Map{
			"message": "User not found",
		})
		return
	}

	tokenPair, err := auth.GetJWTService().GenerateTokenPair(user)
	if err != nil {
		c.Ctx.StopWithJSON(500, iris.Map{
			"message": "Failed to generate token pair",
		})
		return
	}

	c.Ctx.Header("Authorization", "Bearer "+string(tokenPair.AccessToken))
	c.Ctx.Header("X-Refresh-Token", string(tokenPair.RefreshToken))
	c.Ctx.JSON(iris.Map{
		"message": "Token refreshed successfully",
	})
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
