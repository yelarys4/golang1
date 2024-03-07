package services

import (
	"errors"
	"github.com/Hoaper/golang_university/app/models"
	"github.com/Hoaper/golang_university/app/repositories"
	"github.com/sirupsen/logrus"
)

type UserService interface {
	CreateUser(user *models.User) error
	AuthenticateUser(login, password string) (*models.User, error)
	GetUserByLogin(login string) (*models.User, error)
	DeleteUser(login string) error
	GetUserByToken(token string) (*models.User, error)
	UpdateUser(user *models.User) error
}

type AuthService struct {
	UserRepository repositories.UserRepository
}

func NewAuthService(userRepository *repositories.UserRepository) *AuthService {
	return &AuthService{UserRepository: *userRepository}
}

func (s *AuthService) GetUserByToken(token string) (*models.User, error) {
	user, err := s.UserRepository.GetUserByToken(token)
	if err != nil {
		logrus.WithError(err).WithField("token", token).Error("Error getting user by login")
		return nil, err
	}

	logrus.WithField("token", token).Info("User retrieved successfully by login")
	return user, nil

}

func (s *AuthService) UpdateUser(user *models.User) error {
	err := s.UserRepository.UpdateUser(user)
	if err != nil {
		logrus.WithError(err).WithField("user_id", user.ID).Error("Error updating user by id")
		return err
	}
	return nil
}

func (s *AuthService) GetUserByLogin(login string) (*models.User, error) {
	user, err := s.UserRepository.GetUserByLogin(login)
	if err != nil {
		logrus.WithError(err).WithField("login", login).Error("Error getting user by login")
		return nil, err
	}

	logrus.WithField("login", login).Info("User retrieved successfully by login")
	return user, nil
}

func (s *AuthService) CreateUser(user *models.User) error {
	findUser, _ := s.UserRepository.GetUserByLogin(user.Login)

	if findUser != nil {
		err := errors.New("user already exists")
		logrus.WithError(err).WithField("login", user.Login).Error("Error creating user - user already exists")
		return err
	}

	err := s.UserRepository.CreateUser(user)
	if err != nil {
		logrus.WithError(err).WithField("login", user.Login).Error("Error creating user")
		return err
	}

	logrus.WithField("login", user.Login).Info("User created successfully")
	return nil
}

func (s *AuthService) AuthenticateUser(login, password string) (*models.User, error) {
	user, err := s.UserRepository.GetUserByLogin(login)
	if err != nil {
		logrus.WithError(err).WithField("login", login).Error("Error getting user for authentication")
		return nil, err
	}

	if user.Password != password {
		err := errors.New("invalid credentials")
		logrus.WithField("login", login).Error("Invalid credentials")
		return nil, err
	}

	logrus.WithField("login", login).Info("User authenticated successfully")
	return user, nil
}

func (s *AuthService) DeleteUser(login string) error {
	err := s.UserRepository.DeleteUser(login)
	if err != nil {
		logrus.WithError(err).WithField("login", login).Error("Error deleting user")
		return err
	}

	logrus.WithField("login", login).Info("User deleted successfully")
	return nil
}
