package memorystorage

import (
	"context"
	"sync"

	"github.com/Redice1997/http-rest-api/internal/app/storage"
)

var mut sync.Mutex

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

func (s *Storage) BeginTx(_ context.Context) (storage.TxStorage, error) {
	mut.Lock()

	return &TxStorage{
		userRepository: NewUserRepository(s),
	}, nil
}

type TxStorage struct {
	userRepository *UserRepository
}

func (s *TxStorage) User() storage.UserRepository {
	return s.userRepository
}

func (s *TxStorage) Commit() error {
	mut.Unlock()
	return nil
}

func (s *TxStorage) Rollback() error {
	mut.Unlock()
	return nil
}
