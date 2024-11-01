package services

import (
	"errors"

	"github.com/henok-tesfu/expense-manager/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (us *UserService) RegisterUser(username, email, password string) (*models.User, error) {
	if len(password) < 6 {
		return nil, errors.New("password must be at least 6 characters long")
	}

	if user, _ := models.GetUserByEmail(email); user != nil {
		return nil, errors.New("email already registered")
	}

	userID, err := models.RegisterUser(username, email, password)
	if err != nil {
		return nil, err
	}

	return models.GetUserByID(int(userID))
}

func (us *UserService) AuthenticateUser(email, password string) (*models.User, error) {
	user, err := models.GetUserByEmail(email)
	if err != nil || user == nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}
