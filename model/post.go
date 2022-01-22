package model

import (
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
	"gorm.io/gorm"
	"html/template"
	"time"
)

type Post struct {
	gorm.Model
	Id           uint
	Title        string
	Content      string
	AuthorId     uint
	CreatedAt    int64
	UpdatedAt    int64
	IsPublic     bool
	AllowComment bool

	Views int
}

// PostHtml 在Post基础上，针对Html更改了数据类型
type PostHtml struct {
	Id           uint
	Title        string
	ContentHtml  template.HTML
	AuthorId     uint
	CreateAt     string
	UpdateAt     string
	IsPublic     bool
	AllowComment bool
	Views        int
}

func GetPostById(id int64) (*Post, bool) {
	var post Post

	if err := Db.First(&post, id); err != nil {
		return &post, true
	}

	return nil, false
}

func Post2PostHtml(post *Post) *PostHtml {
	var htmlPost PostHtml

	// 复制基本数据
	htmlPost.Id = post.Id
	htmlPost.Title = post.Title
	htmlPost.AuthorId = post.AuthorId
	htmlPost.IsPublic = post.IsPublic
	htmlPost.AllowComment = post.AllowComment
	htmlPost.Views = post.Views

	// 转换时间戳为字符串
	ctm := time.Unix(post.CreatedAt, 0)
	utm := time.Unix(post.UpdatedAt, 0)

	htmlPost.CreateAt = ctm.Format("2006-01-02 15:04:05")
	htmlPost.UpdateAt = utm.Format("2006-01-02 15:04:05")

	// Markdown 转为 Html
	//unsafeHtml := bluemonday.UGCPolicy().Sanitize(post.Content)
	//htmlPost.ContentHtml = string(blackfriday.Run([]byte(unsafeHtml)))
	unsafe := blackfriday.Run([]byte(post.Content))
	htmlPost.ContentHtml = template.HTML(bluemonday.UGCPolicy().SanitizeBytes(unsafe))
	fmt.Println(string(unsafe))

	return &htmlPost
}
