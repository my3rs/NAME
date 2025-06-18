package model_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"NAME/model"
)

func TestTag_JSONSerialization(t *testing.T) {
	tag := &model.Tag{
		ID:         1,
		Slug:       "test-tag",
		Text:       "Test Tag",
		UseCount: 5,
		CreatedAt:  1640995200000,
		UpdatedAt:  1640995300000,
	}

	// Test JSON marshaling
	data, err := json.Marshal(tag)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"id":1`)
	assert.Contains(t, string(data), `"slug":"test-tag"`)
	assert.Contains(t, string(data), `"text":"Test Tag"`)
	assert.Contains(t, string(data), `"useCount":5`)

	// Test JSON unmarshaling
	var newTag model.Tag
	err = json.Unmarshal(data, &newTag)
	assert.NoError(t, err)
	assert.Equal(t, tag.ID, newTag.ID)
	assert.Equal(t, tag.Slug, newTag.Slug)
	assert.Equal(t, tag.Text, newTag.Text)
	assert.Equal(t, tag.UseCount, newTag.UseCount)
}

func TestTag_JSONDeserialization(t *testing.T) {
	jsonStr := `{
		"id": 2,
		"slug": "technology",
		"text": "Technology",
		"useCount": 10,
		"createdAt": 1640995200000,
		"updatedAt": 1640995300000
	}`

	var tag model.Tag
	err := json.Unmarshal([]byte(jsonStr), &tag)
	assert.NoError(t, err)

	assert.Equal(t, uint(2), tag.ID)
	assert.Equal(t, "technology", tag.Slug)
	assert.Equal(t, "Technology", tag.Text)
	assert.Equal(t, 10, tag.UseCount)
	assert.Equal(t, int64(1640995200000), tag.CreatedAt)
	assert.Equal(t, int64(1640995300000), tag.UpdatedAt)
}

func TestTag_EmptyTag(t *testing.T) {
	tag := model.Tag{}
	
	assert.Equal(t, uint(0), tag.ID)
	assert.Equal(t, "", tag.Slug)
	assert.Equal(t, "", tag.Text)
	assert.Equal(t, 0, tag.UseCount)
	assert.Equal(t, int64(0), tag.CreatedAt)
	assert.Equal(t, int64(0), tag.UpdatedAt)
}

func TestTag_ValidData(t *testing.T) {
	testCases := []struct {
		slug string
		text string
		count int
	}{
		{"technology", "Technology", 15},
		{"programming", "Programming", 25},
		{"web-development", "Web Development", 30},
		{"go-lang", "Go Language", 20},
		{"database", "Database", 18},
	}

	for _, tc := range testCases {
		tag := model.Tag{
			Slug:       tc.slug,
			Text:       tc.text,
			UseCount: tc.count,
		}
		
		assert.Equal(t, tc.slug, tag.Slug)
		assert.Equal(t, tc.text, tag.Text)
		assert.Equal(t, tc.count, tag.UseCount)
	}
}