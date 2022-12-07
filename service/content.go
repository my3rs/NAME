package service

import (
	"NAME/database"
	"NAME/dict"
	"NAME/model"
	"errors"
	"log"

	"gorm.io/gorm"
)

type ContentService interface {
	InsertPost(post model.Content) error
	DeletePostByIDs(ids []uint) error
	UpdatePost(content model.Content) error
	GetPostsWithOrder(pageIndex int, pageSize int, order string) []model.Content
	GetPostsCount() int64

	GetContentByID(id int) (model.Content, error)

	GetPageCount() int64
}

type contentService struct {
	Db *gorm.DB
}

func NewContentService() ContentService {
	db, err := database.GetDb()
	if err != nil {
		log.Panic(err.Error())
	}
	return &contentService{Db: db}
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
	result := s.Db.Where("type = ?", model.ContentTypePost).Delete(&model.Content{}, ids)
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
