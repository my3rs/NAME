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
	log.Println("Connecting to sqlite database in " + config.Database.DataPath + "/name.db")
	db, err := gorm.Open(sqlite.Open(config.Database.DataPath+"/name.db"), &gorm.Config{})
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
