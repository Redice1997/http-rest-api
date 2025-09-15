package sqlstorage

import (
	"database/sql"

	"github.com/Redice1997/http-rest-api/internal/app/storage"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) User() storage.UserRepository {
	panic("implement me")
}
