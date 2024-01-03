package service

import (
	"NAME/database"
	"NAME/dict"
	"NAME/model"
	"errors"
	"gorm.io/gorm"
)

type CommentService interface {
	InsertComment(comment model.Comment) error
	GetCommentsByContentID(contentID int, pageIndex int, pageSize int, order string) []model.Comment
	GetComments(pageIndex int, pageSize int, order string) []model.Comment
	GetCommentByID(id int) model.Comment
	UpdateComment(comment model.Comment) error
}

type commentService struct {
	DB *gorm.DB
}

const queryCommentByContentID = `SELECT *
FROM comments 
WHERE path <@ ARRAY(
	(SELECT id::VARCHAR FROM comments WHERE content_id = ? AND status = ? LIMIT ? OFFSET ?)
)::ltree[];
`

func NewCommentService() CommentService {
	db := database.GetDB()

	return &commentService{DB: db}
}

func (s *commentService) InsertComment(comment model.Comment) error {
	if result := s.DB.Create(&comment); result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *commentService) GetCommentsByContentID(contentID int, pageIndex int, pageSize int, order string) []model.Comment {
	var results []model.Comment
	// log.Print(s.DB.Raw(queryCommentByContentID, contentID, model.CommentStatus_Approved, pageSize, pageIndex*pageSize).Statement.SQL.String(), s.DB.Raw(queryCommentByContentID, contentID, model.CommentStatus_Approved, pageSize, pageIndex*pageSize).Statement.Vars)

	s.DB.Raw(queryCommentByContentID, contentID, model.CommentStatus_Approved, pageSize, pageIndex*pageSize).Scan(&results)

	return results
}

func (s *commentService) GetComments(pageIndex int, pageSize int, order string) []model.Comment {
	var results []model.Comment

	s.DB.Model(&model.Comment{}).Offset(pageIndex * pageSize).Limit(pageSize).Order(order).Find(&results)

	return results
}

func (s *commentService) UpdateComment(comment model.Comment) error {
	result := s.DB.First(&model.Comment{}, comment.ID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return dict.ErrContentNotExists
	}

	s.DB.Save(comment)

	return nil
}

func (s *commentService) GetCommentByID(id int) model.Comment {
	var comment model.Comment
	if result := s.DB.Model(&model.Comment{}).Where("id = ?", id).Take(&comment); result.Error != nil {
		return model.Comment{}
	}
	return comment
}
