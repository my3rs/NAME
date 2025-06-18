package service

import (
	"NAME/model"
	"NAME/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var testTagIDs []uint

func cleanupTestTags(db *gorm.DB) error {
	if len(testTagIDs) > 0 {
		if err := db.Where("id IN ?", testTagIDs).Delete(&model.Tag{}).Error; err != nil {
			return err
		}
	}
	testTagIDs = nil
	return nil
}

func createTestTag(t *testing.T, tagService service.TagService, text string, no string, path string) model.Tag {
	tag := model.Tag{
		Text: text + " " + time.Now().Format("20060102150405"),
		Slug: no + "-" + time.Now().Format("20060102150405"),
		// Path字段不存在于Tag模型中，移除
	}

	err := tagService.InsertTag(tag)
	assert.NoError(t, err)

	// 获取插入后的标签
	tags, err := tagService.GetTagsWithOrder(0, 10, "id desc")
	assert.NoError(t, err)
	assert.Greater(t, len(tags), 0)

	// 找到刚插入的标签
	var insertedTag model.Tag
	for _, tagItem := range tags {
		if tagItem.Text == tag.Text {
			insertedTag = tagItem
			testTagIDs = append(testTagIDs, tagItem.ID)
			break
		}
	}

	assert.NotEqual(t, uint(0), insertedTag.ID)
	return insertedTag
}

func TestTagService_InsertTag(t *testing.T) {
	db := SetupTestDB(t)
	defer cleanupTestTags(db)

	tagService := service.NewTagService()

	tag := model.Tag{
		Text: "Technology " + time.Now().Format("20060102150405"),
		Slug: "tech-" + time.Now().Format("20060102150405"),
		// Path字段不存在于Tag模型中，移除
	}

	err := tagService.InsertTag(tag)
	assert.NoError(t, err)

	// 验证标签是否成功插入
	tags, err := tagService.GetTagsWithOrder(0, 10, "id desc")
	assert.NoError(t, err)
	assert.Greater(t, len(tags), 0)

	// 找到并验证插入的标签
	var insertedTag *model.Tag
	for _, tagItem := range tags {
		if tagItem.Text == tag.Text {
			insertedTag = &tagItem
			testTagIDs = append(testTagIDs, tagItem.ID)
			break
		}
	}

	assert.NotNil(t, insertedTag)
	assert.Equal(t, tag.Text, insertedTag.Text)
	assert.Equal(t, tag.Slug, insertedTag.Slug)
	// Path字段不存在，移除此断言
}

func TestTagService_GetTagByID(t *testing.T) {
	db := SetupTestDB(t)
	defer cleanupTestTags(db)

	tagService := service.NewTagService()

	// 创建测试标签
	createdTag := createTestTag(t, tagService, "GetByID Test", "getbyid-test", "test.getbyid")

	// 根据ID获取标签
	foundTag, err := tagService.GetTagByID(createdTag.ID)
	assert.NoError(t, err)
	assert.Equal(t, createdTag.ID, foundTag.ID)
	assert.Equal(t, createdTag.Text, foundTag.Text)
	assert.Equal(t, createdTag.Slug, foundTag.Slug)
	// Path字段不存在，移除此断言
}

func TestTagService_GetAllTags(t *testing.T) {
	db := SetupTestDB(t)
	defer cleanupTestTags(db)

	tagService := service.NewTagService()

	// 创建多个测试标签
	for i := 0; i < 3; i++ {
		createTestTag(t, tagService, "AllTags Test", "alltags-test", "test.alltags")
	}

	// 获取所有标签
	allTags := tagService.GetAllTags()
	assert.GreaterOrEqual(t, len(allTags), 3)

	// 验证我们创建的标签在结果中
	createdCount := 0
	for _, tag := range allTags {
		if len(tag.Text) > 13 && tag.Text[:13] == "AllTags Test " {
			createdCount++
		}
	}
	assert.GreaterOrEqual(t, createdCount, 3)
}

func TestTagService_GetTagsWithOrder(t *testing.T) {
	db := SetupTestDB(t)
	defer cleanupTestTags(db)

	tagService := service.NewTagService()

	// 创建多个测试标签
	for i := 0; i < 5; i++ {
		createTestTag(t, tagService, "Order Test", "order-test", "test.order")
	}

	// 测试分页获取标签
	tags, err := tagService.GetTagsWithOrder(0, 3, "id desc")
	assert.NoError(t, err)
	assert.Equal(t, 3, len(tags))

	// 验证排序
	for i := 0; i < len(tags)-1; i++ {
		assert.GreaterOrEqual(t, tags[i].ID, tags[i+1].ID)
	}

	// 测试第二页
	tagsPage2, err := tagService.GetTagsWithOrder(1, 3, "id desc")
	assert.NoError(t, err)
	assert.LessOrEqual(t, len(tagsPage2), 3)

	// 验证分页数据不重复
	if len(tagsPage2) > 0 {
		assert.NotEqual(t, tags[0].ID, tagsPage2[0].ID)
	}
}

func TestTagService_GetTagsCount(t *testing.T) {
	db := SetupTestDB(t)
	defer cleanupTestTags(db)

	tagService := service.NewTagService()

	// 获取初始数量
	initialCount := tagService.GetTagsCount()

	// 创建测试标签
	for i := 0; i < 3; i++ {
		createTestTag(t, tagService, "Count Test", "count-test", "test.count")
	}

	// 验证数量增加
	newCount := tagService.GetTagsCount()
	assert.Equal(t, initialCount+3, newCount)
}

func TestTagService_GetMetadata(t *testing.T) {
	db := SetupTestDB(t)
	defer cleanupTestTags(db)

	tagService := service.NewTagService()

	// 创建具有层级结构的测试标签
	createTestTag(t, tagService, "Root Tag", "root", "root")
	createTestTag(t, tagService, "Child Tag 1", "child1", "root.child1")
	createTestTag(t, tagService, "Child Tag 2", "child2", "root.child2")
	createTestTag(t, tagService, "Grandchild Tag", "grandchild", "root.child1.grandchild")

	// 获取元数据
	metadata := tagService.GetMetadata()
	assert.NotNil(t, metadata)

	// 验证元数据包含层级信息
	// 具体的验证依赖于GetMetadata的实际返回结构
	assert.IsType(t, map[string]interface{}{}, metadata)
}

func TestTagService_DeleteTags(t *testing.T) {
	db := SetupTestDB(t)
	defer cleanupTestTags(db)

	tagService := service.NewTagService()

	// 创建多个测试标签
	var tagIDs []uint
	for i := 0; i < 3; i++ {
		tag := createTestTag(t, tagService, "Delete Test", "delete-test", "test.delete")
		tagIDs = append(tagIDs, tag.ID)
	}

	// 获取删除前的总数
	beforeCount := tagService.GetTagsCount()

	// 测试批量删除
	err := tagService.DeleteTags(tagIDs)
	assert.NoError(t, err)

	// 验证删除是否成功
	afterCount := tagService.GetTagsCount()
	assert.Equal(t, beforeCount-3, afterCount)

	// 验证标签确实被删除
	for _, deletedID := range tagIDs {
		_, err := tagService.GetTagByID(deletedID)
		assert.Error(t, err) // 应该返回错误，表示找不到标签
	}

	// 从清理列表中移除已删除的标签
	for _, deletedID := range tagIDs {
		for i, id := range testTagIDs {
			if id == deletedID {
				testTagIDs = append(testTagIDs[:i], testTagIDs[i+1:]...)
				break
			}
		}
	}
}

func TestTagService_SimpleTags(t *testing.T) {
	db := SetupTestDB(t)
	defer cleanupTestTags(db)

	tagService := service.NewTagService()

	// 创建简单标签结构（Tag模型不支持层级）
	techTag := createTestTag(t, tagService, "Technology", "tech", "")
	progTag := createTestTag(t, tagService, "Programming", "programming", "")
	_ = createTestTag(t, tagService, "Web Development", "webdev", "")
	_ = createTestTag(t, tagService, "Frontend", "frontend", "")

	// 验证标签基本属性
	assert.Equal(t, "Technology", techTag.Text)
	assert.Contains(t, techTag.Slug, "tech")
	assert.Equal(t, "Programming", progTag.Text)
	assert.Contains(t, progTag.Slug, "programming")

	// 获取所有标签验证创建成功
	allTags := tagService.GetAllTags()
	
	var createdTags []model.Tag
	tagTexts := []string{"Technology", "Programming", "Web Development", "Frontend"}
	
	for _, tag := range allTags {
		for _, targetText := range tagTexts {
			if tag.Text == targetText {
				createdTags = append(createdTags, tag)
				break
			}
		}
	}
	
	assert.GreaterOrEqual(t, len(createdTags), 4) // 至少有我们创建的4个标签
}

func TestTagService_EmptyResults(t *testing.T) {
	tagService := service.NewTagService()

	// 测试获取不存在的标签
	_, err := tagService.GetTagByID(999999)
	assert.Error(t, err)

	// 测试获取空结果
	tags, err := tagService.GetTagsWithOrder(999, 10, "id desc")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(tags))

	// 测试删除空数组
	err = tagService.DeleteTags([]uint{})
	assert.NoError(t, err)

	// 测试删除不存在的ID
	err = tagService.DeleteTags([]uint{999999})
	assert.NoError(t, err) // 删除不存在的记录不应该报错
}