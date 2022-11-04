package model

type User struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	Name           string `gorm:"unique" json:"name"`
	HashedPassword string `gorm:"column:password" json:"-"`
	Mail           string `gorm:"unique" json:"mail"`
	Url            string `json:"url"`
	CreatedAt      int    `json:"createdAt"`
	UpdatedAt      int    `json:"updatedAt"`
	Activated      bool   `json:"activated"`
	Role           string `json:"role"`
}
