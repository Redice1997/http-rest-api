package sqlstorage

import (
	"context"
	"database/sql"

	"github.com/Redice1997/http-rest-api/internal/app/storage"

	_ "github.com/lib/pq"
)

type Storage struct {
	db             *sql.DB
	userRepository *userRepository
}

func New(connectionString string) (*Storage, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	s := new(Storage)
	s.db = db
	s.userRepository = newUserRepository(db)

	return s, nil
}

func (s *Storage) User() storage.UserRepository {
	return s.userRepository
}

func (s *Storage) BeginTx(ctx context.Context) (storage.TxStorage, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &TxStorage{
		userRepository: newUserRepository(tx),
		db:             tx,
	}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

type TxStorage struct {
	db             *sql.Tx
	userRepository *userRepository
}

func (s *TxStorage) User() storage.UserRepository {
	return s.userRepository
}

func (s *TxStorage) Rollback() error {
	return s.db.Rollback()
}

func (s *TxStorage) Commit() error {
	return s.db.Commit()
}

type sqlDB interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}
