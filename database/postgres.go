package database

import (
	"NAME/conf"
	"NAME/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
	"sync"
)

var (
	once sync.Once
	db   *gorm.DB
)

func initDb() (*gorm.DB, error) {
	once.Do(func() {
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

		db.AutoMigrate(&model.Content{}, &model.User{}, &model.Tag{})
	})

	return db, nil
}

func GetDb() (*gorm.DB, error) {
	if db == nil {
		return initDb()
	} else {
		return db, nil
	}
}
