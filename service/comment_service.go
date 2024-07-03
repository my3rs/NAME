package service

import (
	"NAME/database"
	"NAME/dict"
	"NAME/model"
	"errors"
	"gorm.io/gorm"
)

type CommentService interface {
	// 增加评论
	InsertComment(comment model.Comment) error

	// 删
	DeleteComment(id uint) error
	DeleteBatchComment(ids []uint) (int64, error)

	// 查
	GetComments(pageIndex int, pageSize int, order string) []model.Comment
	GetCommentsWithContentTitle(pageIndex int, pageSize int, order string) []model.Comment
	GetCommentsCount(contentID int64) int64
	GetCommentByID(id int) model.Comment
	GetCommentsByContentID(contentID int, pageIndex int, pageSize int, order string) []model.Comment

	// 改
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

func (s *commentService) GetCommentsWithContentTitle(pageIndex int, pageSize int, order string) []model.Comment {
	var results []model.Comment

	s.DB.Model(&model.Comment{}).Offset(pageIndex * pageSize).Limit(pageSize).Order(order).Preload("Content").Find(&results)

	return results
}

// GetCommentsCount 获取评论总数
// - contentID 为 0 时，获取所有评论总数
// - contentID 不为 0 时，获取指定内容的评论总数
func (s *commentService) GetCommentsCount(contentID int64) int64 {
	var count int64

	if contentID != 0 {
		s.DB.Model(&model.Comment{}).Where("content_id = ?", contentID).Count(&count)
	} else {
		s.DB.Model(&model.Comment{}).Count(&count)
	}

	return count
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

func (s *commentService) DeleteComment(id uint) error {
	result := s.DB.Delete(&model.Comment{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *commentService) DeleteBatchComment(ids []uint) (int64, error) {
	result := s.DB.Delete(&model.Comment{}, ids)
	if result.Error != nil {
		return result.RowsAffected, result.Error
	}

	return result.RowsAffected, nil
}
