package model

type Category struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

func (Category) TableName() string {
	return "categories"
}
