package users

import (
	"fmt"
	"net/http"

	"example.com/fintech-app/db"
	"example.com/fintech-app/models"
)

type UserService interface {
	LogInService(string, string, http.ResponseWriter, *http.Request)
	LogoutService(w http.ResponseWriter)
	CreateUserService(*models.User) (*models.User, error)
	DeleteUserService(int) (*models.User, error)
	UpdateUserService(int) error
	GetUserService(int, *models.User) (*models.User, error)
}

func LogInService(email, password string, w http.ResponseWriter, r *http.Request) (string, error) {
	otp, err := db.Login(w, r, email, password)

	if err != nil {
		return "", fmt.Errorf("something went wrong: %v", err)
	}

	if otp != "" {
		return otp, nil
	}

	return "", nil
}

func LogoutService(w http.ResponseWriter) {
	db.Logout(w)
}

func CreateUserService(userData *models.User) (*models.User, error) {
	user, err := db.CreateNewUser(userData)

	if err != nil {
		return nil, fmt.Errorf("error creating new user: %v", err)
	}

	return user, nil
}

func DeleteUserService(userId int64) error {
	err := db.DeleteUser(userId)

	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}

	return nil
}

func UpdateUserService(userId int64, userData map[string]interface{}) error {
	err := db.UpdateUserData(userId, userData)

	if err != nil {
		return fmt.Errorf("error updating user: %v", err)
	}

	return nil
}

func GetUserService(userId int64) (*models.User, error) {
	user, err := db.GetUserById(userId)

	if err != nil {
		return nil, fmt.Errorf("error finding user: %v", err)
	}

	return user, nil
}
