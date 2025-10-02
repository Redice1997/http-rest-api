package memorystorage

import (
	"context"
	"sync"

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

func (s *Storage) BeginTx(ctx context.Context) (storage.TxStorage, error) {
	mut := new(sync.Mutex)
	mut.Lock()

	return &TxStorage{
		userRepository: NewUserRepository(s),
		mut:            mut,
	}, nil
}

type TxStorage struct {
	mut            *sync.Mutex
	userRepository *UserRepository
}

func (s *TxStorage) User() storage.UserRepository {
	return s.userRepository
}

func (s *TxStorage) Commit() error {
	s.mut.Unlock()
	return nil
}

func (s *TxStorage) Rollback() error {
	s.mut.Unlock()
	return nil
}
