package service

import (
	"NAME/database"
	"NAME/model"
	"gorm.io/gorm"
	"log"
)

type Uploader interface {
	InsertAttachment(attachment model.Attachment) error
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

func (u *uploader) InsertAttachment(attachment model.Attachment) error {
	result := u.DB.Create(&attachment)

	return result.Error
}
