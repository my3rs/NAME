package database

import (
	"NAME/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDB(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	db := getTestDB(t)
	assert.NotNil(t, db, "Database connection should not be nil")
}

func TestUpdateSequence(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	db := getTestDB(t)

	// Create test data
	categories := []model.Category{
		{ID: 1, Text: "Category 1", No: "cat-1"},
		{ID: 3, Text: "Category 2", No: "cat-2"}, // Skip ID 2 intentionally
		{ID: 5, Text: "Category 3", No: "cat-3"}, // Skip ID 4 intentionally
	}

	// Insert test data
	err := db.Create(&categories).Error
	require.NoError(t, err)

	// Test UpdateSequence
	err = db.Exec("SELECT setval('categories_id_seq', ?);", 5).Error
	assert.NoError(t, err)

	// Verify sequence is updated by inserting a new record
	newCategory := model.Category{
		Text: "New Category",
		No:   "cat-new",
	}
	err = db.Create(&newCategory).Error
	assert.NoError(t, err)
	assert.Equal(t, uint(6), newCategory.ID, "New record should have ID 6")
}

func TestDatabaseOperations(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	db := getTestDB(t)

	// Test Create
	category := model.Category{
		Text: "Test Category",
		No:   "test-cat",
	}
	err := db.Create(&category).Error
	assert.NoError(t, err)
	assert.NotZero(t, category.ID)

	// Test Read
	var found model.Category
	err = db.First(&found, category.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, category.Text, found.Text)
	assert.Equal(t, category.No, found.No)

	// Test Update
	found.Text = "Updated Category"
	err = db.Save(&found).Error
	assert.NoError(t, err)

	var updated model.Category
	err = db.First(&updated, category.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "Updated Category", updated.Text)

	// Test Delete
	err = db.Delete(&updated).Error
	assert.NoError(t, err)

	var notFound model.Category
	err = db.First(&notFound, category.ID).Error
	assert.Error(t, err, "Record should be deleted")
}

func TestDatabaseRelationships(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	db := getTestDB(t)

	// Create test data
	user := model.User{
		Name: "Test User",
	}
	err := db.Create(&user).Error
	require.NoError(t, err)

	category := model.Category{
		Text: "Test Category",
		No:   "test-cat",
	}
	err = db.Create(&category).Error
	require.NoError(t, err)

	content := model.Content{
		Title:      "Test Content",
		Text:       "Test Content Body",
		AuthorId:   user.ID,
		CategoryID: category.ID,
	}
	err = db.Create(&content).Error
	require.NoError(t, err)

	comment := model.Comment{
		ContentID:  content.ID,
		AuthorName: "Commenter",
		Text:       "Test Comment",
		Status:     uint(model.CommentStatus_Approved),
	}
	err = db.Create(&comment).Error
	require.NoError(t, err)

	// Test loading relationships
	var foundContent model.Content
	err = db.Preload("Author").Preload("Category").Preload("Comments").First(&foundContent, content.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, user.ID, foundContent.Author.ID)
	assert.Equal(t, category.ID, foundContent.Category.ID)
	assert.Equal(t, 1, len(foundContent.Comments))
	assert.Equal(t, comment.Text, foundContent.Comments[0].Text)
}
