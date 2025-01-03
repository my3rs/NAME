package database

import (
	"NAME/conf"
	"NAME/model"
	"strings"
	"sync"

	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

// GetDB returns the database instance
func GetDB() *gorm.DB {
	once.Do(initPostgres)
	return db
}

// SetDB sets the database instance (for testing)
func SetDB(instance *gorm.DB) {
	db = instance
}

func initPostgres() {
	config := conf.GetConfig()
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.Name)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}

	// 创建 ltree 扩展
	err = db.Exec("CREATE EXTENSION IF NOT EXISTS ltree").Error
	if err != nil {
		panic("Failed to create ltree extension")
	}

	// 迁移数据库
	err = db.AutoMigrate(
		&model.Attachment{},
		&model.Comment{},
		&model.Content{},
		&model.User{},
		&model.Tag{},
		&model.Setting{},
		&model.Category{},
	)
	if err != nil {
		panic("Failed to migrate database")
	}
}

// UpdateSequence updates the sequence of a table
func UpdateSequence(tableName string) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	// Get the current maximum ID
	var maxID uint
	result := db.Table(tableName).Select("COALESCE(MAX(id), 0)").Scan(&maxID)
	if result.Error != nil {
		return result.Error
	}

	// Update the sequence
	seqName := strings.ToLower(tableName) + "_id_seq"
	return db.Exec(fmt.Sprintf("SELECT setval('%s', %d)", seqName, maxID)).Error
}
