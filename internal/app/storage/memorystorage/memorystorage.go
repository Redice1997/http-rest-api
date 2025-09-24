package memorystorage

import (
	"github.com/Redice1997/http-rest-api/internal/app/storage"
)

type Storage struct {
	userRepository *UserRepository
}

func New() *Storage {
	s := new(Storage)
	s.userRepository = NewUserRepository(s)

	return s
}

func (s *Storage) User() storage.UserRepository {
	return s.userRepository
}
