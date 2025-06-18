package model

import (
	"NAME/model"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_TableName(t *testing.T) {
	user := model.User{}
	assert.Equal(t, "users", user.TableName())
}

func TestUser_JSONSerialization(t *testing.T) {
	user := model.User{
		ID:             1,
		Username:       "testuser",
		HashedPassword: "hashedpassword123",
		Mail:          "test@example.com",
		Avatar:        "https://example.com/avatar.jpg",
		URL:           "https://example.com",
		Role:          model.UserRole("admin"),
		Activated:     true,
	}

	// 测试序列化
	jsonData, err := json.Marshal(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// 验证敏感字段不被序列化
	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonData, &jsonMap)
	assert.NoError(t, err)
	
	// HashedPassword字段应该被忽略
	_, exists := jsonMap["hashedpassword"]
	assert.False(t, exists)
	_, exists = jsonMap["HashedPassword"]
	assert.False(t, exists)

	// 其他字段应该存在
	assert.Equal(t, float64(1), jsonMap["id"])
	assert.Equal(t, "testuser", jsonMap["username"])
	assert.Equal(t, "test@example.com", jsonMap["mail"])
	assert.Equal(t, "admin", jsonMap["role"])
	assert.Equal(t, true, jsonMap["activated"])
}

func TestUser_JSONDeserialization(t *testing.T) {
	jsonStr := `{
		"id": 1,
		"username": "testuser",
		"mail": "test@example.com",
		"avatar": "https://example.com/avatar.jpg",
		"url": "https://example.com",
		"role": "user",
		"activated": true
	}`

	var user model.User
	err := json.Unmarshal([]byte(jsonStr), &user)
	assert.NoError(t, err)

	assert.Equal(t, uint(1), user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Mail)
	assert.Equal(t, "https://example.com/avatar.jpg", user.Avatar)
	assert.Equal(t, "https://example.com", user.URL)
	assert.Equal(t, model.UserRole("user"), user.Role)
	assert.True(t, user.Activated)
}

func TestUserRole_Values(t *testing.T) {
	adminRole := model.UserRole("admin")
	userRole := model.UserRole("user")

	assert.Equal(t, "admin", string(adminRole))
	assert.Equal(t, "user", string(userRole))
}

func TestUser_EmptyUser(t *testing.T) {
	user := model.User{}
	
	// 测试空用户的默认值
	assert.Equal(t, uint(0), user.ID)
	assert.Equal(t, "", user.Username)
	assert.Equal(t, "", user.Mail)
	assert.Equal(t, "", user.Avatar)
	assert.Equal(t, "", user.URL)
	assert.Equal(t, model.UserRole(""), user.Role)
	assert.False(t, user.Activated)
}

func TestUser_ValidRoles(t *testing.T) {
	testCases := []struct {
		role     model.UserRole
		expected string
	}{
		{model.UserRole("admin"), "admin"},
		{model.UserRole("user"), "user"},
		{model.UserRole("moderator"), "moderator"},
		{model.UserRole(""), ""},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expected, string(tc.role))
	}
}