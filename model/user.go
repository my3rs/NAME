package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID             int
	Name           string
	HashedPassword string `gorm:"column:password"`
	Mail           string `gorm:"column:email"`
	Url            string
	CreatedAt      int
	Activated      bool
	Role           string
}

// ValidatePassword
// @pwd: plain password
func (u *User) ValidatePassword(pwd []byte) bool {
	byteHash := []byte(u.HashedPassword)
	if err := bcrypt.CompareHashAndPassword(byteHash, pwd); err != nil {
		return false
	}
	return true
}

func GetUserById(id int) (*User, bool) {
	var user User

	if result := Db.First(&user, id); result.Error != nil {
		return nil, false
	}

	return &user, true
}

func GetUserByName(name string) (*User, bool) {
	var user User
	if result := Db.Where("name = ?", name).Find(&user); result.Error != nil {
		return nil, false
	}
	return &user, true
}

func GetUserByMail(mail string) (*User, bool) {
	var user User
	if result := Db.Where("email = ?", mail).Find(&user); result.Error != nil {
		return &User{}, false
	}
	return &user, true
}

func CreateUser(user User) (User, error) {
	if result := Db.Create(&user); result.Error != nil {
		return User{}, result.Error
	}

	return user, nil
}

func DeleteUserById(id int) bool {
	Db.Delete(&User{}, id)
	return true
}

func GetAll() []User {
	var users []User
	if result := Db.Find(&users); result.Error != nil {
		return []User{}
	}

	return users
}
