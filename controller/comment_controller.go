package controller

import (
	"NAME/dict"
	"NAME/model"
	"NAME/service"
	"reflect"

	"github.com/kataras/iris/v12"
)

type CommentController struct {
	Ctx            iris.Context
	CommentService service.CommentService
	UserService    service.UserService
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

type PostCommentResponse struct {
	Success bool            `json:"success,omitempty"`
	Message string          `json:"message,omitempty"`
	Data    []model.Comment `json:"data,omitempty"`
	Page    *model.Page     `json:"pagination,omitempty"`
}

type getCommentsRequest struct {
	PageSize    int    `url:"pageSize" json:"pageSize"`
	PageIndex   int    `url:"pageIndex" json:"pageIndex"`
	Order       string `url:"orderBy" json:"orderBy,omitempty"`
	ContentID   uint   `url:"contentID" json:"contentID"`
	WithContent string `url:"withContent" json:"withContent"`
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

func checkComment(comment string) error {
	if comment == "" || len(comment) == 0 {
		return dict.ErrEmptyContent
	}

	return nil
}

func (c *CommentController) Post(req postCommentRequest) PostCommentResponse {
	// 检查评论是否合规
	if checkComment(req.Text) != nil {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return PostCommentResponse{Success: false, Message: dict.ErrEmptyContent.Error()}
	}

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

	// 如果是已登录用户，直接从数据库中获取用户信息
	if comment.AuthorID != 0 {

		user, err := c.UserService.GetUserByID(int(comment.AuthorID))
		if err != nil {
			c.Ctx.Application().Logger().Error("通过ID获取用户失败：", err.Error())
		}

		comment.AuthorName = user.Name
		comment.Mail = user.Mail
		comment.URL = user.Url

		c.Ctx.Application().Logger().Info("评论--用户ID：", comment.AuthorID, " 用户名：", comment.AuthorName)
	}

	// 获取用户IP和User-Agent
	c.Ctx.Application().Logger().Info(c.Ctx.GetHeader("User-Agent"))
	comment.Agent = c.Ctx.GetHeader("User-Agent")
	comment.IP = c.Ctx.RemoteAddr()

	c.Ctx.Application().Logger().Infof("发表评论：%+v", comment)

	if err := c.CommentService.InsertComment(comment); err != nil {
		c.Ctx.Application().Logger().Error(err.Error())
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return PostCommentResponse{Success: false, Message: err.Error()}
	}

	c.Ctx.StatusCode(iris.StatusOK)
	return PostCommentResponse{Success: true}
}

func (c *CommentController) Get(req getCommentsRequest) {
	// 页码参数检查
	if req.PageIndex <= 0 || req.PageSize <= 0 {
		c.Ctx.Application().Logger().Error("Failed to get comments: pageIndex or pageSize <= 0")
		c.Ctx.StatusCode(iris.StatusBadRequest)
		c.Ctx.JSON(iris.Map{
			"success": false,
			"message": dict.ErrInvalidParameters.Error() + ": pageSize or pageIndex",
		})

		return
	}

	// 指定内容ID时，获取该内容下的评论
	if req.ContentID != 0 {
		data := c.CommentService.GetCommentsByContentID(int(req.ContentID), req.PageIndex-1, req.PageSize, req.Order)
		var page model.Page
		page.PageIndex = req.PageIndex
		page.PageSize = len(data)
		page.Order = req.Order
		page.ContentCount = c.CommentService.GetCommentsCount(int64(req.ContentID))

		c.Ctx.StatusCode(iris.StatusOK)
		c.Ctx.JSON(iris.Map{
			"success":    true,
			"data":       data,
			"pagination": &page,
		})

		return

	}

	// 获取所有评论

	var data []model.Comment
	switch req.WithContent {
	case "title":
		data = c.CommentService.GetCommentsWithContentTitle(req.PageIndex-1, req.PageSize, req.Order)
		break
	default:

		data = c.CommentService.GetComments(req.PageIndex-1, req.PageSize, req.Order)
		break
	}

	var page model.Page

	page.PageIndex = req.PageIndex
	page.PageSize = len(data)
	page.Order = req.Order
	page.ContentCount = c.CommentService.GetCommentsCount(0)

	c.Ctx.StatusCode(iris.StatusOK)
	c.Ctx.JSON(iris.Map{
		"success":    true,
		"data":       data,
		"pagination": &page,
	})

	return
}

// PutBy Put handles PUT /api/v1/comments/{:id} 更新评论（完整）
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
	if err := c.CommentService.UpdateComment(comment); err != nil {
		return model.Response{Success: false, Message: err.Error()}
	}

	return model.Response{Success: true, Message: "Succeed to update comment"}
}

// PatchBy handles PATCH /api/v1/comments/{:id} 更新评论（指定字段）
func (c *CommentController) PatchBy(id uint) model.Response {
	var req model.Comment
	if err := c.Ctx.ReadJSON(&req); err != nil {
		c.Ctx.Application().Logger().Error("Failed to read json from PUT comments request: ", err.Error())
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.Response{Success: false, Message: err.Error()}
	}

	comment := c.CommentService.GetCommentByID(int(id))
	if comment.ID == 0 {
		c.Ctx.Application().Logger().Error("Comment doesn't exist, id = ", id)
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.Response{Success: false, Message: "Comment doesn't exist"}
	}

	replaceNonEmptyFields(&req, &comment)
	c.CommentService.UpdateComment(comment)

	return model.Response{Success: true, Message: "Succeed to update comment"}
}

// DeleteBy handles DELETE /api/v1/comments/{:id} 删除评论
func (c *CommentController) DeleteBy(id uint) model.Response {
	if err := c.CommentService.DeleteComment(id); err != nil {
		return model.Response{Success: false, Message: err.Error()}
	}

	return model.Response{Success: true, Message: "Succeed to delete comment"}
}
