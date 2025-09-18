package memorystorage

import "github.com/Redice1997/http-rest-api/internal/app/model"

type UserRepository struct {
	storage *Storage
	users   map[string]*model.User
}
