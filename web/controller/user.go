package controller

import (
	"change/model"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"golang.org/x/crypto/bcrypt"
)

// UserController : handles following requests:
// 	GET 	/user/register
//  POST 	/user/register
//	GET 	/user/login
//	POST	/user/login
//	Get		/user/me
//	All HTTP methods	/user/logout
type UserController struct {
	Ctx iris.Context
}

type UserClaims struct {
	Username string `json:username`
}

var (
	registerStaticView = mvc.View{
		Name: "user/register.html",
		Data: iris.Map{"Title": "用户注册"},
	}

	userIndexStaticView = mvc.View{
		Name: "user/index.html",
		Data: iris.Map{},
	}

	loginStaticView = mvc.View{
		Name: "user/login.html",
		Data: iris.Map{"Title": "用户登录"},
	}
)

func NewUserController() *UserController {
	return &UserController{}
}

func (c *UserController) isLoggedIn() bool {
	//return c.GetCurrentUserId() > 0
	return true
}

// GetLogin handles  GET: http://localhost/user/login
func (c *UserController) GetLogin() mvc.Result {
	return loginStaticView
}

// GetRegister handles  GET: http://localhost/user/register
func (c *UserController) GetRegister() mvc.Result {
	return registerStaticView
}

// PostRegister handles POST: http://localhost/user/register
func (c *UserController) PostRegister() mvc.Result {
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
		Path: "/user/login",
	}
}

// HashAndSalt : Generate hashed password
// @password: plain password
func HashAndSalt(password []byte) string {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {

	}
	return string(hash)
}
