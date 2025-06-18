package controller

import (
	"NAME/dict"
	"NAME/model"
	"NAME/service"
	"reflect"
	"strconv"
	"strings"

	"github.com/kataras/iris/v12"
)

type CommentController struct {
	Ctx            iris.Context
	CommentService service.CommentService
	UserService    service.UserService
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

func (c *CommentController) Post(req model.Comment) model.DetailResponse {
	// 检查评论是否合规
	if err := checkComment(req.Text); err != nil {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.DetailResponse{Success: false, Message: err.Error()}
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
		Status:     model.CommentStatusApproved,
	}

	// 如果是已登录用户，直接从数据库中获取用户信息
	if comment.AuthorID != 0 {
		user, err := c.UserService.GetUserByID(int(comment.AuthorID))
		if err != nil {
			c.Ctx.Application().Logger().Error("通过ID获取用户失败：", err.Error())
		}

		comment.AuthorName = user.Username
		comment.Mail = user.Mail
		comment.URL = user.URL

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
		return model.DetailResponse{Success: false, Message: err.Error()}
	}

	c.Ctx.StatusCode(iris.StatusOK)
	return model.DetailResponse{Success: true, Message: "评论发表成功", Data: comment}
}

func (c *CommentController) Get(req getCommentsRequest) model.PageResponse {
	// 页码参数检查
	if req.PageIndex <= 0 || req.PageSize <= 0 {
		c.Ctx.Application().Logger().Error("Failed to get comments: pageIndex or pageSize <= 0")
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.NewPageResponse(false, dict.ErrInvalidParameters.Error()+": pageSize or pageIndex", nil, req.PageIndex, req.PageSize, 0)
	}

	var data []model.Comment
	var totalCount int64
	var err error

	// 指定内容ID时，获取该内容下的评论
	if req.ContentID != 0 {
		data, err = c.CommentService.GetCommentsByContentID(int(req.ContentID), req.PageIndex-1, req.PageSize, req.Order)
		if err != nil {
			c.Ctx.StatusCode(iris.StatusBadRequest)
			return model.NewPageResponse(false, err.Error(), nil, req.PageIndex, req.PageSize, 0)
		}
		totalCount = c.CommentService.GetCommentsCount(int64(req.ContentID))
	} else {
		// 获取所有评论
		switch req.WithContent {
		case "title":
			data, err = c.CommentService.GetCommentsWithContentTitle(req.PageIndex-1, req.PageSize, req.Order)
		default:
			data, err = c.CommentService.GetComments(req.PageIndex-1, req.PageSize, req.Order)
		}
		if err != nil {
			c.Ctx.StatusCode(iris.StatusBadRequest)
			return model.NewPageResponse(false, err.Error(), nil, req.PageIndex, req.PageSize, 0)
		}
		totalCount = c.CommentService.GetCommentsCount(0)
	}

	return model.NewPageResponse(true, "获取评论成功", data, req.PageIndex, req.PageSize, totalCount)
}

// PutBy Put handles PUT /api/v1/comments/{:id} 更新评论（完整）
func (c *CommentController) PutBy(id uint) model.DetailResponse {
	var req model.Comment
	if err := c.Ctx.ReadJSON(&req); err != nil {
		c.Ctx.Application().Logger().Error("Failed to read json from PUT comments request: ", err.Error())
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.DetailResponse{Success: false, Message: err.Error()}
	}

	comment, err := c.CommentService.GetCommentByID(int(id))
	if err != nil {
		c.Ctx.Application().Logger().Error("Failed to get comment by id: ", err.Error())
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.DetailResponse{Success: false, Message: err.Error()}
	}

	if err := c.CommentService.UpdateComment(comment); err != nil {
		c.Ctx.Application().Logger().Error("Failed to update comment: ", err.Error())
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.DetailResponse{Success: false, Message: err.Error()}
	}

	c.Ctx.StatusCode(iris.StatusOK)
	return model.DetailResponse{Success: true, Message: "评论更新成功", Data: comment}
}

// PatchBy handles PATCH /api/v1/comments/{:id} 更新评论（指定字段）
func (c *CommentController) PatchBy(id uint) model.DetailResponse {
	var req model.Comment
	if err := c.Ctx.ReadJSON(&req); err != nil {
		c.Ctx.Application().Logger().Error("Failed to read json from PATCH comments request: ", err.Error())
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.DetailResponse{Success: false, Message: err.Error()}
	}

	comment, err := c.CommentService.GetCommentByID(int(id))
	if err != nil {
		c.Ctx.Application().Logger().Error("Failed to get comment by id: ", err.Error())
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.DetailResponse{Success: false, Message: err.Error()}
	}

	// 将非空字段复制到目标对象
	replaceNonEmptyFields(&model.Comment{
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
	}, &comment)

	if err := c.CommentService.UpdateComment(comment); err != nil {
		c.Ctx.Application().Logger().Error("Failed to update comment: ", err.Error())
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.DetailResponse{Success: false, Message: err.Error()}
	}

	c.Ctx.StatusCode(iris.StatusOK)
	return model.DetailResponse{Success: true, Message: "评论更新成功", Data: comment}
}

// DeleteBy handles DELETE /api/v1/comments/{:ids_string} 删除评论
func (c *CommentController) DeleteBy(idsReq string) model.BatchResponse {
	if len(idsReq) == 0 {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return model.BatchResponse{Success: false, Message: "bad params in url"}
	}

	// 去掉末尾的逗号
	if idsReq[len(idsReq)-1] == ',' {
		idsReq = idsReq[:len(idsReq)-1]
	}

	// 将字符中的ID转换为数组
	var ids []uint
	str := strings.Split(idsReq, ",")
	for _, item := range str {
		if len(item) == 0 {
			continue
		}
		id, err := strconv.Atoi(item)
		if err != nil {
			c.Ctx.StatusCode(iris.StatusBadRequest)
			return model.BatchResponse{Success: false, Message: err.Error()}
		}
		ids = append(ids, uint(id))
	}

	// 调用 Service 进行删除
	var successList []uint
	var failedList []uint

	for _, id := range ids {
		err := c.CommentService.DeleteComment(id)
		if err != nil {
			failedList = append(failedList, id)
		} else {
			successList = append(successList, id)
		}
	}

	if len(failedList) > 0 {
		return model.NewBatchResponse(false, "部分删除成功", successList, failedList)
	}

	return model.NewBatchResponse(true, "success", successList, nil)
}
