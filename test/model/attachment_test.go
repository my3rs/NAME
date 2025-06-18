package model

import (
	"NAME/model"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAttachment_JSONSerialization(t *testing.T) {
	attachment := model.Attachment{
		ID:        1,
		Title:     "test_image.jpg",
		Path:      "/uploads/test_image.jpg",
		ContentID: 5,
		Content: model.Content{
			ID:    5,
			Title: "Test Post",
			Type:  model.ContentTypePost,
		},
		CreatedAt: time.Now().Unix(),
	}

	// 测试序列化
	jsonData, err := json.Marshal(attachment)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// 验证JSON结构
	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonData, &jsonMap)
	assert.NoError(t, err)

	// 验证字段映射
	assert.Equal(t, float64(1), jsonMap["id"])
	assert.Equal(t, "test_image.jpg", jsonMap["title"]) // Title field
	assert.Equal(t, "/uploads/test_image.jpg", jsonMap["path"])
	
	// ContentID应该被忽略（json:"-"）
	_, exists := jsonMap["contentid"]
	assert.False(t, exists)
	_, exists = jsonMap["ContentID"]
	assert.False(t, exists)

	// Content对象应该存在
	content, exists := jsonMap["content"]
	assert.True(t, exists)
	assert.NotNil(t, content)

	// CreatedAt应该存在
	createdAt, exists := jsonMap["createdAt"]
	assert.True(t, exists)
	assert.NotNil(t, createdAt)
}

func TestAttachment_JSONDeserialization(t *testing.T) {
	jsonStr := `{
		"id": 2,
		"title": "document.pdf",
		"path": "/uploads/document.pdf",
		"content": {
			"id": 10,
			"title": "Important Document",
			"type": "post"
		},
		"createdAt": 1640995200
	}`

	var attachment model.Attachment
	err := json.Unmarshal([]byte(jsonStr), &attachment)
	assert.NoError(t, err)

	assert.Equal(t, uint(2), attachment.ID)
	assert.Equal(t, "document.pdf", attachment.Title)
	assert.Equal(t, "/uploads/document.pdf", attachment.Path)
	assert.Equal(t, int64(1640995200), attachment.CreatedAt)
	
	// 验证关联的Content对象
	assert.Equal(t, uint(10), attachment.Content.ID)
	assert.Equal(t, "Important Document", attachment.Content.Title)
	assert.Equal(t, model.ContentTypePost, attachment.Content.Type)
}

func TestAttachment_EmptyAttachment(t *testing.T) {
	attachment := model.Attachment{}
	
	assert.Equal(t, uint(0), attachment.ID)
	assert.Equal(t, "", attachment.Title)
	assert.Equal(t, "", attachment.Path)
	assert.Equal(t, uint(0), attachment.ContentID)
	assert.Equal(t, int64(0), attachment.CreatedAt)
}

func TestAttachment_WithContent(t *testing.T) {
	content := model.Content{
		ID:    1,
		Title: "Test Content",
		Type:  model.ContentTypePost,
	}

	attachment := model.Attachment{
		ID:        1,
		Title:     "image.png",
		Path:      "/uploads/image.png",
		ContentID: content.ID,
		Content:   content,
		CreatedAt: time.Now().Unix(),
	}

	assert.Equal(t, content.ID, attachment.ContentID)
	assert.Equal(t, content.ID, attachment.Content.ID)
	assert.Equal(t, content.Title, attachment.Content.Title)
	assert.Equal(t, content.Type, attachment.Content.Type)
}

func TestAttachment_FileTypes(t *testing.T) {
	testCases := []struct {
		title string
		path  string
	}{
		{"image.jpg", "/uploads/image.jpg"},
		{"document.pdf", "/uploads/document.pdf"},
		{"video.mp4", "/uploads/video.mp4"},
		{"audio.mp3", "/uploads/audio.mp3"},
		{"archive.zip", "/uploads/archive.zip"},
		{"text.txt", "/uploads/text.txt"},
	}

	for _, tc := range testCases {
		attachment := model.Attachment{
			Title: tc.title,
			Path:  tc.path,
		}
		
		assert.Equal(t, tc.title, attachment.Title)
		assert.Equal(t, tc.path, attachment.Path)
	}
}

func TestAttachment_CreatedAtHandling(t *testing.T) {
	now := time.Now().Unix()
	attachment := model.Attachment{
		ID:        1,
		Title:     "test.jpg",
		Path:      "/uploads/test.jpg",
		CreatedAt: now,
	}

	assert.Equal(t, now, attachment.CreatedAt)
	
	// 测试JSON序列化是否保持时间戳
	jsonData, err := json.Marshal(attachment)
	assert.NoError(t, err)
	
	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonData, &jsonMap)
	assert.NoError(t, err)
	
	assert.Equal(t, float64(now), jsonMap["createdAt"])
}