package service

import (
	"NAME/model"
	"NAME/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserService_InsertUser(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestUsers(db)

	userService := service.NewUserService()
	
	user := model.User{
		Username:       "test_insert_" + time.Now().Format("20060102150405"),
		HashedPassword: "$2a$10$Nmz0WVEsuNT69cljMkb25.ASmyIVHL3vTLy9lZLQOEBRfHYW2I4HC",
		Mail:          "test_insert_" + time.Now().Format("20060102150405") + "@example.com",
		Avatar:        "https://example.com/avatar.jpg",
		URL:           "https://example.com",
		Role:          model.UserRole("user"),
		Activated:     true,
	}

	err := userService.InsertUser(user)
	assert.NoError(t, err)

	// 验证用户是否成功插入
	insertedUser, err := userService.GetUserByName(user.Username)
	assert.NoError(t, err)
	assert.Equal(t, user.Username, insertedUser.Username)
	assert.Equal(t, user.Mail, insertedUser.Mail)
	testUserIDs = append(testUserIDs, insertedUser.ID)
}

func TestUserService_GetUserByName(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestUsers(db)

	userService := service.NewUserService()
	user := CreateTestUser(t, userService)

	// 测试获取用户
	foundUser, err := userService.GetUserByName(user.Username)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, user.Username, foundUser.Username)
	assert.Equal(t, user.Mail, foundUser.Mail)
}

func TestUserService_GetUserByID(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestUsers(db)

	userService := service.NewUserService()
	user := CreateTestUser(t, userService)

	// 测试通过ID获取用户
	foundUser, err := userService.GetUserByID(int(user.ID))
	assert.NoError(t, err)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, user.Username, foundUser.Username)
}

func TestUserService_UpdateUser(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestUsers(db)

	userService := service.NewUserService()
	user := CreateTestUser(t, userService)

	// 测试更新用户
	user.Avatar = "https://example.com/new-avatar.jpg"
	user.URL = "https://new-example.com"
	err := userService.UpdateUser(user)
	assert.NoError(t, err)

	// 验证更新是否成功
	updatedUser, err := userService.GetUserByID(int(user.ID))
	assert.NoError(t, err)
	assert.Equal(t, user.Avatar, updatedUser.Avatar)
	assert.Equal(t, user.URL, updatedUser.URL)
}

func TestUserService_CheckUserExist(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestUsers(db)

	userService := service.NewUserService()
	user := CreateTestUser(t, userService)

	// 通过GetUserByName检查用户是否存在
	foundUser, err := userService.GetUserByName(user.Username)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, foundUser.ID)

	// 测试不存在的用户
	_, err = userService.GetUserByName("nonexistent_user")
	assert.Error(t, err) // 应该返回错误
}

func TestUserService_ValidatePassword(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestUsers(db)

	userService := service.NewUserService()
	user := CreateTestUser(t, userService)

	// 测试密码验证（测试密码是test123）
	err := userService.VerifyPassword(user, "test123")
	assert.NoError(t, err)

	// 测试错误密码
	err = userService.VerifyPassword(user, "wrongpassword")
	assert.Error(t, err)
}

func TestUserService_GetUsersWithOrder(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestUsers(db)

	userService := service.NewUserService()
	
	// 创建多个测试用户
	for i := 0; i < 3; i++ {
		CreateTestUser(t, userService)
	}

	// 测试获取用户列表
	users, err := userService.GetUsersWithOrder(0, 10, "id desc")
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(users), 3)
}
