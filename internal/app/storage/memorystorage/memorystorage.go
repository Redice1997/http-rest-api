package memorystorage

import (
	"context"

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

func (s *Storage) BeginTx(_ context.Context) (storage.TxStorage, error) {
	return &TxStorage{
		userRepository: s.userRepository,
	}, nil
}

type TxStorage struct {
	userRepository *UserRepository
}

func (s *TxStorage) User() storage.UserRepository {
	return s.userRepository
}

func (s *TxStorage) Commit() error {
	return nil
}

func (s *TxStorage) Rollback() error {
	return nil
}
