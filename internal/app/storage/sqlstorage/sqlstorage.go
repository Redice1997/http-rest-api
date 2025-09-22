package sqlstorage

import (
	"database/sql"

	"github.com/Redice1997/http-rest-api/internal/app/storage"

	_ "github.com/lib/pq"
)

type Storage struct {
	db             *sql.DB
	userRepository *UserRepository
}

func New(connectionStirng string) (*Storage, error) {
	db, err := sql.Open("postgres", connectionStirng)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	s := new(Storage)
	s.db = db
	s.userRepository = NewUserRepository(s)

	return s, nil
}

func (s *Storage) User() storage.UserRepository {
	return s.userRepository
}

func (s *Storage) Close() error {
	return s.db.Close()
}
