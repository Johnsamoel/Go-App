package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID          int64  `json:"Id,omitempty" bson:"Id,omitempty"`
	Name        string `json:"name,omitempty" bson:"name,omitempty" validate:"required"`
	Email       string `json:"email,omitempty" bson:"email,omitempty" validate:"required"`
	Password    string `json:"password,omitempty" bson:"password,omitempty" validate:"required"`
	PhoneNumber string `json:"phoneNumber,omitempty" bson:"phoneNumber,omitempty" validate:"required"`
}

type AuthUser struct {
	Email       string `json:"email,omitempty" bson:"email,omitempty" validate:"required"`
	Password    string `json:"password,omitempty" bson:"password,omitempty" validate:"required"`
}

func (u *User) HashPassword() error {
	// Generate a bcrypt hash of the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

