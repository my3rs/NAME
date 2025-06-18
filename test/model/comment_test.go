package model_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"NAME/model"
)

func TestCommentStatus(t *testing.T) {
	assert.Equal(t, model.CommentStatus(0), model.CommentStatusApproved)
	assert.Equal(t, model.CommentStatus(1), model.CommentStatusUnreviewed)
	assert.Equal(t, model.CommentStatus(2), model.CommentStatusRefused)
	assert.Equal(t, model.CommentStatus(3), model.CommentStatusTrash)
}

func TestComment(t *testing.T) {
	comment := &model.Comment{
		ID:         1,
		ContentID:  100,
		CreatedAt:  1234567890,
		AuthorID:   1,
		AuthorName: "Test User",
		ParentID:   0,
		Path:       "1",
		Mail:       "test@example.com",
		URL:        "https://example.com",
		Text:       "Test Comment",
		Status:     model.CommentStatusApproved,
		IP:         "127.0.0.1",
		Agent:      "Test Agent",
		Points:     10,
	}

	// Test JSON marshaling
	data, err := json.Marshal(comment)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"id":1`)
	assert.Contains(t, string(data), `"authorName":"Test User"`)
	assert.Contains(t, string(data), `"text":"Test Comment"`)
	assert.Contains(t, string(data), `"status":0`)
	assert.Contains(t, string(data), `"points":10`)

	// Verify that sensitive fields are not included in JSON
	assert.NotContains(t, string(data), `"IP"`)
	assert.NotContains(t, string(data), `"127.0.0.1"`)

	// Test JSON unmarshaling
	var newComment model.Comment
	err = json.Unmarshal(data, &newComment)
	assert.NoError(t, err)
	assert.Equal(t, comment.ID, newComment.ID)
	assert.Equal(t, comment.AuthorName, newComment.AuthorName)
	assert.Equal(t, comment.Text, newComment.Text)
	assert.Equal(t, comment.Status, newComment.Status)
	assert.Equal(t, comment.Points, newComment.Points)

	// Test default status
	defaultComment := &model.Comment{}
	assert.Equal(t, model.CommentStatusApproved, defaultComment.Status) // Should default to Approved
}

func TestCommentWithContent(t *testing.T) {
	comment := &model.Comment{
		ID:        1,
		ContentID: 100,
		Content: model.Content{
			ID:    100,
			Title: "Test Content",
			Type:  model.ContentTypePost,
		},
		Text: "Test Comment",
	}

	// Test JSON marshaling with content
	data, err := json.Marshal(comment)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"id":1`)
	assert.Contains(t, string(data), `"content":`)
	assert.Contains(t, string(data), `"title":"Test Content"`)

	// Test JSON unmarshaling with content
	var newComment model.Comment
	err = json.Unmarshal(data, &newComment)
	assert.NoError(t, err)
	assert.Equal(t, comment.ID, newComment.ID)
	assert.Equal(t, comment.Content.ID, newComment.Content.ID)
	assert.Equal(t, comment.Content.Title, newComment.Content.Title)
}
