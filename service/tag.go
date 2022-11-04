package service

import (
	"NAME/database"
	"NAME/dict"
	"NAME/model"
	"gorm.io/gorm"
	"log"
)

type TagService interface {
	GetTagByID(id uint) (model.Tag, error)
	GetAllTags() ([]model.Tag, error)
	GetTagsWithOrder(pageIndex int, pageSize int, order string) ([]model.Tag, error)
	GetTagsCount() int64
}

type tagService struct {
	DB *gorm.DB
}

func NewTagService() TagService {
	db, err := database.GetDb()
	if err != nil {
		log.Panic(err.Error())
	}
	return &tagService{DB: db}
}

func (s *tagService) GetTagByID(id uint) (model.Tag, error) {
	var tag model.Tag
	if result := s.DB.First(&tag, id); result.Error != nil {
		return model.Tag{}, result.Error
	}

	return tag, nil
}

func (s *tagService) GetAllTags() ([]model.Tag, error) {
	var tags []model.Tag
	if result := s.DB.Order("path").Find(&tags); result.Error != nil {
		return []model.Tag{}, result.Error
	}

	return tags, nil
}

func (s *tagService) GetTagsWithOrder(pageIndex int, pageSize int, order string) ([]model.Tag, error) {
	if pageSize <= 0 || pageSize >= 100 || pageIndex < 0 {
		return nil, dict.ErrInvalidParameters
	}

	var results []model.Tag
	s.DB.Order(order).Limit(pageSize).Offset(pageIndex * pageSize).Find(&results)

	return results, nil
}

func (s *tagService) GetTagsCount() int64 {
	var count int64
	s.DB.Model(&model.Tag{}).Count(&count)

	return count
}

//func (s *tagService) GetTagsForest() []model.TagTreeNode {
//	var trees []model.TagTreeNode
//
//	for i := 1; i <= 100; i++ {
//		results := s.DB.Where("path ~ ")
//	}
//
//}
