package sqlstorage_test

import (
	"context"
	"testing"

	"github.com/Redice1997/http-rest-api/internal/app/model"
	"github.com/Redice1997/http-rest-api/internal/app/storage"
	"github.com/Redice1997/http-rest-api/internal/app/storage/sqlstorage"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, clear := sqlstorage.TestDB(t, connectionString)
	defer clear("users")
	s := sqlstorage.New(db)

	err := s.User().Create(context.Background(), &model.User{
		Email: "user@example.org",
	})

	assert.NoError(t, err)
}

func TestUserRepository_GetByEmail(t *testing.T) {

	db, clear := sqlstorage.TestDB(t, connectionString)
	defer clear("users")
	s := sqlstorage.New(db)

	email := "user@example.org"

	_, err := s.User().GetByEmail(context.Background(), email)

	assert.ErrorIs(t, err, storage.ErrRecordNotFound)

	s.User().Create(context.Background(), &model.User{
		Email: email,
	})

	u, err := s.User().GetByEmail(context.Background(), email)

	assert.NoError(t, err)
	assert.Equal(t, email, u.Email)
}
