package model

type Setting struct {
	ID    uint   `gorm:"primaryKey;autoIncrement" json:"-"`
	Key   string `gorm:"uniqueIndex;size:50;not null" json:"key"`
	Value string `gorm:"size:1000;not null" json:"value"`
}

const (
	EnvironmentDev  = "dev"
	EnvironmentProd = "prod"
)

func (Setting) TableName() string {
	return "settings"
}
