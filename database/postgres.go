package database

import (
	"NAME/conf"
	"NAME/model"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initPostgres() *gorm.DB {
	config := conf.GetConfig()
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.Name,
		config.Database.SSLMode,
		config.Database.TimeZone)

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

	return db
}
