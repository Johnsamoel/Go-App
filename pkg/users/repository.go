package users

import (
	"database/sql"
	"fmt"
	"net/http"

	dataSource "example.com/fintech-app/db"
	"example.com/fintech-app/models"
)

type UserRepo interface {
	CreateUser(*models.User) (*models.User, error)
	Login(w http.ResponseWriter, r *http.Request, email, password string) (string, error)
	Logout(w http.ResponseWriter)
	GetUser(int64) (*models.User, error)
	DeleteUser(int64) error
	UpdateUser(int64, map[string]interface{}) error
}

type userRepo struct {
	db *sql.DB
	Methods dataSource.UserRepository
}

// NewUserRepo initializes and returns an instance of userRepo implementing UserRepo
func NewUserRepo(db *sql.DB) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(userData *models.User) (*models.User, error) {
	newUserData , err := r.Methods.CreateNewUser(userData)

	return newUserData, err
}

func (r *userRepo) Login(w http.ResponseWriter, req *http.Request, email, password string) (string, error) {
	token, err := r.Methods.Login(w, req, email, password)

	if err != nil {
		return "", fmt.Errorf("something went wrong : %v", err)
	}

	return token, err
}

func (r *userRepo) Logout(w http.ResponseWriter) {
	r.Methods.Logout(w)
}

func (r *userRepo) GetUser(userId int64) (*models.User, error) {
	user , err := r.Methods.GetUserById(userId)

	return user, err
}

func (r *userRepo) DeleteUser(userId int64) error {
	err := r.Methods.DeleteUser(userId)
	return err
}

func (r *userRepo) UpdateUser(userId int64, userData map[string]interface{}) error {
	err := r.Methods.UpdateUserData(userId, userData)
	return err
}
