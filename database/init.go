package database

import "NAME/model"

func init() {
	initPostgres()
	model.LoadSettingsToCache(db)
}
