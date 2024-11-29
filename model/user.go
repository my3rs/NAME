package model

type UserRole string

const (
	UserRoleAdmin  UserRole = "admin"
	UserRoleReader UserRole = "reader"
)

type User struct {
	ID             uint     `gorm:"primaryKey" json:"id"`
	Name           string   `gorm:"unique" json:"name"`
	HashedPassword string   `gorm:"column:password" json:"-"`
	Mail           string   `gorm:"unique" json:"mail"`
	Avatar         string   `json:"avatar"`
	Url            string   `json:"url"`
	CreatedAt      int64    `json:"createdAt" gorm:"autoCreateTime:milli"`
	UpdatedAt      int64    `json:"updatedAt" gorm:"autoUpdateTime:milli"`
	Activated      bool     `json:"activated"`
	Role           UserRole `json:"role"`
}
