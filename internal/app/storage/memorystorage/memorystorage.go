package memorystorage

import "github.com/Redice1997/http-rest-api/internal/app/storage"

type Storage struct {
	userRepository *UserRepository
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) User() storage.UserRepository {
	panic("implement me")
}
