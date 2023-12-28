package service

import (
	"NAME/database"
	"NAME/dict"
	"NAME/model"
	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

const sqlGetTagsWithPath = `SELECT tag1.id,tag1.no,tag1.parent_id, tag1.path, array_to_string(array_agg(tag2.text ORDER BY tag2.path), ' / ') As text
FROM public.tags As tag1 
INNER JOIN public.tags As tag2 ON (tag2.path @> tag1.path)
GROUP BY tag1.id, tag1.path, tag1.text
ORDER BY text;
`

const sqlGetTagByIDWithReadablePath = `SELECT 
  array_to_string(array_agg(tag2.text ORDER BY tag2.path), ' / ') AS text
FROM 
  public.tags AS tag1 
INNER JOIN 
  public.tags AS tag2 ON (tag2.path @> tag1.path)
WHERE 
  tag1.id = ?
GROUP BY 
  tag1.id, tag1.path, tag1.text;
`

type TagService interface {
	GetTagByID(id uint) (model.Tag, error)
	GetAllTags() []model.Tag
	GetAllTagsWithPath() []model.Tag
	GetTagsWithOrder(pageIndex int, pageSize int, order string) ([]model.Tag, error)
	GetTagsCount() int64

	GetTagByIDWithReadablePath(id uint) (model.Tag, error)
	GetTagReadablePath(id uint) string

	// GetMetadata 查看标签的统计信息
	GetMetadata() iris.Map

	InsertTag(tag model.Tag) error
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

func (s *tagService) GetTagByIDWithReadablePath(id uint) (model.Tag, error) {
	var result model.Tag
	s.DB.Raw(sqlGetTagByIDWithReadablePath, id).Scan(&result)

	return result, nil
}

func (s *tagService) GetTagReadablePath(id uint) string {
	var path string
	s.DB.Raw(sqlGetTagByIDWithReadablePath, id).Scan(&path)
	return path
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
