package service

import (
	"NAME/database"
	"NAME/dict"
	"NAME/model"
	"errors"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
	"html/template"

	"gorm.io/gorm"
)

type ContentService interface {
	InsertPost(post model.Content) error
	DeletePostByIDs(ids []uint) error
	UpdatePost(content model.Content) error

	GetPostsWithOrder(pageIndex int, pageSize int, order string) []model.Content
	GetPostsCount() int64
	GetPostByID(id int) model.Content
	GetFormattedPostByID(id int) model.Content

	GetContentByID(id int) (model.Content, error)
	GetPureContentByID(id int) model.Content

	GetPageCount() int64
}

type contentService struct {
	Db *gorm.DB
}

func NewContentService() ContentService {
	db, err := database.GetDb()
	if err != nil {
		panic(err.Error())
	}
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

	// result := s.Db.First(&model.Content{}, post.ID)
	// if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
	// 	return dict.ErrContentAlreadyExists
	// }

	result := s.Db.Create(&post)

	if result.Error != nil {
		return result.Error
	}
	//s.Db.Save(&post)

	return nil
}

func (s *contentService) DeletePostByIDs(ids []uint) error {
	// result := s.Db.Select(clause.Associations).Delete(&model.Content{}, ids)
	// result := s.Db.Select("Tags").Delete(&model.Content{}, ids)
	// s.Db.Model(&model.Content{}, ids).Association("Tags").Delete(&model.Tag{}, ids)
	// result := s.Db.Where("type = ?", model.ContentTypePost).Delete(&model.Content{}, ids)
	// result := s.Db.Model(&model.Content{}).Select("Tags").Delete(&model.Content{}, ids)
	// if result.Error != nil {
	// 	return result.Error
	// }

	var contents []model.Content
	s.Db.Where("type = ?", model.ContentTypePost).Find(&contents, ids)
	for _, content := range contents {
		s.Db.Model(&content).Association("Tags").Clear()
	}
	result := s.Db.Delete(&contents)
	if result.Error != nil {
		return result.Error
	}

	return nil
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
	s.Db.Preload("Author").Preload("Tags").Order(order).Limit(pageSize).Offset(pageIndex*pageSize).Find(&results, "type = ?", model.ContentTypePost)

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
