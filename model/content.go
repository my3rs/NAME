package model

import (
	"html/template"
	"time"
)

type ContentType string

const (
	ContentTypePost ContentType = "post"
	ContentTypeDigu ContentType = "digu"
	ContentTypePage ContentType = "page"
)

type ContentStatus string

const (
	ContentStatusDraft     ContentStatus = "draft"
	ContentStatusPublished ContentStatus = "published"
	ContentStatusPending   ContentStatus = "pending"
)

type Content struct {
	ID            uint          `gorm:"primaryKey" json:"id"`
	Type          ContentType   `json:"type"`
	Title         string        `json:"title"`
	Abstract      string        `json:"abstract"`
	Text          string        `json:"text"`
	TextHTML      template.HTML `json:"textHtml" gorm:"-"`
	FeaturedImage string        `json:"featuredImage" gorm:"default:null"`

	AuthorId uint `json:"-"`
	Author   User `json:"author"`

	CategoryID uint     `json:"-"`
	Category   Category `json:"category"`

	CreatedAt    int64         `json:"createdAt" gorm:"autoCreateTime:milli"`
	UpdatedAt    int64         `json:"updatedAt" gorm:"autoUpdateTime:milli"`
	PublishAt    int64         `json:"publishAt"`
	Status       ContentStatus `json:"status"`
	AllowComment bool          `json:"allowComment"`
	Password     string        `json:"-"`

	Tags []Tag `json:"tags" gorm:"many2many:content_tags;constraint:OnDelete:CASCADE"`

	ViewsNum    uint      `json:"viewsNum"`
	CommentsNum uint      `json:"commentsNum"`
	Comments    []Comment `json:"comments" gorm:"constraint:OnDelete:CASCADE"`
}

func (c *Content) GetAuthor() User {
	return c.Author
}

func (c *Content) GetAbstract() string {
	if len(c.Abstract) > 0 {
		return c.Abstract
	}

	maxLen := 140
	if len(c.Text) < maxLen {
		return c.Text
	}
	return c.Text[0:maxLen]
}

func (c *Content) GetDate() string {
	date := time.Unix(c.CreatedAt, 0).Format("2006-01-02")
	return date
}

func (c *Content) GetTime() string {
	t := time.Unix(c.CreatedAt, 0).Format("15:04:05")
	return t
}

func (c *Content) GetDateAndTime() string {
	t := time.Unix(c.CreatedAt, 0).Format("2006-01-02 15:04:05")
	return t
}
