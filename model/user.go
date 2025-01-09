package model

import "github.com/kataras/iris/v12/middleware/jwt"

type UserRole string

const (
	UserRoleAdmin  UserRole = "admin"
	UserRoleReader UserRole = "reader"
)

func (u UserRole) String() string {
	return string(u)
}

type User struct {
	ID             uint     `gorm:"primaryKey" json:"id"`
	Username       string   `gorm:"unique" json:"username"`
	Password       string   `gorm:"-" json:"password,omitempty"`
	HashedPassword string   `gorm:"column:password" json:"-"`
	Mail           string   `gorm:"unique" json:"mail"`
	Avatar         string   `json:"avatar"`
	Url            string   `json:"url"`
	CreatedAt      int64    `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt      int64    `json:"updatedAt" gorm:"autoUpdateTime"`
	Activated      bool     `json:"activated"`
	Role           UserRole `json:"role"`
}

// Claims is a custom JWT claims struct.

type Claims struct {
	jwt.Claims
	Role UserRole `json:"role"`
}
