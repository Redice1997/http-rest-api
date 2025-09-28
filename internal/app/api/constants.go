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
)
