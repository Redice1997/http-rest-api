package model

import (
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                int64  `json:"id"`
	Email             string `json:"email" validate:"required,email"`
	Password          string `json:"password,omitempty" validate:"required,min=8,max=100"`
	EncryptedPassword string `json:"-"`
}

func (u *User) Validate() error {
	return validate.Struct(u)
}

func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		encryptedPassword, err := encryptPassword(u.Password)
		if err != nil {
			return err
		}
		u.EncryptedPassword = encryptedPassword
	}
	return nil
}

func encryptPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

var validate = validator.New(validator.WithRequiredStructEnabled())
