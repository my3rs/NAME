package service

import (
	"NAME/model"
	"NAME/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var testCommentIDs []uint

func cleanupTestComments(db *gorm.DB) error {
	if len(testCommentIDs) > 0 {
		if err := db.Where("id IN ?", testCommentIDs).Delete(&model.Comment{}).Error; err != nil {
			return err
		}
	}
	testCommentIDs = nil
	return nil
}

func createTestComment(t *testing.T, commentService service.CommentService, contentService service.ContentService, userService service.UserService, parentPath string) model.Comment {
	// 创建测试内容
	content := createTestContent(t, contentService, userService)

	comment := model.Comment{
		AuthorName: "Test Author " + time.Now().Format("20060102150405"),
		Mail:      "test" + time.Now().Format("20060102150405") + "@example.com",
		URL:       "https://example.com",
		Text:      "This is a test comment " + time.Now().Format("20060102150405"),
		IP:        "127.0.0.1",
		Agent:     "Test User Agent",
		ContentID: content.ID,
		ParentID:  0,
		Path:      parentPath,
		Status:    model.CommentStatusApproved,
	}

	err := commentService.InsertComment(comment)
	assert.NoError(t, err)

	// 获取插入后的评论
	comments, err := commentService.GetComments(0, 10, "created_at desc")
	assert.NoError(t, err)
	assert.Greater(t, len(comments), 0)

	// 找到刚插入的评论
	var insertedComment model.Comment
	for _, c := range comments {
		if c.Text == comment.Text {
			insertedComment = c
			testCommentIDs = append(testCommentIDs, c.ID)
			break
		}
	}

	assert.NotEqual(t, uint(0), insertedComment.ID)
	return insertedComment
}

func TestCommentService_InsertComment(t *testing.T) {
	db := SetupTestDB(t)
	defer cleanupTestComments(db)
	defer cleanupTestContent(db)
	defer CleanupTestUsers(db)

	commentService := service.NewCommentService()
	contentService := service.NewContentService()
	userService := service.NewUserService()

	// 创建测试内容
	content := createTestContent(t, contentService, userService)

	comment := model.Comment{
		AuthorName: "John Doe",
		Mail:      "john" + time.Now().Format("20060102150405") + "@example.com",
		URL:       "https://johndoe.com",
		Text:      "This is a test comment",
		IP:        "192.168.1.1",
		Agent:     "Mozilla/5.0",
		ContentID: content.ID,
		ParentID:  0,
		Path:      "",
		Status:    model.CommentStatusApproved,
	}

	err := commentService.InsertComment(comment)
	assert.NoError(t, err)

	// 验证评论是否成功插入
	comments, err := commentService.GetComments(0, 10, "created_at desc")
	assert.NoError(t, err)
	assert.Greater(t, len(comments), 0)

	// 找到并验证插入的评论
	var insertedComment *model.Comment
	for _, c := range comments {
		if c.Text == comment.Text {
			insertedComment = &c
			testCommentIDs = append(testCommentIDs, c.ID)
			break
		}
	}

	assert.NotNil(t, insertedComment)
	assert.Equal(t, comment.AuthorName, insertedComment.AuthorName)
	assert.Equal(t, comment.Mail, insertedComment.Mail)
	assert.Equal(t, comment.Text, insertedComment.Text)
	assert.Equal(t, content.ID, insertedComment.ContentID)
}

func TestCommentService_GetComments(t *testing.T) {
	db := SetupTestDB(t)
	defer cleanupTestComments(db)
	defer cleanupTestContent(db)
	defer CleanupTestUsers(db)

	commentService := service.NewCommentService()
	contentService := service.NewContentService()
	userService := service.NewUserService()

	// 创建多个测试评论
	for i := 0; i < 5; i++ {
		createTestComment(t, commentService, contentService, userService, "")
	}

	// 测试分页获取评论
	comments, err := commentService.GetComments(0, 3, "created_at desc")
	assert.NoError(t, err)
	assert.Equal(t, 3, len(comments))

	// 验证排序
	for i := 0; i < len(comments)-1; i++ {
		assert.GreaterOrEqual(t, comments[i].CreatedAt, comments[i+1].CreatedAt)
	}

	// 测试第二页
	commentsPage2, err := commentService.GetComments(1, 3, "created_at desc")
	assert.NoError(t, err)
	assert.LessOrEqual(t, len(commentsPage2), 3)
}

func TestCommentService_GetCommentsByContentID(t *testing.T) {
	db := SetupTestDB(t)
	defer cleanupTestComments(db)
	defer cleanupTestContent(db)
	defer CleanupTestUsers(db)

	commentService := service.NewCommentService()
	contentService := service.NewContentService()
	userService := service.NewUserService()

	// 创建测试内容
	content := createTestContent(t, contentService, userService)

	// 为该内容创建多个评论
	for i := 0; i < 3; i++ {
		comment := model.Comment{
			AuthorName: "Author " + string(rune(65+i)),
			Mail:      "author" + string(rune(97+i)) + time.Now().Format("20060102150405") + "@example.com",
			Text:      "Comment " + string(rune(48+i)) + " for content",
			ContentID: content.ID,
			Status:    model.CommentStatusApproved,
		}
		err := commentService.InsertComment(comment)
		assert.NoError(t, err)
	}

	// 获取该内容的评论
	comments, err := commentService.GetCommentsByContentID(int(content.ID), 0, 10, "created_at desc")
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(comments), 3)

	// 记录评论ID用于清理
	for _, comment := range comments {
		testCommentIDs = append(testCommentIDs, comment.ID)
	}

	// 验证所有评论都属于该内容
	for _, comment := range comments {
		assert.Equal(t, content.ID, comment.ContentID)
	}
}

func TestCommentService_GetCommentsCount(t *testing.T) {
	db := SetupTestDB(t)
	defer cleanupTestComments(db)
	defer cleanupTestContent(db)
	defer CleanupTestUsers(db)

	commentService := service.NewCommentService()
	contentService := service.NewContentService()
	userService := service.NewUserService()

	// 创建测试内容
	content := createTestContent(t, contentService, userService)

	// 获取初始评论数量
	initialCount := commentService.GetCommentsCount(int64(content.ID))

	// 创建评论
	for i := 0; i < 3; i++ {
		comment := model.Comment{
			AuthorName: "Counter " + string(rune(65+i)),
			Mail:      "counter" + string(rune(97+i)) + time.Now().Format("20060102150405") + "@example.com",
			Text:      "Count test comment " + string(rune(48+i)),
			ContentID: content.ID,
			Status:    model.CommentStatusApproved,
		}
		err := commentService.InsertComment(comment)
		assert.NoError(t, err)
	}

	// 验证评论数量增加
	newCount := commentService.GetCommentsCount(int64(content.ID))
	assert.Equal(t, initialCount+3, newCount)

	// 获取所有评论以便清理
	comments, err := commentService.GetCommentsByContentID(int(content.ID), 0, 100, "id desc")
	assert.NoError(t, err)
	for _, comment := range comments {
		testCommentIDs = append(testCommentIDs, comment.ID)
	}
}

func TestCommentService_GetCommentByID(t *testing.T) {
	db := SetupTestDB(t)
	defer cleanupTestComments(db)
	defer cleanupTestContent(db)
	defer CleanupTestUsers(db)

	commentService := service.NewCommentService()
	contentService := service.NewContentService()
	userService := service.NewUserService()

	// 创建测试评论
	createdComment := createTestComment(t, commentService, contentService, userService, "")

	// 根据ID获取评论
	foundComment, err := commentService.GetCommentByID(int(createdComment.ID))
	assert.NoError(t, err)
	assert.Equal(t, createdComment.ID, foundComment.ID)
	assert.Equal(t, createdComment.AuthorName, foundComment.AuthorName)
	assert.Equal(t, createdComment.Text, foundComment.Text)
}

func TestCommentService_UpdateComment(t *testing.T) {
	db := SetupTestDB(t)
	defer cleanupTestComments(db)
	defer cleanupTestContent(db)
	defer CleanupTestUsers(db)

	commentService := service.NewCommentService()
	contentService := service.NewContentService()
	userService := service.NewUserService()

	// 创建测试评论
	comment := createTestComment(t, commentService, contentService, userService, "")

	// 更新评论
	comment.Text = "Updated comment text " + time.Now().Format("20060102150405")
	comment.Status = model.CommentStatusUnreviewed
	err := commentService.UpdateComment(comment)
	assert.NoError(t, err)

	// 验证更新是否成功
	updatedComment, err := commentService.GetCommentByID(int(comment.ID))
	assert.NoError(t, err)
	assert.Equal(t, comment.Text, updatedComment.Text)
	assert.Equal(t, comment.Status, updatedComment.Status)
}

func TestCommentService_DeleteComment(t *testing.T) {
	db := SetupTestDB(t)
	defer cleanupTestComments(db)
	defer cleanupTestContent(db)
	defer CleanupTestUsers(db)

	commentService := service.NewCommentService()
	contentService := service.NewContentService()
	userService := service.NewUserService()

	// 创建测试评论
	comment := createTestComment(t, commentService, contentService, userService, "")

	// 删除评论
	err := commentService.DeleteComment(comment.ID)
	assert.NoError(t, err)

	// 验证评论是否被删除
	_, err = commentService.GetCommentByID(int(comment.ID))
	assert.Error(t, err) // 应该返回错误，表示找不到评论

	// 从清理列表中移除，因为已经删除
	for i, id := range testCommentIDs {
		if id == comment.ID {
			testCommentIDs = append(testCommentIDs[:i], testCommentIDs[i+1:]...)
			break
		}
	}
}

func TestCommentService_DeleteBatchComment(t *testing.T) {
	db := SetupTestDB(t)
	defer cleanupTestComments(db)
	defer cleanupTestContent(db)
	defer CleanupTestUsers(db)

	commentService := service.NewCommentService()
	contentService := service.NewContentService()
	userService := service.NewUserService()

	// 创建多个测试评论
	var commentIDs []uint
	for i := 0; i < 3; i++ {
		comment := createTestComment(t, commentService, contentService, userService, "")
		commentIDs = append(commentIDs, comment.ID)
	}

	// 批量删除评论
	deletedCount, err := commentService.DeleteBatchComment(commentIDs)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), deletedCount)

	// 验证评论是否被删除
	for _, id := range commentIDs {
		_, err := commentService.GetCommentByID(int(id))
		assert.Error(t, err) // 应该返回错误，表示找不到评论
	}

	// 清理列表中移除已删除的评论
	testCommentIDs = []uint{}
}

func TestCommentService_GetCommentsWithContentTitle(t *testing.T) {
	db := SetupTestDB(t)
	defer cleanupTestComments(db)
	defer cleanupTestContent(db)
	defer CleanupTestUsers(db)

	commentService := service.NewCommentService()
	contentService := service.NewContentService()
	userService := service.NewUserService()

	// 创建测试评论
	createTestComment(t, commentService, contentService, userService, "")

	// 获取包含内容标题的评论
	comments, err := commentService.GetCommentsWithContentTitle(0, 10, "created_at desc")
	assert.NoError(t, err)
	assert.Greater(t, len(comments), 0)

	// 验证评论包含内容信息
	for _, comment := range comments {
		assert.NotEqual(t, uint(0), comment.ContentID)
		// 注意：这里需要检查Content字段是否被正确加载
		// 具体的验证依赖于实际的数据库关联设置
	}
}