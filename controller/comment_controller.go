package controller

import (
	"NAME/dict"
	"NAME/model"
	"NAME/service"
	"reflect"

	"github.com/kataras/iris/v12"
)

type CommentController struct {
	Ctx     iris.Context
	Service service.CommentService
}

type postCommentRequest struct {
	ID         uint   `json:"id"`
	ContentID  uint   `json:"contentID"`
	AuthorID   uint   `json:"authorID"`
	AuthorName string `json:"authorName"`
	Mail       string
	URL        string `json:"url"`
	Text       string
	IP         string `json:"ip"`
	Agent      string `json:"agnet"`

	ParentID  uint `json:"parentID"`
	CreatedAt int  `json:"createdAt"`
}

type postCommnetResponse struct {
	Success bool            `json:"success,omitempty"`
	Message string          `json:"message,omitempty"`
	Data    []model.Comment `json:"data,omitempty"`
	Page    *model.Page     `json:"pagination,omitempty"`
}

type getCommentsRequest struct {
	PageSize  int    `url:"pageSize" json:"pageSize"`
	PageIndex int    `url:"pageIndex" json:"pageIndex"`
	Order     string `url:"orderBy" json:"orderBy,omitempty"`
	ContentID uint   `url:"contentID" json:"contentID"`
}

func NewCommentController() *CommentController {
	return &CommentController{Service: service.NewCommentService()}
}

func replaceNonEmptyFields(src, dst *model.Comment) {
	t := reflect.TypeOf(*src)
	v := reflect.ValueOf(*src)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// 不可覆盖的字段
		if field.Name == "Content" || field.Name == "ParentID" || field.Name == "ID" {
			continue
		}

		if value.Interface() != reflect.Zero(value.Type()).Interface() {
			dstField := reflect.ValueOf(dst).Elem().FieldByName(field.Name)
			dstField.Set(value)
		}
	}
}

func (c *CommentController) Post(req postCommentRequest) postCommnetResponse {
	var comment = model.Comment{
		ID:         req.ID,
		ContentID:  req.ContentID,
		AuthorID:   req.AuthorID,
		AuthorName: req.AuthorName,
		Mail:       req.Mail,
		URL:        req.URL,
		Text:       req.Text,
		IP:         req.IP,
		Agent:      req.Agent,
		ParentID:   req.ParentID,
	}

	if err := c.Service.InsertCommnet(comment); err != nil {
		c.Ctx.Application().Logger().Error(err.Error())
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return postCommnetResponse{Success: false, Message: err.Error()}
	}

	c.Ctx.StatusCode(iris.StatusOK)
	return postCommnetResponse{Success: true}
}

func (c *CommentController) Get(req getCommentsRequest) postCommnetResponse {
	if req.PageIndex <= 0 || req.PageSize <= 0 {
		c.Ctx.Application().Logger().Error("Failed to get comments: pageIndex or pageSize <= 0")
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return postCommnetResponse{Success: false, Message: dict.ErrInvalidParameters.Error() + ": pageSize or pageIndex"}
	}

	var rsp postCommnetResponse
	var page model.Page

	rsp.Success = true
	rsp.Data = c.Service.GetComments(int(req.ContentID), req.PageIndex-1, req.PageSize, req.Order)
	page.PageIndex = req.PageIndex
	page.PageSize = len(rsp.Data)
	page.Order = req.Order

	rsp.Page = &page

	return rsp
}

// Put handles PUT /api/v1/comments/{:id} 更新评论（完整）
func (c *CommentController) PutBy(id uint) model.Response {
	var req postCommentRequest
	if err := c.Ctx.ReadJSON(&req); err != nil {
		c.Ctx.Application().Logger().Error("Failed to read json from PUT comments request: ", err.Error())
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.Response{Success: false, Message: err.Error()}
	}

	var comment = model.Comment{
		ID:         id,
		ContentID:  req.ContentID,
		AuthorID:   req.AuthorID,
		AuthorName: req.AuthorName,
		Mail:       req.Mail,
		URL:        req.URL,
		Text:       req.Text,
		IP:         req.IP,
		Agent:      req.Agent,
		ParentID:   req.ParentID,
	}
	if err := c.Service.UpdateComment(comment); err != nil {
		return model.Response{Success: false, Message: err.Error()}
	}

	return model.Response{Success: true, Message: "Succeed to update comment"}
}

// Patch handles PATCH /api/v1/comments/{:id} 更新评论（指定字段）
func (c *CommentController) PatchBy(id uint) model.Response {
	var req model.Comment
	if err := c.Ctx.ReadJSON(&req); err != nil {
		c.Ctx.Application().Logger().Error("Failed to read json from PUT comments request: ", err.Error())
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.Response{Success: false, Message: err.Error()}
	}

	comment := c.Service.GetCommentByID(int(id))
	if comment.ID == 0 {
		c.Ctx.Application().Logger().Error("Comment doesn't exist, id = ", id)
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.Response{Success: false, Message: "Comment doesn't exist"}
	}

	replaceNonEmptyFields(&req, &comment)
	c.Service.UpdateComment(comment)

	return model.Response{Success: true, Message: "Succeed to update comment"}
}
