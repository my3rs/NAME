package service

import (
	"NAME/model"
	"NAME/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var testCategoryIDs []uint

func cleanupTestCategories(db *gorm.DB) error {
	if len(testCategoryIDs) > 0 {
		if err := db.Where("id IN ?", testCategoryIDs).Delete(&model.Category{}).Error; err != nil {
			return err
		}
	}
	testCategoryIDs = nil
	return nil
}

func TestCategoryService_InsertCategory(t *testing.T) {
	db := SetupTestDB(t)
	defer cleanupTestCategories(db)

	categoryService := service.NewCategoryService()

	category := model.Category{
		Text: "Technology " + time.Now().Format("20060102150405"),
		Slug: "tech-" + time.Now().Format("20060102150405"),
	}

	err := categoryService.InsertCategory(category)
	assert.NoError(t, err)

	// 验证分类是否成功插入
	categories := categoryService.GetCategories(0, 10, "id desc")
	assert.Greater(t, len(categories), 0)
	
	// 找到刚插入的分类
	var insertedCategory *model.Category
	for _, cat := range categories {
		if cat.Text == category.Text {
			insertedCategory = &cat
			testCategoryIDs = append(testCategoryIDs, cat.ID)
			break
		}
	}
	
	assert.NotNil(t, insertedCategory)
	assert.Equal(t, category.Text, insertedCategory.Text)
	assert.Equal(t, category.Slug, insertedCategory.Slug)
}

func TestCategoryService_GetCategories(t *testing.T) {
	db := SetupTestDB(t)
	defer cleanupTestCategories(db)

	categoryService := service.NewCategoryService()

	// 创建多个测试分类
	for i := 0; i < 5; i++ {
		category := model.Category{
			Text: "Category " + time.Now().Format("20060102150405") + string(rune(65+i)),
			Slug: "cat-" + time.Now().Format("20060102150405") + string(rune(97+i)),
		}
		err := categoryService.InsertCategory(category)
		assert.NoError(t, err)
	}

	// 测试分页获取分类
	categories := categoryService.GetCategories(0, 3, "id desc")
	assert.Equal(t, 3, len(categories))

	// 记录测试数据ID用于清理
	for _, cat := range categories {
		if len(testCategoryIDs) < 5 { // 避免重复添加
			testCategoryIDs = append(testCategoryIDs, cat.ID)
		}
	}

	// 测试第二页
	categoriesPage2 := categoryService.GetCategories(1, 3, "id desc")
	assert.LessOrEqual(t, len(categoriesPage2), 3)

	// 验证分页数据不重复
	if len(categoriesPage2) > 0 {
		assert.NotEqual(t, categories[0].ID, categoriesPage2[0].ID)
	}
}

func TestCategoryService_GetCategoriesCount(t *testing.T) {
	db := SetupTestDB(t)
	defer cleanupTestCategories(db)

	categoryService := service.NewCategoryService()

	// 获取初始数量
	initialCount := categoryService.GetCategoriesCount()

	// 创建测试分类
	for i := 0; i < 3; i++ {
		category := model.Category{
			Text: "Count Test " + time.Now().Format("20060102150405") + string(rune(48+i)),
			Slug: "count-" + time.Now().Format("20060102150405") + string(rune(48+i)),
		}
		err := categoryService.InsertCategory(category)
		assert.NoError(t, err)
	}

	// 验证数量增加
	newCount := categoryService.GetCategoriesCount()
	assert.Equal(t, initialCount+3, newCount)

	// 获取所有分类以便清理
	allCategories := categoryService.GetCategories(1, 100, "id desc")
	for _, cat := range allCategories {
		testCategoryIDs = append(testCategoryIDs, cat.ID)
	}
}

func TestCategoryService_UpdateCategory(t *testing.T) {
	db := SetupTestDB(t)
	defer cleanupTestCategories(db)

	categoryService := service.NewCategoryService()

	// 创建测试分类
	category := model.Category{
		Text: "Original Category " + time.Now().Format("20060102150405"),
		Slug: "original-" + time.Now().Format("20060102150405"),
	}

	err := categoryService.InsertCategory(category)
	assert.NoError(t, err)

	// 获取插入后的分类（包含ID）
	categories := categoryService.GetCategories(0, 10, "id desc")
	var insertedCategory *model.Category
	for _, cat := range categories {
		if cat.Text == category.Text {
			insertedCategory = &cat
			testCategoryIDs = append(testCategoryIDs, cat.ID)
			break
		}
	}
	assert.NotNil(t, insertedCategory)

	// 更新分类
	insertedCategory.Text = "Updated Category " + time.Now().Format("20060102150405")
	insertedCategory.Slug = "updated-" + time.Now().Format("20060102150405")
	err = categoryService.UpdateCategory(*insertedCategory)
	assert.NoError(t, err)

	// 验证更新是否成功
	updatedCategories := categoryService.GetCategories(0, 10, "id desc")
	var updatedCategory *model.Category
	for _, cat := range updatedCategories {
		if cat.ID == insertedCategory.ID {
			updatedCategory = &cat
			break
		}
	}
	
	assert.NotNil(t, updatedCategory)
	assert.Equal(t, insertedCategory.Text, updatedCategory.Text)
	assert.Equal(t, insertedCategory.Slug, updatedCategory.Slug)
}

func TestCategoryService_DeleteCategories(t *testing.T) {
	db := SetupTestDB(t)
	defer cleanupTestCategories(db)

	categoryService := service.NewCategoryService()

	// 创建多个测试分类
	var categoryIDs []uint
	for i := 0; i < 3; i++ {
		category := model.Category{
			Text: "Delete Test " + time.Now().Format("20060102150405") + string(rune(48+i)),
			Slug: "delete-" + time.Now().Format("20060102150405") + string(rune(48+i)),
		}
		err := categoryService.InsertCategory(category)
		assert.NoError(t, err)
	}

	// 获取刚创建的分类ID
	categories := categoryService.GetCategories(0, 10, "id desc")
	for i, cat := range categories {
		if i < 3 { // 只取前3个
			categoryIDs = append(categoryIDs, cat.ID)
		}
	}
	assert.Equal(t, 3, len(categoryIDs))

	// 获取删除前的总数
	beforeCount := categoryService.GetCategoriesCount()

	// 测试批量删除
	err := categoryService.DeleteCategories(categoryIDs)
	assert.NoError(t, err)

	// 验证删除是否成功
	afterCount := categoryService.GetCategoriesCount()
	assert.Equal(t, beforeCount-3, afterCount)

	// 验证分类确实被删除
	remainingCategories := categoryService.GetCategories(0, 10, "id desc")
	for _, cat := range remainingCategories {
		for _, deletedID := range categoryIDs {
			assert.NotEqual(t, deletedID, cat.ID)
		}
	}
}

func TestCategoryService_EmptyResults(t *testing.T) {
	categoryService := service.NewCategoryService()

	// 测试获取空结果
	categories := categoryService.GetCategories(999, 10, "id desc")
	assert.Equal(t, 0, len(categories))

	// 测试删除空数组
	err := categoryService.DeleteCategories([]uint{})
	assert.NoError(t, err)

	// 测试删除不存在的ID
	err = categoryService.DeleteCategories([]uint{999999})
	assert.NoError(t, err) // 删除不存在的记录不应该报错
}