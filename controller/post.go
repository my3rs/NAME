package controller

import (
	"NAME/dict"
	"NAME/model"
	"NAME/service"
	"math/rand"
	"strconv"
	"strings"

	"github.com/kataras/iris/v12"
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
		req.Order = "created_at desc"
	}
	page.Order = req.Order

	posts := c.Service.GetPostsWithOrder(req.PageIndex-1, req.PageSize, req.Order)

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

func (c *PostController) Post(req model.PostRequest) model.Response {
	var post = model.Content{
		Type:         model.ContentTypePost,
		Title:        req.Title,
		Text:         req.Text,
		Abstract:     req.Abstract,
		AuthorId:     req.AuthorID,
		PublishAt:    req.PublishAt,
		Status:       req.Status,
		AllowComment: req.AllowComment,
		Tags:         req.Tags,
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

func (c *PostController) PutBy(id int) model.Response {
	var req model.PostRequest
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		c.Ctx.StatusCode(400)
		return model.Response{Success: false, Message: err.Error()}
	}

	if id != int(req.ID) {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.Response{Success: false, Message: "Check IDs in URL and body"}
	}

	var post = model.Content{
		ID:           req.ID,
		Type:         model.ContentTypePost,
		Title:        req.Title,
		Text:         req.Text,
		Abstract:     req.Abstract,
		AuthorId:     req.AuthorID,
		CreatedAt:    req.CreatedAt,
		UpdatedAt:    req.UpdatedAt,
		PublishAt:    req.PublishAt,
		Status:       req.Status,
		AllowComment: req.AllowComment,
		Tags:         req.Tags,
	}
	err = c.Service.UpdatePost(post)

	if err != nil {
		c.Ctx.Application().Logger().Info(err.Error())
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.Response{Success: false, Message: err.Error()}
	}

	c.Ctx.StatusCode(iris.StatusOK)
	return model.Response{Success: true}
}

// func (c *PostController) Put(req model.PostRequest) model.Response {
// 	var post = model.Content{
// 		ID:           req.ID,
// 		Type:         model.ContentTypePost,
// 		Title:        req.Title,
// 		Text:         req.Text,
// 		Abstract:     req.Abstract,
// 		AuthorId:     req.AuthorID,
// 		IsPublic:     req.IsPublic,
// 		PublishAt:    req.PublishAt,
// 		Status:       req.Status,
// 		AllowComment: req.AllowComment,
// 	}

// 	err := c.Service.UpdatePost(post)
// 	if err != nil {
// 		c.Ctx.Application().Logger().Info(err.Error())
// 		c.Ctx.StatusCode(iris.StatusBadRequest)
// 		return model.Response{Success: false, Message: err.Error()}
// 	}

// 	c.Ctx.StatusCode(iris.StatusOK)

// 	return model.Response{
// 		Success: true,
// 	}
// }

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
		if len(item) == 0 {
			continue
		}
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
func (c *PostController) PostInit() iris.Map {
	var json map[string]int
	c.Ctx.ReadJSON(&json)

	c.Ctx.Application().Logger().Info(json)

	count := json["count"]
	if count <= 0 {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return iris.Map{"message": dict.ErrInvalidParameters}
	}

	c.Ctx.Application().Logger().Info("Init ", count, " posts")

	for i := 0; i < count; i++ {

		var tmp = model.Content{
			Type:     model.ContentTypePost,
			Title:    RandStringRunes(10),
			Abstract: RandStringRunes(70),
			Text:     RandStringRunes(200),
			Author: model.User{
				ID: 1,
			},
			AllowComment: true,
			ViewsNum:     0,
			CommentsNum:  0,
		}
		err := c.Service.InsertPost(tmp)
		if err != nil {
			c.Ctx.StatusCode(iris.StatusBadRequest)
			return iris.Map{"message": err.Error()}
		}

	}

	c.Ctx.StatusCode(200)
	return iris.Map{"message": "Success to init posts"}
}
