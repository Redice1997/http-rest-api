package api

import "errors"

var (
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
	ErrAlreadyExists          = errors.New("already exists")
	ErrUnauthorized           = errors.New("unauthorized")
)

type ContextKey int8

const (
	SessionName string     = "http-rest-api-session"
	CtxUserKey  ContextKey = iota
	CtxRequestIDKey
)

type UserCreateRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"securepassword"`
}

type UserResponse struct {
	ID    int64  `json:"id" example:"1"`
	Email string `json:"email" example:"user@example.com"`
}

type SessionCreateRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"securepassword"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"invalid request"`
}
