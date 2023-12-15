package service

import (
	"NAME/database"
	"NAME/dict"
	"NAME/model"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	GetUserByID(int) (model.User, error)
	GetUserByName(string) (model.User, error)
	GetUserByMail(string) (model.User, error)
	GetUsersWithOrder(pageIndex int, pageSize int, order string) ([]model.User, error)
	GetUserNum() int64
	InsertUser(model.User) error
	UpdateUser(model.User) error
	DeleteUser(model.User) error
	DeleteUserById(int) error
	VerifyPassword(model.User, string) error
}

type userService struct {
	Db *gorm.DB
}

func NewUserService() UserService {
	db, err := database.GetDb()
	if err != nil {
		panic(err.Error())
	}

	return &userService{
		Db: db,
	}
}

func (s *userService) GetUserByID(id int) (model.User, error) {
	var user model.User

	if result := s.Db.First(&user, id); result.Error != nil {
		return model.User{}, result.Error
	}

	return user, nil
}

func (s *userService) GetUserByName(name string) (model.User, error) {
	var user model.User
	if result := s.Db.First(&user, "name = ?", name); result.Error != nil {
		return model.User{}, result.Error
	}
	return user, nil
}

func (s *userService) GetUserByMail(mail string) (model.User, error) {
	var user model.User
	if result := s.Db.First(&user, "mail = ?", mail); result.Error != nil {
		return model.User{}, result.Error
	}
	return user, nil
}

func (s *userService) GetUserNum() int64 {
	var count int64
	s.Db.Model(&model.User{}).Count(&count)

	return count
}

func (s *userService) GetUsersWithOrder(pageIndex int, pageSize int, order string) ([]model.User, error) {
	if pageSize <= 0 || pageSize >= 100 || pageIndex < 0 {
		return nil, dict.ErrInvalidParameters
	}

	var results []model.User
	s.Db.Order(order).Limit(pageSize).Offset(pageIndex * pageSize).Find(&results)

	return results, nil
}

func (s *userService) InsertUser(user model.User) error {
	if result := s.Db.First(&model.User{}, user.ID); !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return dict.ErrUserAlreadyExists
	}

	s.Db.Create(&user)

	return nil
}

func (s *userService) UpdateUser(user model.User) error {
	if result := s.Db.First(&model.User{}, user.ID); errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return dict.ErrUserNotExists
	}
	s.Db.Save(&user)

	return nil
}

func (s *userService) DeleteUser(user model.User) error {
	result := s.Db.Delete(user)

	return result.Error
}

func (s *userService) DeleteUserById(id int) error {
	result := s.Db.Delete(&model.User{}, id)

	return result.Error
}

func (s *userService) VerifyPassword(user model.User, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return dict.ErrWrongPassword
	}

	return nil
}
