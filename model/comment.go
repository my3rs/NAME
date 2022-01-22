package model

import (
	"gorm.io/gorm"
)

const (
	CommentStatus_Normal     = 0 // 正常评论
	CommentStatus_Unreviewed = 1 // 等待审核
	CommentStatus_Refused    = 2 // 审核未通过
	CommentStatus_Trash      = 3 // 放入垃圾箱
)

type CommentDb struct {
	gorm.Model
	Id         uint
	CreatedAt  int
	AuthorId   uint
	AuthorName string
	Path       string
	Mail       string
	Url        string
	Ip         string
	Text       string
	Status     uint
}

func GetCommentById(id uint) (*CommentDb, bool) {
	var comment CommentDb

	if err := Db.First(&comment, id); err != nil {
		return &comment, true
	}

	return nil, false
}
