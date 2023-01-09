package service

import (
	"pendekin/dto"
	"pendekin/model"
	"pendekin/repository"
)

type UserService interface {
	CheckUserByUsernameEmail(username string, email string) (model.User, error)
	CheckUserByUsernameEmailPassword(username string, email string, password string) (model.User, error)
	StoreUser(user *dto.UserRegister) error
}

type userService struct {
	userRepository repository.UserRepository
}

func InitUserService(userRepository repository.UserRepository) *userService {
	return &userService{userRepository}
}

func (us *userService) CheckUserByUsernameEmail(username string, email string) (model.User, error) {
	user, err := us.userRepository.GetUserByUsernameEmail(username, email)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (us *userService) CheckUserByUsernameEmailPassword(username string, email string, password string) (model.User, error) {
	user, err := us.userRepository.GetUserByUsernameEmailPassword(username, email, password)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (us *userService) StoreUser(user *dto.UserRegister) error {
	if err := us.userRepository.InsertUser(user); err != nil {
		return err
	}
	return nil
}

