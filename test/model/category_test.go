package model_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"NAME/model"
)

func TestCategory(t *testing.T) {
	category := &model.Category{
		ID:   1,
		Text: "Test Category",
		Slug: "test-category",
	}

	// Test JSON marshaling
	data, err := json.Marshal(category)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"id":1`)
	assert.Contains(t, string(data), `"text":"Test Category"`)
	assert.Contains(t, string(data), `"slug":"test-category"`)

	// Test JSON unmarshaling
	var newCategory model.Category
	err = json.Unmarshal(data, &newCategory)
	assert.NoError(t, err)
	assert.Equal(t, category.ID, newCategory.ID)
	assert.Equal(t, category.Text, newCategory.Text)
	assert.Equal(t, category.Slug, newCategory.Slug)

	// Test TableName
	assert.Equal(t, "categories", category.TableName())
}
