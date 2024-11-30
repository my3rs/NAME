package service

import (
	"NAME/customerror"
	"NAME/database"
	"NAME/dict"
	"NAME/model"
	"errors"
	"fmt"
	"html/template"

	"github.com/kataras/iris/v12"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
	"gorm.io/gorm"
)

type ContentService interface {
	InsertPost(post model.Content) error
	DeletePostByIDs(ids []uint) error
	DeletePostByID(id uint) error
	UpdatePost(content model.Content) error

	GetPostsWithOrder(pageIndex int, pageSize int, order string) []model.Content
	GetPostsCount() int64
	GetPostByID(id int) model.Content
	GetFormattedPostByID(id int) model.Content

	GetContentByID(id int) (model.Content, error)
	GetPureContentByID(id int) model.Content

	GetPageCount() int64

	GetMeta() iris.Map
}

type contentService struct {
	Db *gorm.DB
}

func NewContentService() ContentService {
	db := database.GetDB()

	return &contentService{Db: db}
}

func Markdown2Html(markdown string) template.HTML {
	unsafe := blackfriday.Run([]byte(markdown))
	html := template.HTML(bluemonday.UGCPolicy().SanitizeBytes(unsafe))

	return html
}

func (s *contentService) InsertPost(post model.Content) error {
	if post.Type != model.ContentTypePost {
		return dict.ErrWrongContentType
	}

	result := s.Db.Create(&post)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// 抽取删除关联数据的逻辑到独立函数
func deleteAssociatedData(tx *gorm.DB, ids []uint) error {
	// 删除标签关联
	if err := tx.Table("content_tags").
		Where("content_id IN ?", ids).
		Delete(&struct{}{}).Error; err != nil {
		return customerror.NewDatabaseError("failed to delete content_tags", err)
	}

	// 删除评论
	if err := tx.Where("content_id IN ?", ids).
		Delete(&model.Comment{}).Error; err != nil {
		return customerror.NewDatabaseError("failed to delete comments", err)
	}

	return nil
}

// 删除文章
func (s *contentService) DeletePostByIDs(ids []uint) error {
	if len(ids) == 0 {
		return customerror.NewValidationError("no post ids provided")
	}

	return s.Db.Transaction(func(tx *gorm.DB) error {
		// 检查文章是否都存在
		var count int64
		if err := tx.Model(&model.Content{}).
			Where("id IN ? AND type = ?", ids, model.ContentTypePost).
			Count(&count).Error; err != nil {
			return customerror.NewDatabaseError("failed to check posts existence", err)
		}

		if count != int64(len(ids)) {
			return customerror.NewNotFoundError(
				fmt.Sprintf("some posts do not exist: expected %d, found %d", len(ids), count),
			)
		}

		// 删除关联数据
		if err := deleteAssociatedData(tx, ids); err != nil {
			return err
		}

		// 删除内容
		result := tx.Where("id IN ? AND type = ?", ids, model.ContentTypePost).
			Delete(&model.Content{})
		if result.Error != nil {
			return customerror.NewDatabaseError("failed to delete contents", result.Error)
		}

		if result.RowsAffected == 0 {
			return customerror.NewNotFoundError("no posts were deleted")
		}

		return nil
	})
}

func (s *contentService) DeletePostByID(id uint) error {
	return s.Db.Transaction(func(tx *gorm.DB) error {
		// 检查文章是否都存在
		var count int64
		if err := tx.Model(&model.Content{}).
			Where("id = ? AND type = ?", id, model.ContentTypePost).
			Count(&count).Error; err != nil {
			return customerror.NewDatabaseError("failed to check post existence", err)
		}

		if count != 1 {
			return customerror.NewNotFoundError(
				fmt.Sprintf("post does not exist: expected 1, found %d", count),
			)
		}

		// 删除关联数据
		if err := deleteAssociatedData(tx, []uint{id}); err != nil {
			return err
		}

		// 删除内容
		result := tx.Where("id = ? AND type = ?", id, model.ContentTypePost).
			Delete(&model.Content{})
		if result.Error != nil {
			return customerror.NewDatabaseError("failed to delete contents", result.Error)
		}

		if result.RowsAffected == 0 {
			return customerror.NewNotFoundError("no posts were deleted")
		}

		return nil
	})
}

func (s *contentService) UpdatePost(post model.Content) error {
	if post.Type != model.ContentTypePost {
		return dict.ErrWrongContentType
	}

	result := s.Db.First(&model.Content{}, post.ID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return dict.ErrContentNotExists
	}

	s.Db.Save(&post)

	return nil
}

func (s *contentService) GetPostsWithOrder(pageIndex int, pageSize int, order string) []model.Content {
	var results []model.Content
	s.Db.Preload("Author").Preload("Tags").Preload("Category").Order(order).Limit(pageSize).Offset(pageIndex*pageSize).Find(&results, "type = ?", model.ContentTypePost)

	return results
}

func (s *contentService) GetPostsCount() int64 {
	var count int64
	s.Db.Model(&model.Content{}).Where("type = ?", model.ContentTypePost).Count(&count)

	return count
}

func (s *contentService) GetPostByID(id int) model.Content {
	var post model.Content
	result := s.Db.Model(&model.Content{}).Preload("Author").Where("id = ?", id).Take(&post)

	if result.Error != nil || post.Type != model.ContentTypePost {
		return model.Content{}
	}

	return post
}

func (s *contentService) GetPureContentByID(id int) model.Content {
	var content model.Content
	result := s.Db.First(&content, id)
	if result.Error != nil {
		return model.Content{}
	}
	return content
}

func (s *contentService) GetFormattedPostByID(id int) model.Content {
	post := s.GetPostByID(id)
	if post.ID > 0 {
		post.TextHTML = Markdown2Html(post.Text)
	}

	return post
}

func (s *contentService) GetContentByID(id int) (model.Content, error) {
	var content model.Content
	result := s.Db.Preload("Author").Preload("Tags").First(&content, id)
	if result.Error != nil {
		return model.Content{}, result.Error
	}
	return content, nil
}

func (s *contentService) GetPageCount() int64 {
	var count int64
	s.Db.Model(&model.Content{}).Where("type = ?", model.ContentTypePage).Count(&count)

	return count
}

func (s *contentService) GetMeta() iris.Map {
	var postCount, pageCount, commentCount, tagCount, categoryCount int64

	s.Db.Model(&model.Content{}).Where("type = ?", model.ContentTypePost).Count(&postCount)
	s.Db.Model(&model.Content{}).Where("type = ?", model.ContentTypePage).Count(&pageCount)
	s.Db.Model(&model.Comment{}).Count(&commentCount)
	s.Db.Model(&model.Tag{}).Count(&tagCount)
	s.Db.Model(&model.Category{}).Count(&categoryCount)

	return iris.Map{
		"posts_count":      postCount,
		"pages_count":      pageCount,
		"comments_count":   commentCount,
		"tags_count":       tagCount,
		"categories_count": categoryCount,
	}
}
