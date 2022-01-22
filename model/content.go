package model

import (
	"fmt"
	"gorm.io/gorm"
	"html/template"
	"time"
)

const (
	ContentTypePost = 0
	ContentTypeDigu = 1
	ContentTypePage = 2
)

const (
	ContentStatusDraft = 0
)

type Content struct {
	gorm.Model
	ID           uint
	Type         uint
	Title        string
	Abstract     string
	Text         string
	AuthorId     uint
	TemplateId   uint
	CreatedAt    int64
	UpdatedAt    int64
	Status       uint
	IsPublic     bool
	AllowComment bool
	Password     string

	ViewsNum    uint
	CommentsNum uint
}

func (c *Content) GetAuthor() (*User, bool) {
	var user User

	if err := Db.First(&user, c.ID); err != nil {
		return &user, false
	}

	return &user, true
}

func (c *Content) GetAbstract() string {
	if len(c.Abstract) == 0 {
		return c.Text[0:140]
	}

	return c.Abstract
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

func (c *Content) GetAbstractHtml() template.HTML {
	md := c.GetAbstract()

	return Markdown2Html(md)
}

// GetContentByPage 根据页码和每页数量返回Content Slice
func GetContentByPage(page int, count int) []Content {
	if page <= 0 || count >= 100 {
		return []Content{}
	}

	leftIndex := (page - 1) * count
	rightIndex := page * count

	allContent := GetContentSliceOrderedByCreatedTime()

	if len(allContent) < leftIndex-1 ||
		len(allContent) == 0 {
		return []Content{}
	}

	return allContent[leftIndex:rightIndex]
}

// GetContentSliceOrderedByCreatedTime 获得按”创建时间“降序排列的Content Slice
func GetContentSliceOrderedByCreatedTime() []Content {
	var list []Content
	if result := Db.Order("created_at desc").Find(&list); result.Error != nil {
		fmt.Println(result.Error)
	}

	return list
}

func GetContentSlice() []Content {
	var list []Content

	if result := Db.Find(&list); result.Error != nil {
		fmt.Println(result.Error)
	}

	return list
}

func GetContentById(id int64) (*Content, bool) {
	var content Content

	if err := Db.First(&content, id); err != nil {
		return &content, false
	}

	return &content, true
}

func SaveContent(content Content) bool {
	if (content == Content{}) {
		return false
	}

	if result := Db.Create(&content); result.Error != nil {
		return false
	}

	return true
}
