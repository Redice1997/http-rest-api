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
	s, clear := sqlstorage.NewTestStorage(t, connectionString)
	defer clear("users")

	err := s.User().Create(context.Background(), model.TestUser(t))

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
