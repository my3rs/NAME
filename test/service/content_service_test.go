package service

import (
	"NAME/database"
	"NAME/model"
	"NAME/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	testContentIDs []uint
)

func setupContentTestDB(t *testing.T) *gorm.DB {
	db := database.GetDB()
	testContentIDs = []uint{} // Reset test content IDs list
	return db
}

func cleanupTestContent(db *gorm.DB) error {
	if len(testContentIDs) > 0 {
		if err := db.Where("id IN ?", testContentIDs).Delete(&model.Content{}).Error; err != nil {
			return err
		}
	}
	testContentIDs = nil
	return nil
}

// 移除cleanupTestCategories函数，使用category_service_test.go中的版本

func createTestCategory(t *testing.T, db *gorm.DB) model.Category {
	category := model.Category{
		Text: "Test Category " + time.Now().Format("20060102150405"),
		Slug: "test-" + time.Now().Format("20060102150405"),
	}

	err := db.Create(&category).Error
	assert.NoError(t, err)
	return category
}

func createTestContent(t *testing.T, contentService service.ContentService, userService service.UserService) model.Content {
	// Create test user
	user := CreateTestUser(t, userService)

	// Create test category
	category := createTestCategory(t, database.GetDB())

	// Create test content
	content := model.Content{
		Type:          model.ContentTypePost,
		Title:         "Test Post " + time.Now().Format("20060102150405"),
		Abstract:      "Test Abstract " + time.Now().Format("20060102150405"),
		Text:          "# Test Content\n\nThis is a test post. " + time.Now().Format("20060102150405"),
		AuthorID:      user.ID,
		CategoryID:    category.ID,
		CreatedAt:     time.Now().UnixMilli(),
		UpdatedAt:     time.Now().UnixMilli(),
		PublishAt:     time.Now().UnixMilli(),
		Status:        model.ContentStatusPublished,
		AllowComment:  true,
		FeaturedImage: "https://example.com/image.jpg",
	}

	err := contentService.InsertPost(content)
	assert.NoError(t, err)

	// Get the inserted content by searching for it since ID might not be set
	posts := contentService.GetPostsWithOrder(0, 10, "id desc")
	assert.Greater(t, len(posts), 0)
	
	// Find the content we just created
	var insertedContent model.Content
	for _, post := range posts {
		if post.Title == content.Title {
			insertedContent = post
			break
		}
	}
	
	assert.NotEmpty(t, insertedContent.ID)
	testContentIDs = append(testContentIDs, insertedContent.ID)
	return insertedContent
}

func TestContentService_InsertPost(t *testing.T) {
	db := setupContentTestDB(t)
	defer cleanupTestContent(db)
	defer CleanupTestUsers(db)

	contentService := service.NewContentService()
	userService := service.NewUserService()
	user := CreateTestUser(t, userService)

	// 测试插入新文章
	content := model.Content{
		Type:          model.ContentTypePost,
		Title:         "New Post " + time.Now().Format("20060102150405"),
		Abstract:      "New Abstract " + time.Now().Format("20060102150405"),
		Text:          "# New Content\n\nThis is a new post. " + time.Now().Format("20060102150405"),
		AuthorID:      user.ID,
		Status:        model.ContentStatusPublished,
		AllowComment:  true,
		PublishAt:     time.Now().UnixMilli(),
		FeaturedImage: "https://example.com/new-image.jpg",
	}

	err := contentService.InsertPost(content)
	assert.NoError(t, err)

	// 获取并记录创建的内容ID - 通过查找最新的内容
	posts := contentService.GetPostsWithOrder(0, 10, "id desc")
	var insertedContent model.Content
	for _, post := range posts {
		if post.Title == content.Title {
			insertedContent = post
			testContentIDs = append(testContentIDs, post.ID)
			break
		}
	}

	// 验证文章是否成功插入
	assert.Equal(t, content.Title, insertedContent.Title)
	assert.Equal(t, content.Text, insertedContent.Text)
}

func TestContentService_UpdatePost(t *testing.T) {
	db := setupContentTestDB(t)
	defer cleanupTestContent(db)
	defer CleanupTestUsers(db)

	contentService := service.NewContentService()
	userService := service.NewUserService()
	content := createTestContent(t, contentService, userService)

	// 测试更新文章
	content.Title = "Updated Title " + time.Now().Format("20060102150405")
	content.Abstract = "Updated Abstract " + time.Now().Format("20060102150405")
	content.Text = "# Updated Content\n\nThis is an updated post. " + time.Now().Format("20060102150405")
	err := contentService.UpdatePost(content)
	assert.NoError(t, err)

	// 验证更新是否成功
	updatedContent := contentService.GetPostByID(int(content.ID))
	assert.Equal(t, content.Title, updatedContent.Title)
	assert.Equal(t, content.Abstract, updatedContent.Abstract)
}

func TestContentService_DeletePost(t *testing.T) {
	db := setupContentTestDB(t)
	defer cleanupTestContent(db)
	defer CleanupTestUsers(db)

	contentService := service.NewContentService()
	userService := service.NewUserService()
	content := createTestContent(t, contentService, userService)

	// 测试删除文章
	err := contentService.DeletePostByID(content.ID)
	assert.NoError(t, err)

	// 验证文章是否已被删除
	deletedContent := contentService.GetPostByID(int(content.ID))
	assert.Equal(t, uint(0), deletedContent.ID)
}

func TestContentService_GetFormattedPostByID(t *testing.T) {
	db := setupContentTestDB(t)
	defer cleanupTestContent(db)
	defer CleanupTestUsers(db)

	contentService := service.NewContentService()
	userService := service.NewUserService()
	content := createTestContent(t, contentService, userService)

	// 测试获取格式化后的文章
	formattedContent := contentService.GetFormattedPostByID(int(content.ID))
	assert.NotEmpty(t, formattedContent.TextHTML)
	assert.Contains(t, string(formattedContent.TextHTML), "<h1>Test Content</h1>")
}

func TestContentService_GetPostsWithOrder(t *testing.T) {
	db := setupContentTestDB(t)
	defer cleanupTestContent(db)
	defer CleanupTestUsers(db)

	contentService := service.NewContentService()
	userService := service.NewUserService()

	// 创建多篇测试文章
	for i := 0; i < 5; i++ {
		createTestContent(t, contentService, userService)
	}

	// 测试分页获取文章
	posts := contentService.GetPostsWithOrder(0, 3, "created_at desc")
	assert.Equal(t, 3, len(posts))

	// 验证排序是否正确
	for i := 0; i < len(posts)-1; i++ {
		assert.GreaterOrEqual(t, posts[i].CreatedAt, posts[i+1].CreatedAt)
	}
}

func TestContentService_GetMeta(t *testing.T) {
	db := setupContentTestDB(t)
	defer cleanupTestContent(db)
	defer CleanupTestUsers(db)

	contentService := service.NewContentService()
	userService := service.NewUserService()

	// 创建一些测试数据
	createTestContent(t, contentService, userService)
	createTestContent(t, contentService, userService)

	// 测试获取元信息
	meta := contentService.GetMeta()
	assert.NotNil(t, meta)
	assert.Contains(t, meta, "posts")
	assert.Contains(t, meta, "pages")
	assert.Contains(t, meta, "categories")
}
