package service

import (
	"NAME/database"
	"NAME/dict"
	"NAME/model"
	"gorm.io/gorm"
)

const sqlGetTagsWithPath = `SELECT tag1.id,tag1.no, array_to_string(array_agg(tag2.text ORDER BY tag2.path), ' / ') As text
FROM public.tags As tag1 
INNER JOIN public.tags As tag2 ON (tag2.path @> tag1.path)
GROUP BY tag1.id, tag1.path, tag1.text
ORDER BY text;
`

type TagService interface {
	GetTagByID(id uint) (model.Tag, error)
	GetAllTags() []model.Tag
	GetAllTagsWithPath() []model.Tag
	GetTagsWithOrder(pageIndex int, pageSize int, order string) ([]model.Tag, error)
	GetTagsCount() int64
}

type tagService struct {
	DB *gorm.DB
}

func NewTagService() TagService {
	db, err := database.GetDb()
	if err != nil {
		panic(err.Error())
	}

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

func (s *tagService) GetAllTagsWithPath() []model.Tag {
	var tags []model.Tag
	s.DB.Raw(sqlGetTagsWithPath).Scan(&tags)

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

//func (s *tagService) GetTagsForest() []model.TagTreeNode {
//	var trees []model.TagTreeNode
//
//	for i := 1; i <= 100; i++ {
//		results := s.DB.Where("path ~ ")
//	}
//
//}
