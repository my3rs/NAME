package controller

import (
	"NAME/dict"
	"NAME/model"
	"NAME/service"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"strconv"
	"strings"

	"github.com/kataras/iris/v12/mvc"

	"github.com/kataras/iris/v12"
)

type PostController struct {
	Ctx             iris.Context
	Service         service.ContentService
	TagService      service.TagService
	CategoryService service.CategoryService
}

type postContentRequest struct {
	Title        string              `json:"title"`
	Abstract     string              `json:"abstract"`
	Text         string              `json:"text"`
	Author       model.User          `json:"author"`
	Category     model.Category      `json:"category"`
	Status       model.ContentStatus `json:"status"`
	CreatedAt    int64               `json:"createdAt"`
	UpdatedAt    int64               `json:"updatedAt"`
	PublishAt    int64               `json:"publishAt"`
	AllowComment bool                `json:"allowComment"`
	Tags         []model.Tag         `json:"tags"`
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (c *PostController) BeforeActivation(b mvc.BeforeActivation) {
	log.Print("before ", b.GetRoutes("GET"))
}

func (c *PostController) AfterActivation(b mvc.BeforeActivation) {
	log.Print("after ", b.GetRoutes("GET"))
}

func replaceContentNonEmptyFields(src, dst *model.Content) {
	t := reflect.TypeOf(*src)
	v := reflect.ValueOf(*src)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// 不可覆盖的字段
		if field.Name == "ID" || field.Name == "Type" || field.Name == "AuthorId" ||
			field.Name == "ViewsNum" || field.Name == "Tags" {
			continue
		}

		if value.Interface() != reflect.Zero(value.Type()).Interface() {
			dstField := reflect.ValueOf(dst).Elem().FieldByName(field.Name)
			dstField.Set(value)
		}
	}
}

func (c *PostController) Get(req model.QueryRequest) model.PageResponse {
	// check parameters
	if req.PageSize <= 0 || req.PageIndex < 0 {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.NewPageResponse(false, dict.ErrInvalidParameters.Error()+": pageSize or pageIndex", nil, req.PageIndex, req.PageSize, 0)
	}

	// get total count
	count := c.Service.GetPostsCount()
	if count == 0 {
		return model.NewPageResponse(true, "success", []model.Content{}, req.PageIndex, req.PageSize, 0)
	}

	// check pageIndex
	maxPageIndex := (count - 1) / int64(req.PageSize)
	if int64(req.PageIndex) > maxPageIndex {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.NewPageResponse(false, "pageIndex too large", nil, req.PageIndex, req.PageSize, count)
	}

	// get posts
	posts := c.Service.GetPostsWithOrder(req.PageIndex, req.PageSize, req.Order)

	return model.NewPageResponse(true, "success", posts, req.PageIndex, req.PageSize, count)
}

func (c *PostController) GetBy(id int) model.DetailResponse {
	post, err := c.Service.GetContentByID(id)
	if err != nil {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.DetailResponse{Success: false, Message: err.Error()}
	}

	readablePath, err := strconv.ParseBool(c.Ctx.URLParam("readablePath"))
	if err != nil {
		readablePath = false
	}

	if readablePath {
		debug, found := model.GetSettingsItem("environment")
		if debug.Value == model.EnvironmentDev || !found {
			c.Ctx.Application().Logger().Info("查找标签的可读化路径")
		}

		for i, tag := range post.Tags {
			post.Tags[i].Text = c.TagService.GetTagReadablePath(tag.ID)
		}
	}

	return model.DetailResponse{Success: true, Message: "success", Data: post}
}

func (c *PostController) Post(req postContentRequest) model.EmptyResponse {
	var post = model.Content{
		Type:         model.ContentTypePost,
		Title:        req.Title,
		Text:         req.Text,
		Abstract:     req.Abstract,
		PublishAt:    req.PublishAt,
		Status:       req.Status,
		AllowComment: req.AllowComment,
		Category:     req.Category,
		Author:       req.Author,
		Tags:         req.Tags,
	}

	err := c.Service.InsertPost(post)
	if err != nil {
		c.Ctx.Application().Logger().Info(err.Error())
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.EmptyResponse{Success: false, Message: err.Error()}
	}

	return model.EmptyResponse{Success: true, Message: "success"}
}

func (c *PostController) PatchBy(id uint) model.EmptyResponse {
	var req model.Content
	if err := c.Ctx.ReadJSON(&req); err != nil {
		c.Ctx.Application().Logger().Error("Failed to read json from PUT post request: ", err.Error())
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.EmptyResponse{Success: false, Message: err.Error()}
	}

	post := c.Service.GetPureContentByID(int(id))
	if post.ID == 0 {
		c.Ctx.Application().Logger().Error("Post doesn't exist, id = ", id)
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.EmptyResponse{Success: false, Message: "Post doesn't exist"}
	}

	replaceContentNonEmptyFields(&req, &post)
	c.Service.UpdatePost(post)

	return model.EmptyResponse{Success: true, Message: "success"}
}

func (c *PostController) PutBy(id int) model.EmptyResponse {
	var req model.Content
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		c.Ctx.Application().Logger().Error("Failed to read json from PUT post request: ", err.Error())
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.EmptyResponse{Success: false, Message: err.Error()}
	}

	if id != int(req.ID) {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.EmptyResponse{Success: false, Message: "ID mismatch"}
	}

	post := model.Content{
		ID:           req.ID,
		Type:         model.ContentTypePost,
		Title:        req.Title,
		Text:         req.Text,
		Abstract:     req.Abstract,
		AuthorId:     req.AuthorId,
		PublishAt:    req.PublishAt,
		Status:       req.Status,
		AllowComment: req.AllowComment,
	}

	err = c.Service.UpdatePost(post)
	if err != nil {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.EmptyResponse{Success: false, Message: err.Error()}
	}

	return model.EmptyResponse{Success: true, Message: "success"}
}

// DeleteBy handles DELETE /api/v1/post/1,2,3 批量删除文章
func (c *PostController) DeleteBy(idsReq string) model.BatchResponse {
	if idsReq == "" {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.BatchResponse{Success: false, Message: "ids is empty"}
	}

	ids := strings.Split(idsReq, ",")
	var successList []uint
	var failedList []uint

	for _, idStr := range ids {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			failedList = append(failedList, 0)
			continue
		}

		err = c.Service.DeletePostByID(uint(id))
		if err != nil {
			failedList = append(failedList, uint(id))
		} else {
			successList = append(successList, uint(id))
		}
	}

	if len(successList) == 0 {
		return model.NewBatchResponse(false, "删除失败", nil, failedList)
	}

	if len(failedList) > 0 {
		return model.NewBatchResponse(false, "部分删除成功", successList, failedList)
	}

	return model.NewBatchResponse(true, "success", successList, nil)
}

func (c *PostController) PostInit() model.EmptyResponse {
	for i := 0; i < 10; i++ {
		post := model.Content{
			Type:         model.ContentTypePost,
			Title:        fmt.Sprintf("Post %d", i),
			Text:         fmt.Sprintf("This is post %d", i),
			Abstract:     fmt.Sprintf("Abstract of post %d", i),
			AuthorId:     1,
			Status:       model.ContentStatusPublished,
			AllowComment: true,
		}

		err := c.Service.InsertPost(post)
		if err != nil {
			return model.EmptyResponse{Success: false, Message: err.Error()}
		}
	}

	return model.EmptyResponse{Success: true, Message: "success"}
}
