package database

import (
	"NAME/conf"
	"NAME/model"
	"log"
	"strings"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	once sync.Once
	db   *gorm.DB
)

func initPostgres() {

	// dsn example : host=127.0.0.1 user=postgres password=postgres dbname=nuwa post=5432
	dsn := "host=" + strings.Split(conf.Config().DB.Host, ":")[0] +
		" user=" + conf.Config().DB.User +
		" password= " + conf.Config().DB.Password +
		" dbname=" + conf.Config().DB.Name +
		" port=" + strings.Split(conf.Config().DB.Host, ":")[1] +
		" sslmode=disable TimeZone=Asia/Shanghai"

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}

	err = db.AutoMigrate(
		&model.Attachment{},
		&model.Comment{},
		&model.Content{},
		&model.User{},
		&model.Tag{},
		&model.Setting{},
		&model.Category{},
	)
	if err != nil {
		log.Println("数据库迁移失败：", err)
	}
}

func GetDB() *gorm.DB {
	if db == nil {
		panic("数据库连接为 nil")
	}

	return db
}
