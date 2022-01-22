package api

import (
	"change/conf"
	"change/model"
	"errors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/kataras/iris/v12/mvc"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserClaims struct {
	ID       int    `json:id`
	Username string `json:username`
	Role     string `json:"role"`
}

type Controller struct {
	Ctx iris.Context
}

var (
	signer             = jwt.NewSigner(jwt.HS256, conf.GetSecretKey(), 15*time.Minute)
	Verifier           = jwt.NewVerifier(jwt.HS256, conf.GetSecretKey())
	VerifierMiddleware = Verifier.Verify(func() interface{} {
		return new(UserClaims)
	})
)

// PostRegister handles POST: http://localhost/api/v1/auth/register
func (c *Controller) PostRegister() mvc.Result {
	var (
		name     = c.Ctx.FormValue("name")
		mail     = c.Ctx.FormValue("mail")
		password = c.Ctx.FormValue("password")
	)
	u, err := model.CreateUser(model.User{Name: name, Mail: mail, HashedPassword: HashAndSalt([]byte(password))})
	if err != nil {
		// todo
	}

	c.Ctx.Application().Logger().Infof("Register: %s %s %s", u.Name, u.Mail, password)

	return mvc.Response{
		Err: nil,
		// redirect to /user/index
		Path: "/user/index",
	}
}

// PostLogin handles POST: https://localhost/api/v1/auth/login
func (c *Controller) PostLogin() {
	response := ResponseJSON{}

	var (
		username         = c.Ctx.FormValue("name")
		password         = c.Ctx.FormValue("password")
		user, userExists = model.GetUserByName(username)
		claims           = UserClaims{}
	)

	if userExists {
		claims.ID = user.ID
		claims.Username = user.Name
		claims.Role = user.Role
	}

	if err := VerityUserByName(username, password); err != nil {
		response.Message = err.Error()
		// 用户输入的可能是用户名，也可能是邮箱
		if err := VerityUserByMail(username, password); err != nil {
			response.Code = 401
			response.Message += "&" + err.Error()
			c.Ctx.JSON(response)
			return
		}
	}

	token, err := signer.Sign(claims)
	if err != nil {
		c.Ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}
	response.Code = 0
	response.Message = "用户验证成功，请从http header中获取token"
	response.Path = "/admin/index"

	c.Ctx.Header("Authorization", "Bearer "+string(token))
	c.Ctx.JSON(response)

	//c.Ctx.SetCookieKV("Authorization", string(token))
	//
	//return mvc.Response{
	//	Code: 302,
	//	Path: "/admin/index",
	//}
}

func VerityUserByName(username string, password string) error {
	user, succ := model.GetUserByName(username)
	if succ != true {
		return errors.New("用户名不存在")
	}

	if succ := user.ValidatePassword([]byte(password)); succ != true {
		return errors.New("密码与用户名不匹配")
	}

	return nil
}

func VerityUserByMail(mail string, password string) error {
	user, succ := model.GetUserByMail(mail)
	if succ != true {
		return errors.New("邮箱不存在")
	}

	if succ = user.ValidatePassword([]byte(password)); succ != true {
		return errors.New("密码与邮箱不匹配")
	}

	return nil
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
