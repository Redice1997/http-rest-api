package user

import "errors"

var (
	ErrUserValidation         = errors.New("user validation error")
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
	ErrUserExists             = errors.New("user already exists")
	ErrUserUnauthorized       = errors.New("unauthorized")
)
