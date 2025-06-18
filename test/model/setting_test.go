package model

import (
	"NAME/model"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetting_TableName(t *testing.T) {
	setting := model.Setting{}
	assert.Equal(t, "settings", setting.TableName())
}

func TestSetting_JSONSerialization(t *testing.T) {
	setting := model.Setting{
		ID:    1,
		Key:   "site_title",
		Value: "My Blog",
	}

	// 测试序列化
	jsonData, err := json.Marshal(setting)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// 验证JSON结构
	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonData, &jsonMap)
	assert.NoError(t, err)

	// ID字段应该被忽略（json:"-"）
	_, exists := jsonMap["id"]
	assert.False(t, exists)
	_, exists = jsonMap["ID"]
	assert.False(t, exists)

	// 其他字段应该存在
	assert.Equal(t, "site_title", jsonMap["key"])
	assert.Equal(t, "My Blog", jsonMap["value"])
}

func TestSetting_JSONDeserialization(t *testing.T) {
	jsonStr := `{
		"key": "site_description",
		"value": "A wonderful blog"
	}`

	var setting model.Setting
	err := json.Unmarshal([]byte(jsonStr), &setting)
	assert.NoError(t, err)

	assert.Equal(t, "site_description", setting.Key)
	assert.Equal(t, "A wonderful blog", setting.Value)
	assert.Equal(t, uint(0), setting.ID) // ID不会从JSON反序列化
}

func TestSetting_Constants(t *testing.T) {
	assert.Equal(t, "dev", model.EnvironmentDev)
	assert.Equal(t, "prod", model.EnvironmentProd)
}

func TestSetting_KeyValuePairs(t *testing.T) {
	// Test typical key-value configurations
	settings := []model.Setting{
		{Key: "site_title", Value: "My Blog"},
		{Key: "posts_per_page", Value: "10"},
		{Key: "allow_comments", Value: "true"},
	}
	
	for _, setting := range settings {
		assert.NotEmpty(t, setting.Key)
		assert.NotEmpty(t, setting.Value)
	}
}

func TestSetting_EmptySetting(t *testing.T) {
	setting := model.Setting{}
	
	assert.Equal(t, uint(0), setting.ID)
	assert.Equal(t, "", setting.Key)
	assert.Equal(t, "", setting.Value)
}

func TestSetting_ValidData(t *testing.T) {
	testCases := []struct {
		key   string
		value string
	}{
		{"site_title", "My Blog"},
		{"site_description", "A wonderful blog"},
		{"theme", "default"},
		{"posts_per_page", "10"},
		{"allow_comments", "true"},
	}

	for _, tc := range testCases {
		setting := model.Setting{
			Key:   tc.key,
			Value: tc.value,
		}
		
		assert.Equal(t, tc.key, setting.Key)
		assert.Equal(t, tc.value, setting.Value)
	}
}