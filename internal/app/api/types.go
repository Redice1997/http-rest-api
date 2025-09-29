package api

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
