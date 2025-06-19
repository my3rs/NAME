package service

import (
	"NAME/database"
	"NAME/model"

	"gorm.io/gorm"
	"strconv"
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

func NewCommentService() CommentService {
	db := database.GetDB()

	return &commentService{DB: db}
}

func (s *commentService) InsertComment(comment model.Comment) error {
	// 如果有父评论，构建新评论的路径
	if comment.ParentID != 0 {
		var parentComment model.Comment
		if err := s.DB.First(&parentComment, comment.ParentID).Error; err != nil {
			return err
		}
		// 父评论的路径加上当前评论的ID
		comment.Path = parentComment.Path
	}

	// 创建评论
	if err := s.DB.Create(&comment).Error; err != nil {
		return err
	}

	// 更新评论的路径
	// 如果是根评论，路径就是自己的ID
	// 如果是子评论，路径是父评论的路径加上自己的ID
	if comment.ParentID == 0 {
		comment.Path = strconv.FormatUint(uint64(comment.ID), 10)
	} else {
		comment.Path = comment.Path + "." + strconv.FormatUint(uint64(comment.ID), 10)
	}

	// 更新路径
	return s.DB.Model(&comment).Update("path", comment.Path).Error
}

func (s *commentService) GetCommentsByContentID(contentID int, pageIndex int, pageSize int, order string) ([]model.Comment, error) {
	var comments []model.Comment
	offset := pageIndex * pageSize

	// 先获取根评论
	if err := s.DB.Where("content_id = ? AND parent_id = 0", contentID).
		Offset(offset).
		Limit(pageSize).
		Order(order).
		Find(&comments).Error; err != nil {
		return nil, err
	}

	// 对于每个根评论，获取其所有子评论
	for i := range comments {
		var children []model.Comment
		if err := s.DB.Where("content_id = ? AND path LIKE ?", contentID, comments[i].Path+"_%").
			Order(order).
			Find(&children).Error; err != nil {
			return nil, err
		}
		comments[i].Children = children
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
	// 如果没有ID要删除，直接返回成功
	if len(ids) == 0 {
		return 0, nil
	}
	
	result := s.DB.Delete(&model.Comment{}, ids)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
