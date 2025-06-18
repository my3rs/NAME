package service

import (
	"NAME/database"
	"NAME/model"
	"NAME/service"
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
	user := model.User{
		Username:         "test_user_" + time.Now().Format("20060102150405"),
		HashedPassword: "$2a$10$IVxZxP.b7Ey9VPdCkKV4UeHHdZhafQ7x2qFGFMpVCGbfFtCHBZnGK", // 密码: test123
		Mail:           "test_" + time.Now().Format("20060102150405") + "@example.com",
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
