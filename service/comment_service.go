package service

import (
	"NAME/database"
	"NAME/dict"
	"NAME/model"
	"errors"
	"log"

	"gorm.io/gorm"
)

type CommentService interface {
	InsertCommnet(comment model.Comment) error
	GetComments(contentID int, pageIndex int, pageSize int, order string) []model.Comment
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
	db, err := database.GetDb()
	if err != nil {
		log.Panic(err.Error())
	}

	return &commentService{DB: db}
}

func (s *commentService) InsertCommnet(comment model.Comment) error {
	if result := s.DB.Create(&comment); result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *commentService) GetComments(contentID int, pageIndex int, pageSize int, order string) []model.Comment {
	var results []model.Comment
	// log.Print(s.DB.Raw(queryCommentByContentID, contentID, model.CommentStatus_Approved, pageSize, pageIndex*pageSize).Statement.SQL.String(), s.DB.Raw(queryCommentByContentID, contentID, model.CommentStatus_Approved, pageSize, pageIndex*pageSize).Statement.Vars)

	s.DB.Raw(queryCommentByContentID, contentID, model.CommentStatus_Approved, pageSize, pageIndex*pageSize).Scan(&results)

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
