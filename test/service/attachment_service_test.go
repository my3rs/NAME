package service

import (
	"NAME/model"
	"NAME/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var testAttachmentIDs []uint

func cleanupTestAttachments(db *gorm.DB) error {
	if len(testAttachmentIDs) > 0 {
		if err := db.Where("id IN ?", testAttachmentIDs).Delete(&model.Attachment{}).Error; err != nil {
			return err
		}
	}
	testAttachmentIDs = nil
	return nil
}

func createTestAttachment(t *testing.T, attachmentService service.AttachmentService, contentService service.ContentService, userService service.UserService) model.Attachment {
	// 创建测试内容
	content := createTestContent(t, contentService, userService)

	attachment := model.Attachment{
		Title:     "test_file_" + time.Now().Format("20060102150405") + ".jpg",
		Path:      "/uploads/test_file_" + time.Now().Format("20060102150405") + ".jpg",
		ContentID: content.ID,
		Content:   content,
		CreatedAt: time.Now().Unix(),
	}

	err := attachmentService.InsertAttachment(attachment)
	assert.NoError(t, err)

	// 获取插入后的附件（通过查询数据库验证）
	db := SetupTestDB(t)
	var insertedAttachment model.Attachment
	err = db.Where("title = ?", attachment.Title).First(&insertedAttachment).Error
	assert.NoError(t, err)
	
	testAttachmentIDs = append(testAttachmentIDs, insertedAttachment.ID)
	return insertedAttachment
}

func TestAttachmentService_InsertAttachment(t *testing.T) {
	db := SetupTestDB(t)
	defer cleanupTestAttachments(db)
	defer cleanupTestContent(db)
	defer CleanupTestUsers(db)

	attachmentService := service.NewAttachmentService()
	contentService := service.NewContentService()
	userService := service.NewUserService()

	// 创建测试内容
	content := createTestContent(t, contentService, userService)

	attachment := model.Attachment{
		Title:     "sample_image.jpg",
		Path:      "/uploads/sample_image.jpg",
		ContentID: content.ID,
		Content:   content,
		CreatedAt: time.Now().Unix(),
	}

	err := attachmentService.InsertAttachment(attachment)
	assert.NoError(t, err)

	// 验证附件是否成功插入
	var insertedAttachment model.Attachment
	err = db.Where("title = ?", attachment.Title).First(&insertedAttachment).Error
	assert.NoError(t, err)

	assert.Equal(t, attachment.Title, insertedAttachment.Title)
	assert.Equal(t, attachment.Path, insertedAttachment.Path)
	assert.Equal(t, content.ID, insertedAttachment.ContentID)
	
	testAttachmentIDs = append(testAttachmentIDs, insertedAttachment.ID)
}

func TestAttachmentService_InsertAttachment_WithoutContent(t *testing.T) {
	db := SetupTestDB(t)
	defer cleanupTestAttachments(db)

	attachmentService := service.NewAttachmentService()

	// 测试不关联内容的附件
	attachment := model.Attachment{
		Title:     "standalone_file.pdf",
		Path:      "/uploads/standalone_file.pdf",
		ContentID: 0, // 不关联任何内容
		CreatedAt: time.Now().Unix(),
	}

	err := attachmentService.InsertAttachment(attachment)
	assert.NoError(t, err)

	// 验证附件是否成功插入
	var insertedAttachment model.Attachment
	err = db.Where("title = ?", attachment.Title).First(&insertedAttachment).Error
	assert.NoError(t, err)

	assert.Equal(t, attachment.Title, insertedAttachment.Title)
	assert.Equal(t, attachment.Path, insertedAttachment.Path)
	assert.Equal(t, uint(0), insertedAttachment.ContentID)
	
	testAttachmentIDs = append(testAttachmentIDs, insertedAttachment.ID)
}

func TestAttachmentService_InsertAttachment_VariousFileTypes(t *testing.T) {
	db := SetupTestDB(t)
	defer cleanupTestAttachments(db)
	defer cleanupTestContent(db)
	defer CleanupTestUsers(db)

	attachmentService := service.NewAttachmentService()
	contentService := service.NewContentService()
	userService := service.NewUserService()

	// 创建测试内容
	content := createTestContent(t, contentService, userService)

	// 测试不同类型的文件
	fileTypes := []struct {
		title string
		path  string
	}{
		{"document.pdf", "/uploads/document.pdf"},
		{"image.png", "/uploads/image.png"},
		{"video.mp4", "/uploads/video.mp4"},
		{"audio.mp3", "/uploads/audio.mp3"},
		{"archive.zip", "/uploads/archive.zip"},
		{"text.txt", "/uploads/text.txt"},
	}

	for _, fileType := range fileTypes {
		attachment := model.Attachment{
			Title:     fileType.title + "_" + time.Now().Format("20060102150405"),
			Path:      fileType.path + "_" + time.Now().Format("20060102150405"),
			ContentID: content.ID,
			CreatedAt: time.Now().Unix(),
		}

		err := attachmentService.InsertAttachment(attachment)
		assert.NoError(t, err, "Failed to insert %s", fileType.title)

		// 验证文件是否成功插入
		var insertedAttachment model.Attachment
		err = db.Where("title = ?", attachment.Title).First(&insertedAttachment).Error
		assert.NoError(t, err, "Failed to find inserted %s", fileType.title)

		assert.Equal(t, attachment.Title, insertedAttachment.Title)
		assert.Equal(t, attachment.Path, insertedAttachment.Path)
		
		testAttachmentIDs = append(testAttachmentIDs, insertedAttachment.ID)
	}
}

func TestAttachmentService_InsertAttachment_WithContentAssociation(t *testing.T) {
	db := SetupTestDB(t)
	defer cleanupTestAttachments(db)
	defer cleanupTestContent(db)
	defer CleanupTestUsers(db)

	attachmentService := service.NewAttachmentService()
	contentService := service.NewContentService()
	userService := service.NewUserService()

	// 创建多个测试内容
	content1 := createTestContent(t, contentService, userService)
	content2 := createTestContent(t, contentService, userService)

	// 为不同内容创建附件
	attachment1 := model.Attachment{
		Title:     "content1_image.jpg",
		Path:      "/uploads/content1_image.jpg",
		ContentID: content1.ID,
		Content:   content1,
		CreatedAt: time.Now().Unix(),
	}

	attachment2 := model.Attachment{
		Title:     "content2_document.pdf",
		Path:      "/uploads/content2_document.pdf",
		ContentID: content2.ID,
		Content:   content2,
		CreatedAt: time.Now().Unix(),
	}

	// 插入附件
	err := attachmentService.InsertAttachment(attachment1)
	assert.NoError(t, err)

	err = attachmentService.InsertAttachment(attachment2)
	assert.NoError(t, err)

	// 验证附件与内容的关联
	var insertedAttachment1, insertedAttachment2 model.Attachment
	
	err = db.Where("title = ?", attachment1.Title).First(&insertedAttachment1).Error
	assert.NoError(t, err)
	assert.Equal(t, content1.ID, insertedAttachment1.ContentID)
	
	err = db.Where("title = ?", attachment2.Title).First(&insertedAttachment2).Error
	assert.NoError(t, err)
	assert.Equal(t, content2.ID, insertedAttachment2.ContentID)

	testAttachmentIDs = append(testAttachmentIDs, insertedAttachment1.ID, insertedAttachment2.ID)
}

func TestAttachmentService_InsertAttachment_DuplicateTitle(t *testing.T) {
	db := SetupTestDB(t)
	defer cleanupTestAttachments(db)

	attachmentService := service.NewAttachmentService()

	// 创建第一个附件
	attachment1 := model.Attachment{
		Title:     "duplicate_test.jpg",
		Path:      "/uploads/duplicate_test1.jpg",
		CreatedAt: time.Now().Unix(),
	}

	err := attachmentService.InsertAttachment(attachment1)
	assert.NoError(t, err)

	// 尝试创建相同标题的附件
	attachment2 := model.Attachment{
		Title:     "duplicate_test.jpg", // 相同的标题
		Path:      "/uploads/duplicate_test2.jpg", // 不同的路径
		CreatedAt: time.Now().Unix(),
	}

	err = attachmentService.InsertAttachment(attachment2)
	// 根据实际业务逻辑，这里可能成功或失败
	// 如果数据库没有唯一约束，应该成功
	// 如果有唯一约束，应该失败
	// 这里假设没有唯一约束，所以应该成功
	assert.NoError(t, err)

	// 获取所有同名附件
	var attachments []model.Attachment
	err = db.Where("title = ?", "duplicate_test.jpg").Find(&attachments).Error
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(attachments), 2)

	// 记录ID用于清理
	for _, att := range attachments {
		testAttachmentIDs = append(testAttachmentIDs, att.ID)
	}
}

func TestAttachmentService_InsertAttachment_EmptyFields(t *testing.T) {
	db := SetupTestDB(t)
	defer cleanupTestAttachments(db)

	attachmentService := service.NewAttachmentService()

	// 测试空字段的处理
	attachment := model.Attachment{
		Title:     "", // 空标题
		Path:      "", // 空路径
		ContentID: 0,
		CreatedAt: time.Now().Unix(),
	}

	err := attachmentService.InsertAttachment(attachment)
	// 根据业务逻辑，可能允许或不允许空字段
	// 这里假设允许空字段
	assert.NoError(t, err)

	// 验证插入
	var insertedAttachment model.Attachment
	err = db.Where("created_at = ?", attachment.CreatedAt).First(&insertedAttachment).Error
	assert.NoError(t, err)
	
	testAttachmentIDs = append(testAttachmentIDs, insertedAttachment.ID)
}

func TestAttachmentService_InsertAttachment_LongPaths(t *testing.T) {
	db := SetupTestDB(t)
	defer cleanupTestAttachments(db)

	attachmentService := service.NewAttachmentService()

	// 测试长路径
	longPath := "/uploads/very/long/path/structure/that/might/be/used/in/real/applications/with/deep/directory/nesting/test_file_" + time.Now().Format("20060102150405") + ".jpg"
	
	attachment := model.Attachment{
		Title:     "long_path_test.jpg",
		Path:      longPath,
		CreatedAt: time.Now().Unix(),
	}

	err := attachmentService.InsertAttachment(attachment)
	assert.NoError(t, err)

	// 验证长路径是否正确保存
	var insertedAttachment model.Attachment
	err = db.Where("title = ?", attachment.Title).First(&insertedAttachment).Error
	assert.NoError(t, err)
	assert.Equal(t, longPath, insertedAttachment.Path)
	
	testAttachmentIDs = append(testAttachmentIDs, insertedAttachment.ID)
}