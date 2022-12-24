package service

import (
	"NAME/database"
	"gorm.io/gorm"
	"log"
)

type Uploader interface {
}

type uploader struct {
	DB *gorm.DB
}

func NewUploader() Uploader {
	db, err := database.GetDb()
	if err != nil {
		log.Panic(err.Error())
	}

	return &uploader{DB: db}
}
