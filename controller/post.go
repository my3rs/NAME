package controller

import (
	"NAME/dict"
	"NAME/model"
	"NAME/service"
	"github.com/kataras/iris/v12"
	"math/rand"
	"strconv"
	"strings"
)

type PostController struct {
	Ctx     iris.Context
	Service service.ContentService
}

func NewPostController() *PostController {
	return &PostController{Service: service.NewContentService()}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// Get handles GET /api/v1/posts
// @Summary Get posts list
// @Description Get posts list
// @Accept  json
// @Produce  json
// @Param 	Authorization	header	string	true	"Access token with the prefix `Bearer `"
// @Param   pageSize    query  	int      true        "page size"	 default(1)
// @Param 	pageIndex	query	int 		true		"page index"	default(1)
// @Param	orderBy		query	string  	false 		"order"		default("created_at desc")
// @Success 200 		{object} model.Response	"success"
// @Failure 400 		{object} model.Response "Bad request"
// @Router /api/v1/posts [get]
func (c *PostController) Get(req model.QueryRequest) model.Response {

	if req.PageSize <= 0 || req.PageIndex <= 0 {
		c.Ctx.Application().Logger().Info("request: pageIndex=", req.PageIndex, ",pageSize=", req.PageSize)
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.Response{Success: false, Message: dict.ErrInvalidParameters.Error() + ": pageSize or pageIndex"}
	}

	var rsp model.Response
	var page model.Page

	page.Total = c.Service.GetPostsCount()
	page.PageIndex = req.PageIndex
	page.PageSize = req.PageSize
	page.TotalPage = page.Total/int64(req.PageSize) + 1

	if int64(req.PageIndex) > page.TotalPage {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.Response{Success: false, Message: "pageIndex too large"}
	}

	if req.PageIndex > 1 {
		page.Pre = "http://localhost:8000/api/v1/posts&pageIndex=" + strconv.Itoa(req.PageIndex-1) + "pageSize=" + strconv.Itoa(req.PageSize)
	}
	if int64(req.PageIndex) < page.TotalPage {
		page.Next = "http://localhost:8000/api/v1/posts&pageIndex=" + strconv.Itoa(req.PageIndex+1) + "pageSize=" + strconv.Itoa(req.PageSize)
	}

	if len(req.Order) == 0 {
		req.Order = "id asc"
	}
	page.Order = req.Order

	posts, err := c.Service.GetPostsWithOrder(req.PageIndex-1, req.PageSize, req.Order)
	if err != nil {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.Response{Success: false, Message: err.Error()}
	}

	rsp.Data = posts
	rsp.Page = &page

	rsp.Success = true

	return rsp
}

func (c *PostController) GetBy(id int) model.Response {
	post, err := c.Service.GetContentByID(id)
	if err != nil {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.Response{Success: false, Message: err.Error()}
	}

	return model.Response{
		Success: true,
		Data:    []model.Content{post},
	}
}

// Post handles POST /api/v1/posts
// @Summary Create a new post
// @Description Create a new post
// @Accept  json
// @Produce  json
// @Param 	Authorization	header	string	true	"Access token with the prefix `Bearer `"
// @Param   post    	body  	model.PostRequest      true        "Post"
// @Success 200 		{object} model.Response	"success"
// @Failure 400 		{object} model.Response "Bad request"
// @Security     ApiKeyAuth
// @Security     OAuth2Implicit[admin, write]
// @Router /api/v1/posts [post]
func (c *PostController) Post(req model.PostRequest) model.Response {
	var post = model.Content{
		Type:         model.ContentTypePost,
		Title:        req.Title,
		Text:         req.Text,
		Abstract:     req.Abstract,
		AuthorId:     req.AuthorID,
		IsPublic:     req.IsPublic,
		PublishAt:    req.PublishAt,
		Status:       req.Status,
		AllowComment: req.AllowComment,
	}

	err := c.Service.InsertPost(post)
	if err != nil {
		c.Ctx.Application().Logger().Info(err.Error())
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.Response{Success: false, Message: err.Error()}
	}

	c.Ctx.StatusCode(iris.StatusOK)
	return model.Response{Success: true}
}

func (c *PostController) Put(req model.PostRequest) model.Response {
	var post = model.Content{
		ID:           req.ID,
		Type:         model.ContentTypePost,
		Title:        req.Title,
		Text:         req.Text,
		Abstract:     req.Abstract,
		AuthorId:     req.AuthorID,
		IsPublic:     req.IsPublic,
		PublishAt:    req.PublishAt,
		Status:       req.Status,
		AllowComment: req.AllowComment,
	}

	err := c.Service.UpdatePost(post)
	if err != nil {
		c.Ctx.Application().Logger().Info(err.Error())
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.Response{Success: false, Message: err.Error()}
	}

	c.Ctx.StatusCode(iris.StatusOK)

	return model.Response{
		Success: true,
	}
}

// DeleteBy handles DELETE /api/v1/posts/{id1,id2,id3...:string}
func (c *PostController) DeleteBy(idsReq string) model.Response {
	// check request parameters
	if len(idsReq) == 0 {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.Response{Success: true, Message: dict.ErrInvalidParameters.Error()}
	}

	// convert request parameters from string("1,2,3") to array([1,2,3])
	var ids []uint
	idsString := strings.Split(idsReq, ",")
	for _, item := range idsString {
		id, err := strconv.Atoi(item)
		if err != nil {
			c.Ctx.StatusCode(iris.StatusBadRequest)
			return model.Response{Success: false, Message: err.Error()}
		}
		ids = append(ids, uint(id))
	}

	err := c.Service.DeletePostByIDs(ids)
	if err != nil {
		c.Ctx.StatusCode(400)
		return model.Response{Message: err.Error()}
	}

	c.Ctx.StatusCode(200)
	return model.Response{Success: true}
}

// PostInit handles POST /api/v1/posts/init
func (c *PostController) PostInit() {
	var json map[string]int
	c.Ctx.ReadJSON(&json)

	count := json["count"]
	if count <= 0 {
		c.Ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{"message": dict.ErrInvalidParameters})
	}

	for i := 0; i < count; i++ {

		var tmp = model.Content{
			Type:     model.ContentTypePost,
			Title:    RandStringRunes(10),
			Abstract: RandStringRunes(140),
			Text:     RandStringRunes(500),
			Author: model.User{
				ID: 1,
			},
			IsPublic:     true,
			AllowComment: true,
			ViewsNum:     0,
			CommentsNum:  0,
		}
		err := c.Service.InsertPost(tmp)
		if err != nil {
			c.Ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{"message": err.Error()})
		}

	}

	c.Ctx.StatusCode(200)
	c.Ctx.JSON(iris.Map{"message": "Success to init posts"})
}
