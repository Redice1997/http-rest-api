package sqlstorage_test

import (
	"context"
	"log"
	"testing"

	"github.com/Redice1997/http-rest-api/internal/app/model"
	"github.com/Redice1997/http-rest-api/internal/app/storage"
	"github.com/Redice1997/http-rest-api/internal/app/storage/sqlstorage"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	s, clear := sqlstorage.NewTestStorage(t, connectionString)
	defer clear("users")

	err := s.User().Create(context.Background(), model.TestUser(t))

	log.Println(err)
	assert.NoError(t, err)
}

func TestUserRepository_GetByEmail(t *testing.T) {
	s, clear := sqlstorage.NewTestStorage(t, connectionString)
	defer clear("users")

	email := "user@example.org"

	_, err := s.User().GetByEmail(context.Background(), email)

	assert.ErrorIs(t, err, storage.ErrRecordNotFound)

	u := model.TestUser(t)
	u.Email = email
	s.User().Create(context.Background(), u)

	u, err = s.User().GetByEmail(context.Background(), email)

	assert.NoError(t, err)
	assert.Equal(t, email, u.Email)
}

func TestUserRepository_GetByID(t *testing.T) {
	s, clear := sqlstorage.NewTestStorage(t, connectionString)
	defer clear("users")
	u := model.TestUser(t)

	_, err := s.User().GetByID(context.Background(), u.ID)
	assert.ErrorIs(t, err, storage.ErrRecordNotFound)

	s.User().Create(context.Background(), u)
	u, err = s.User().GetByID(context.Background(), u.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}
