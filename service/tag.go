package service

import (
	"NAME/database"
	"NAME/dict"
	"NAME/model"
	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

type TagService interface {
	// 查
	GetTagByID(id uint) (model.Tag, error)
	GetAllTags() []model.Tag
	GetTagsWithOrder(pageIndex int, pageSize int, order string) ([]model.Tag, error)
	GetTagsCount() int64

	// GetMetadata 查看标签的统计信息
	GetMetadata() iris.Map

	// 增
	InsertTag(tag model.Tag) error

	// 删
	DeleteTags(ids []uint) error
}

type tagService struct {
	DB *gorm.DB
}

func NewTagService() TagService {
	db := database.GetDB()

	service := &tagService{DB: db}

	return service
}

func (s *tagService) GetTagByID(id uint) (model.Tag, error) {
	var tag model.Tag
	if result := s.DB.First(&tag, id); result.Error != nil {
		return model.Tag{}, result.Error
	}

	return tag, nil
}

func (s *tagService) GetAllTags() []model.Tag {
	var tags []model.Tag
	if result := s.DB.Order("path").Find(&tags); result.Error != nil {
		return []model.Tag{}
	}

	return tags
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

func (s *tagService) GetMetadata() iris.Map {
	var l1cnt, l2cnt, l3cnt, l4cnt int64
	s.DB.Model(&model.Tag{}).Where("nlevel(path) = 1").Count(&l1cnt)
	s.DB.Model(&model.Tag{}).Where("nlevel(path) = 2").Count(&l2cnt)
	s.DB.Model(&model.Tag{}).Where("nlevel(path) = 3").Count(&l3cnt)
	s.DB.Model(&model.Tag{}).Where("nlevel(path) = 4").Count(&l4cnt)

	return iris.Map{
		"sum":         l1cnt + l2cnt + l3cnt + l4cnt,
		"level1Count": l1cnt,
		"level2Count": l2cnt,
		"level3Count": l3cnt,
		"level4Count": l4cnt,
	}
}

func (s *tagService) InsertTag(tag model.Tag) error {
	result := s.DB.Create(&tag)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *tagService) DeleteTags(ids []uint) error {
	// 如果没有ID要删除，直接返回成功
	if len(ids) == 0 {
		return nil
	}
	
	result := s.DB.Delete(&model.Tag{}, ids)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
