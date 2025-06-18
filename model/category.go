package model

type Category struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Text      string `json:"text"`
	Slug      string `json:"slug"`
	CreatedAt int64  `json:"createdAt,omitempty" gorm:"autoCreateTime"`
}

func (Category) TableName() string {
	return "categories"
}
