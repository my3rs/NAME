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
	ID            uint          `gorm:"primaryKey;comment:内容ID" json:"id"`
	Slug          string        `json:"slug" gorm:"comment:短标题"`
	Type          ContentType   `json:"type" gorm:"comment:内容类型：post文章、digu嘀咕、page页面"`
	Title         string        `json:"title" gorm:"comment:内容标题"`
	Abstract      string        `json:"abstract" gorm:"comment:内容摘要，如果为空则自动截取正文"`
	Text          string        `json:"text" gorm:"comment:内容正文，Markdown格式"`
	TextHTML      template.HTML `json:"textHtml" gorm:"-;comment:内容的HTML格式，运行时生成"`
	FeaturedImage string        `json:"featuredImage" gorm:"default:null;comment:特色图片URL"`

	AuthorId uint `json:"-" gorm:"comment:作者ID"`
	Author   User `json:"author" gorm:"comment:作者信息"`

	CategoryID uint     `json:"-" gorm:"comment:分类ID"`
	Category   Category `json:"category" gorm:"comment:分类信息"`

	CreatedAt    int64         `json:"createdAt" gorm:"autoCreateTime;comment:创建时间（毫秒时间戳）"`
	UpdatedAt    int64         `json:"updatedAt" gorm:"autoUpdateTime;comment:更新时间（毫秒时间戳）"`
	PublishAt    int64         `json:"publishAt" gorm:"comment:发布时间（毫秒时间戳）"`
	Status       ContentStatus `json:"status" gorm:"comment:内容状态：draft草稿、published已发布、pending待审核"`
	AllowComment bool          `json:"allowComment" gorm:"comment:是否允许评论"`
	Password     string        `json:"-" gorm:"comment:访问密码，为空则不需要密码"`

	Tags []Tag `json:"tags" gorm:"many2many:content_tags;constraint:OnDelete:CASCADE;comment:关联的标签列表"`

	ViewsNum    uint      `json:"viewsNum" gorm:"comment:浏览次数"`
	CommentsNum uint      `json:"commentsNum" gorm:"comment:评论数量"`
	Comments    []Comment `json:"comments" gorm:"constraint:OnDelete:CASCADE;comment:关联的评论列表"`
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
