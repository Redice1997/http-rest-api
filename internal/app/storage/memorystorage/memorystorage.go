package memorystorage

import "github.com/Redice1997/http-rest-api/internal/app/storage"

type Storage struct {
	userRepository *UserRepository
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) User() storage.UserRepository {
	if s.userRepository == nil {
		s.userRepository = NewUserRepository(s)
	}

	return s.userRepository
}
