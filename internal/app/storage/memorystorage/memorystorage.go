package memorystorage

import "github.com/Redice1997/http-rest-api/internal/app/storage"

type MemoryStorage struct {
	// Add fields if necessary
}

func New() *MemoryStorage {
	return &MemoryStorage{}
}

func (s *MemoryStorage) User() storage.UserRepository {
	panic("implement me")
}
