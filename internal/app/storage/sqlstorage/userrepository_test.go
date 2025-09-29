package sqlstorage_test

import (
	"context"
	"testing"

	"github.com/Redice1997/http-rest-api/internal/app/model"
	"github.com/Redice1997/http-rest-api/internal/app/storage/sqlstorage"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	// arrange
	s, clear := sqlstorage.NewTestStorage(t, connectionString)
	defer clear("users")

	// act
	err := s.User().Create(context.Background(), model.TestUser(t))

	// assert
	assert.NoError(t, err)
}

func TestUserRepository_GetByEmail(t *testing.T) {
	// arrange
	s, clear := sqlstorage.NewTestStorage(t, connectionString)
	defer clear("users")

	email := "user@example.org"

	// act
	_, err := s.User().GetByEmail(context.Background(), email)

	// assert
	assert.ErrorIs(t, err, model.ErrRecordNotFound)

	// arrange
	u := model.TestUser(t)
	u.Email = email
	s.User().Create(context.Background(), u)

	// act
	u, err = s.User().GetByEmail(context.Background(), email)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, email, u.Email)
}

func TestUserRepository_GetByID(t *testing.T) {
	// arrange
	s, clear := sqlstorage.NewTestStorage(t, connectionString)
	defer clear("users")
	u := model.TestUser(t)

	// act
	_, err := s.User().GetByID(context.Background(), u.ID)

	// assert
	assert.ErrorIs(t, err, model.ErrRecordNotFound)

	// arrange
	s.User().Create(context.Background(), u)

	// act
	u, err = s.User().GetByID(context.Background(), u.ID)

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, u)
}
