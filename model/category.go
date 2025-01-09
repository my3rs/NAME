package model

type Category struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Text string `json:"text"`
	No   string `json:"no"`

	CreatedAt int `json:"createdAt"  gorm:"autoCreateTime"`
}

func (Category) TableName() string {
	return "categories"
}
