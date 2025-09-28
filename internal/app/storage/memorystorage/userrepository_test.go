package memorystorage

import (
	"context"
	"testing"

	"github.com/Redice1997/http-rest-api/internal/app/model"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	// arrange
	s := New()
	user := model.TestUser(t)

	// act
	err := s.User().Create(context.Background(), user)

	// assert
	assert.NoError(t, err)
	assert.NotEmpty(t, user.ID)
}

func TestUserRepository_GetByID(t *testing.T) {
	// arrange
	s := New()
	user := model.TestUser(t)
	s.User().Create(context.Background(), user)

	// act
	user, err := s.User().GetByID(context.Background(), user.ID)

	// assert
	assert.NoError(t, err)
	assert.NotEmpty(t, user)
}

func TestUserRepository_GetByEmail(t *testing.T) {
	// arrange
	s := New()
	user := model.TestUser(t)
	s.User().Create(context.Background(), user)

	// act
	user, err := s.User().GetByEmail(context.Background(), user.Email)

	// assert
	assert.NoError(t, err)
	assert.NotEmpty(t, user)
}
