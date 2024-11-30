package service

import (
	"NAME/database"
	"NAME/model"

	"gorm.io/gorm"
)

type CommentService interface {
	// 增加评论
	InsertComment(comment model.Comment) error

	// 删
	DeleteComment(id uint) error
	DeleteBatchComment(ids []uint) (int64, error)

	// 查
	GetComments(pageIndex int, pageSize int, order string) ([]model.Comment, error)
	GetCommentsWithContentTitle(pageIndex int, pageSize int, order string) ([]model.Comment, error)
	GetCommentsCount(contentID int64) int64
	GetCommentByID(id int) (model.Comment, error)
	GetCommentsByContentID(contentID int, pageIndex int, pageSize int, order string) ([]model.Comment, error)

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

func (s *commentService) GetCommentsByContentID(contentID int, pageIndex int, pageSize int, order string) ([]model.Comment, error) {
	var comments []model.Comment
	if result := s.DB.Where("content_id = ?", contentID).
		Offset(pageIndex * pageSize).
		Limit(pageSize).
		Order(order).
		Find(&comments); result.Error != nil {
		return nil, result.Error
	}
	return comments, nil
}

func (s *commentService) GetComments(pageIndex int, pageSize int, order string) ([]model.Comment, error) {
	var comments []model.Comment
	if result := s.DB.Offset(pageIndex * pageSize).
		Limit(pageSize).
		Order(order).
		Find(&comments); result.Error != nil {
		return nil, result.Error
	}
	return comments, nil
}

func (s *commentService) GetCommentsWithContentTitle(pageIndex int, pageSize int, order string) ([]model.Comment, error) {
	var comments []model.Comment
	if result := s.DB.Joins("Content").
		Offset(pageIndex * pageSize).
		Limit(pageSize).
		Order(order).
		Find(&comments); result.Error != nil {
		return nil, result.Error
	}
	return comments, nil
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
	if result := s.DB.Save(&comment); result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *commentService) GetCommentByID(id int) (model.Comment, error) {
	var comment model.Comment
	if result := s.DB.First(&comment, id); result.Error != nil {
		return model.Comment{}, result.Error
	}
	return comment, nil
}

func (s *commentService) DeleteComment(id uint) error {
	if result := s.DB.Delete(&model.Comment{}, id); result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *commentService) DeleteBatchComment(ids []uint) (int64, error) {
	result := s.DB.Delete(&model.Comment{}, ids)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
