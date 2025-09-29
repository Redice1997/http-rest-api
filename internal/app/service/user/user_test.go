package user_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Redice1997/http-rest-api/internal/app/model"
	"github.com/Redice1997/http-rest-api/internal/app/service/user"
	"github.com/stretchr/testify/assert"
)

func TestMe(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		u := model.TestUser(t)
		ctx := context.WithValue(context.Background(), model.CtxUserKey, u)

		// act
		got, err := user.Me(ctx)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, u.ID, got.ID)
		assert.Equal(t, u.Email, got.Email)
	})

	t.Run("error", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// act
		_, err := user.Me(ctx)

		// assert
		assert.ErrorIs(t, err, user.ErrUserUnauthorized)
	})
}

func TestSave(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		s := user.TestNew(t)
		u := model.TestUser(t)
		ctx := context.Background()

		// act
		err := s.Save(ctx, u)

		// assert
		assert.NoError(t, err)
		assert.NotEmpty(t, u.ID)
		assert.NotEmpty(t, u.EncryptedPassword)
	})

	t.Run("already exists", func(t *testing.T) {
		// arrange
		s := user.TestNew(t)
		u := model.TestUser(t)
		ctx := context.Background()
		s.Save(ctx, u)

		// act
		err := s.Save(ctx, u)

		// assert
		assert.ErrorIs(t, err, user.ErrUserExists)
	})

	t.Run("invalid user", func(t *testing.T) {
		// arrange
		s := user.TestNew(t)
		ctx := context.Background()
		u := model.TestUser(t)
		u.Email = "invalid"

		// act
		err := s.Save(ctx, u)

		// assert
		assert.ErrorIs(t, err, user.ErrUserValidation)
	})
}

func TestCreateHttpSession(t *testing.T) {
	t.Run("user exists", func(t *testing.T) {
		// arrange
		s := user.TestNew(t)
		u := model.TestUser(t)
		ctx := context.Background()
		s.Save(ctx, u)
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		// act
		err := s.CreateHttpSession(w, r, u)

		// assert
		assert.NoError(t, err)
		assert.NotEmpty(t, w.Header().Get("Set-Cookie"))
	})

	t.Run("user not exists", func(t *testing.T) {
		// arrange
		s := user.TestNew(t)
		u := model.TestUser(t)
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		// act
		err := s.CreateHttpSession(w, r, u)

		// assert
		assert.ErrorIs(t, err, user.ErrInvalidEmailOrPassword)
		assert.Empty(t, w.Header().Get("Set-Cookie"))
	})

	t.Run("user invalid password", func(t *testing.T) {
		// arrange
		s := user.TestNew(t)
		u := model.TestUser(t)
		ctx := context.Background()
		s.Save(ctx, u)

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		u.Password = "invalid"

		// act
		err := s.CreateHttpSession(w, r, u)

		// assert
		assert.ErrorIs(t, err, user.ErrInvalidEmailOrPassword)
		assert.Empty(t, w.Header().Get("Set-Cookie"))
	})
}
