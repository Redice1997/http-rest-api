package app

type ContextKey int8

const (
	SessionName  string     = "http-rest-api-session"
	SessonUserID string     = "user_id"
	CtxUserKey   ContextKey = iota
	CtxRequestIDKey
)
