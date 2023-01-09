package repository

import (
	"pendekin/dto"
	"pendekin/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByUsernameEmail(username string, email string) (model.User, error)
	GetUserByUsernameEmailPassword(username string, email string, password string) (model.User, error)
	InsertUser(user *dto.UserRegister) error
}

type userRepository struct {
	db *gorm.DB
}

func InitUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (ur *userRepository) GetUserByUsernameEmail(username string, email string) (model.User, error) {
	var user model.User
	if err := ur.db.Table("users").Where("username = ? or email = ?", username, email).First(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (ur *userRepository) GetUserByUsernameEmailPassword(username string, email string, password string) (model.User, error) {
	var user model.User
	if err := ur.db.Table("users").Where("username = ? or email = ?", username, email).First(&user).Error; err != nil {
		return model.User{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (ur *userRepository) InsertUser(user *dto.UserRegister) error {
	if err := ur.db.Table("users").Create(&user).Error; err != nil {
		return err
	}
	return nil
}
