package users

import (
	models "example.com/fintech-app/models"
)

type UserRepo interface {
	CreateUser(*models.User) (*models.User, error)
	GetUser(int) (*models.User, error)
	DeleteUser(int) error
	UpdateUser(int, *models.User) (*models.User, error)
}


