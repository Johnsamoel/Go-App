package users

import (
	"fmt"
	"net/http"
	"example.com/fintech-app/models"
)

type UserService interface {
	LogInService(string, string, http.ResponseWriter, *http.Request) (string, error)
	LogoutService(w http.ResponseWriter)
	CreateUserService(*models.User) (*models.User, error)
	DeleteUserService(int64) (error)
	UpdateUserService(int64 , map[string]interface{}) error
	GetUserService(int64) (*models.User, error)
}

type userService struct {
	repo UserRepo
}

func NewUserService(repo UserRepo) UserService {
	return &userService{repo: repo}
}

func (s *userService) LogInService(email, password string, w http.ResponseWriter, r *http.Request) (string, error) {
	otp, err :=  s.repo.Login(w, r, email, password)

	if err != nil {
		return "", fmt.Errorf("something went wrong: %v", err)
	}

	if otp != "" {
		return otp, nil
	}

	return "", nil
}

func (s *userService) LogoutService(w http.ResponseWriter) {
	s.repo.Logout(w)
}

func (s *userService) CreateUserService(userData *models.User) (*models.User, error) {
	user, err := s.repo.CreateUser(userData)

	if err != nil {
		return nil, fmt.Errorf("error creating new user: %v", err)
	}

	return user, nil
}

func (s *userService) DeleteUserService(userId int64) error {
	err := s.repo.DeleteUser(userId)

	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}

	return nil
}

func (s *userService) UpdateUserService(userId int64, userData map[string]interface{}) error {
	err := s.repo.UpdateUser(userId, userData)

	if err != nil {
		return fmt.Errorf("error updating user: %v", err)
	}

	return nil
}

func (s *userService) GetUserService(userId int64) (*models.User, error) {
	user, err := s.repo.GetUser(userId)

	if err != nil {
		return nil, fmt.Errorf("error finding user: %v", err)
	}

	return user, nil
}
