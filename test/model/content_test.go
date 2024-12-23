package model_test

import (
	"encoding/json"
	"html/template"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"NAME/model"
)

func TestContentType(t *testing.T) {
	assert.Equal(t, model.ContentType("post"), model.ContentTypePost)
	assert.Equal(t, model.ContentType("digu"), model.ContentTypeDigu)
	assert.Equal(t, model.ContentType("page"), model.ContentTypePage)
}

func TestContentStatus(t *testing.T) {
	assert.Equal(t, model.ContentStatus("draft"), model.ContentStatusDraft)
	assert.Equal(t, model.ContentStatus("published"), model.ContentStatusPublished)
	assert.Equal(t, model.ContentStatus("pending"), model.ContentStatusPending)
}

func TestContent(t *testing.T) {
	now := time.Now().UnixMilli()
	content := &model.Content{
		ID:            1,
		Type:          model.ContentTypePost,
		Title:         "Test Post",
		Abstract:      "Test Abstract",
		Text:          "Test Content",
		TextHTML:      template.HTML("<p>Test Content</p>"),
		FeaturedImage: "test.jpg",
		AuthorId:      1,
		Author: model.User{
			ID:   1,
			Name: "Test User",
		},
		CategoryID: 1,
		Category: model.Category{
			ID:   1,
			Text: "Test Category",
		},
		CreatedAt:    now,
		UpdatedAt:    now,
		PublishAt:    now,
		Status:       model.ContentStatusPublished,
		AllowComment: true,
		Tags: []model.Tag{
			{
				ID:   1,
				No:   "test-tag",
				Text: "Test Tag",
			},
		},
		ViewsNum:    100,
		CommentsNum: 10,
		Comments: []model.Comment{
			{
				ID:   1,
				Text: "Test Comment",
			},
		},
	}

	// Test JSON marshaling
	data, err := json.Marshal(content)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"id":1`)
	assert.Contains(t, string(data), `"type":"post"`)
	assert.Contains(t, string(data), `"title":"Test Post"`)
	assert.Contains(t, string(data), `"textHtml":"\u003cp\u003eTest Content\u003c/p\u003e"`)
	assert.Contains(t, string(data), `"status":"published"`)

	// Test JSON unmarshaling
	var newContent model.Content
	err = json.Unmarshal(data, &newContent)
	assert.NoError(t, err)
	assert.Equal(t, content.ID, newContent.ID)
	assert.Equal(t, content.Type, newContent.Type)
	assert.Equal(t, content.Title, newContent.Title)
	assert.Equal(t, content.Status, newContent.Status)
	assert.Equal(t, content.AllowComment, newContent.AllowComment)
	assert.Equal(t, content.ViewsNum, newContent.ViewsNum)
	assert.Equal(t, content.CommentsNum, newContent.CommentsNum)
}

func TestContentMethods(t *testing.T) {
	now := time.Now()
	content := &model.Content{
		ID:        1,
		Title:     "Test Post",
		Abstract:  "",
		Text:      "This is a test content that should be truncated for the abstract. We need to make sure it works correctly with the GetAbstract method.",
		CreatedAt: now.UnixMilli(),
		PublishAt: now.UnixMilli(),
		Author:    model.User{ID: 1, Name: "Test User"},
	}

	// Test GetAuthor
	author := content.GetAuthor()
	assert.Equal(t, uint(1), author.ID)
	assert.Equal(t, "Test User", author.Name)

	// Test GetAbstract
	abstract := content.GetAbstract()
	assert.NotEmpty(t, abstract)
	assert.LessOrEqual(t, len(abstract), 140)

	// Test GetDate
	date := content.GetDate()
	assert.NotEmpty(t, date)

	// Test GetTime
	timeStr := content.GetTime()
	assert.NotEmpty(t, timeStr)

	// Test GetDateAndTime
	dateTime := content.GetDateAndTime()
	assert.NotEmpty(t, dateTime)
}
