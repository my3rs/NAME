package model

import "gorm.io/gorm"

type Setting struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Key   string `gorm:"uniqueIndex;size:50;not null"`
	Value string `gorm:"size:1000;not null"`
}

var settingsMap = make(map[string]Setting)

func (Setting) TableName() string {
	return "settings"
}

func GetSettingsItem(key string) (Setting, bool) {
	setting, found := settingsMap[key]
	return setting, found
}

// LoadSettingsToCache 将数据库中的配置项读取到内存中
func LoadSettingsToCache(db *gorm.DB) {
	var settings []Setting
	result := db.Find(&settings)

	if result.Error != nil {
		panic("读取配置项失败")
	}

	for _, setting := range settings {
		settingsMap[setting.Key] = setting
	}
}
