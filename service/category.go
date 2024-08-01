package service

import (
	"NAME/database"
	"NAME/model"
	"gorm.io/gorm"
)

type CategoryService interface {
	GetCategoriesCount() int64
	GetCategories(pageIndex int, pageSize int, order string) []model.Category
	InsertCategory(category model.Category) error
	UpdateCategory(category model.Category) error
	DeleteCategories(ids []uint) error
}

type categoryService struct {
	DB *gorm.DB
}

func NewCategoryService() CategoryService {
	db := database.GetDB()

	return &categoryService{DB: db}
}

func (s *categoryService) GetCategories(pageIndex int, pageSize int, order string) []model.Category {
	var results []model.Category
	s.DB.Model(&model.Category{}).Offset(pageIndex * pageSize).Limit(pageSize).Order(order).Find(&results)

	return results
}

func (s *categoryService) GetCategoriesCount() int64 {
	var result int64
	s.DB.Model(&model.Category{}).Count(&result)

	return result
}

func (s *categoryService) InsertCategory(category model.Category) error {
	result := s.DB.Create(&category)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *categoryService) UpdateCategory(category model.Category) error {
	result := s.DB.Save(&category)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *categoryService) DeleteCategories(ids []uint) error {
	result := s.DB.Delete(&model.Category{}, ids)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
