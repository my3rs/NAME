package model

import (
	"change/conf"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"html/template"
	"strings"
)

var Db *gorm.DB

func Init() {
	// dsn example : host=127.0.0.1 user=postgres password=postgres dbname=nuwa post=5432
	dsn := "host=" + strings.Split(conf.Config().DB.Host, ":")[0] +
		" user=" + conf.Config().DB.User +
		" password= " + conf.Config().DB.Password +
		" dbname=" + conf.Config().DB.Name +
		" port=" + strings.Split(conf.Config().DB.Host, ":")[1] +
		" sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}

	db.AutoMigrate(&Content{})
	db.AutoMigrate(&Post{})
	db.AutoMigrate(&User{})

	Db = db
}

func Markdown2Html(markdown string) template.HTML {
	unsafe := blackfriday.Run([]byte(markdown))
	html := template.HTML(bluemonday.UGCPolicy().SanitizeBytes(unsafe))

	return html
}
