package service

import (
	"NAME/database"
	"NAME/model"
	"NAME/service"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var testUserIDs []uint

func SetupTestDB(t *testing.T) *gorm.DB {
	db := database.GetDB()
	testUserIDs = []uint{} // 重置测试用户ID列表
	return db
}

func CleanupTestUsers(db *gorm.DB) error {
	if len(testUserIDs) > 0 {
		return db.Delete(&model.User{}, testUserIDs).Error
	}
	return nil
}

func CreateTestUser(t *testing.T, userService service.UserService) model.User {
	timestamp := time.Now().UnixNano() // 使用纳秒时间戳确保唯一性
	time.Sleep(1 * time.Millisecond) // 确保时间戳唯一
	user := model.User{
		Username:         fmt.Sprintf("test_user_%d", timestamp),
		HashedPassword: "$2a$10$Nmz0WVEsuNT69cljMkb25.ASmyIVHL3vTLy9lZLQOEBRfHYW2I4HC", // 密码: test123
		Mail:           fmt.Sprintf("test_%d@example.com", timestamp),
		Avatar:         "https://example.com/avatar.jpg",
		URL:            "https://example.com",
		Role:           model.UserRole("admin"),
		Activated:      true,
	}

	err := userService.InsertUser(user)
	assert.NoError(t, err)

	// 获取插入后的用户（包含ID）
	insertedUser, err := userService.GetUserByName(user.Username)
	assert.NoError(t, err)
	testUserIDs = append(testUserIDs, insertedUser.ID)
	return insertedUser
}
