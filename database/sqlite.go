package database

import (
	"NAME/conf"
	"NAME/model"
	// "gorm.io/driver/sqlite"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"log"
)

func initSQLite() *gorm.DB {
	config := conf.GetConfig()
	dbPath := config.Database.DataPath + "/" + config.Database.FileName
	log.Println("Connecting to sqlite database in " + dbPath)
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(
		&model.User{},
		&model.Attachment{},
		&model.Comment{},
		&model.Content{},
		&model.Tag{},
		&model.Setting{},
		&model.Category{},
	)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
