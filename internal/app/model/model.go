package model

type User struct {
	ID                int64
	Email             string
	Password          string
	EncryptedPassword string
}

type ContextKey int8

const (
	SessionName     string     = "http-rest-api-session"
	SessionUserID   string     = "user_id"
	HeaderRequestID string     = "X-Request-ID"
	CtxUserKey      ContextKey = iota
	CtxRequestIDKey
)
