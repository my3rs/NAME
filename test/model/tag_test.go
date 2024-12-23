package model_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"NAME/model"
)

func TestTag(t *testing.T) {
	tag := &model.Tag{
		ID:       1,
		ParentID: 0,
		No:       "test-tag",
		Text:     "Test Tag",
		Path:     "1",
	}

	// Test JSON marshaling
	data, err := json.Marshal(tag)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"id":1`)
	assert.Contains(t, string(data), `"no":"test-tag"`)
	assert.Contains(t, string(data), `"text":"Test Tag"`)

	// Test JSON unmarshaling
	var newTag model.Tag
	err = json.Unmarshal(data, &newTag)
	assert.NoError(t, err)
	assert.Equal(t, tag.ID, newTag.ID)
	assert.Equal(t, tag.ParentID, newTag.ParentID)
	assert.Equal(t, tag.No, newTag.No)
	assert.Equal(t, tag.Text, newTag.Text)
	assert.Equal(t, tag.Path, newTag.Path)
}

func TestTagExt(t *testing.T) {
	tagExt := &model.TagExt{
		Tag: model.Tag{
			ID:       1,
			ParentID: 0,
			No:       "test-tag",
			Text:     "Test Tag",
			Path:     "1",
		},
		ReadablePath: "Root/Test Tag",
	}

	// Test JSON marshaling
	data, err := json.Marshal(tagExt)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"id":1`)
	assert.Contains(t, string(data), `"no":"test-tag"`)
	assert.Contains(t, string(data), `"text":"Test Tag"`)
	assert.Contains(t, string(data), `"readablePath":"Root/Test Tag"`)

	// Test JSON unmarshaling
	var newTagExt model.TagExt
	err = json.Unmarshal(data, &newTagExt)
	assert.NoError(t, err)
	assert.Equal(t, tagExt.ID, newTagExt.ID)
	assert.Equal(t, tagExt.ParentID, newTagExt.ParentID)
	assert.Equal(t, tagExt.No, newTagExt.No)
	assert.Equal(t, tagExt.Text, newTagExt.Text)
	assert.Equal(t, tagExt.Path, newTagExt.Path)
	assert.Equal(t, tagExt.ReadablePath, newTagExt.ReadablePath)
}
