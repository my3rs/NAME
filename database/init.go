package database

import (
	"NAME/conf"
	"sync"

	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

// GetDB returns the database instance
func GetDB() *gorm.DB {
	once.Do(func() {
		config := conf.GetConfig()
		if config.Database.Driver == conf.DATABASE_DRIVER_SQLITE {
			db = initSQLite()
		} else if config.Database.Driver == conf.DATABASE_DRIVER_POSTGRES {
			db = initPostgres()
		}

		// Settings are now loaded on demand by services
	})

	return db
}

//func init() {
//	config := conf.GetConfig()
//
//	if config.Database.Driver == conf.DATABASE_DRIVER_SQLITE {
//		db = initSQLite()
//	} else if config.Database.Driver == conf.DATABASE_DRIVER_POSTGRES {
//		db = initPostgres()
//	}
//
//	model.LoadSettingsToCache(db)
//}
