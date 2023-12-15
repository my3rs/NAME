package service

import (
	"NAME/database"
	"NAME/model"
	"gorm.io/gorm"
	"log"
)

type AttachmentService interface {
	InsertAttachment(attachment model.Attachment) error
}

type attachmentService struct {
	DB *gorm.DB
}

func NewAttachmentService() AttachmentService {
	db, err := database.GetDb()
	if err != nil {
		log.Panic(err.Error())
	}

	return &attachmentService{DB: db}
}

func (u *attachmentService) InsertAttachment(attachment model.Attachment) error {
	log.Printf("InsertAttachment: %+v\n", attachment)
	result := u.DB.Create(&attachment)

	return result.Error
}
