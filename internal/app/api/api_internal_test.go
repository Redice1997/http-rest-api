package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Redice1997/http-rest-api/internal/app/model"
	"github.com/Redice1997/http-rest-api/internal/app/storage/memorystorage"
	"github.com/Redice1997/http-rest-api/internal/app/storage/sqlstorage"
	"github.com/stretchr/testify/assert"
)

func TestAPI_HandleUserCreate(t *testing.T) {
	cfg := NewConfig()
	db, cleanup := sqlstorage.NewTestStorage(t, cfg.DbConnectionString)
	defer cleanup("users")
	s := New(
		cfg,
		db,
	)

	testCases := []struct {
		name     string
		request  any
		expected int
	}{
		{
			name: "valid",
			request: map[string]any{
				"email":    "john.doe@example.com",
				"password": "password123",
			},
			expected: http.StatusCreated,
		},
		{
			name: "invalid",
			request: map[string]any{
				"email":    "john.doe",
				"password": "password123",
			},
			expected: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := bytes.NewBuffer([]byte{})
			json.NewEncoder(b).Encode(tc.request)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/users", b)

			s.srv.Handler.ServeHTTP(rec, req)

			assert.Equal(t, tc.expected, rec.Code)
		})
	}
}

func TestAPI_HandleSessionCreate(t *testing.T) {
	db := memorystorage.New()
	s := New(NewConfig(), db)
	u := model.TestUser(t)

	db.User().Create(context.Background(), u)

	testCases := []struct {
		name     string
		JSON     any
		expected int
	}{
		{
			name: "valid",
			JSON: map[string]string{
				"email":    u.Email,
				"password": u.Password,
			},
			expected: http.StatusOK,
		},
		{
			name:     "invalid payload",
			JSON:     "invalid",
			expected: http.StatusBadRequest,
		},
		{
			name: "invalid email",
			JSON: map[string]string{
				"email":    "john@example.org",
				"password": u.Password,
			},
			expected: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := bytes.NewBuffer([]byte{})
			json.NewEncoder(b).Encode(tc.JSON)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sessions", b)

			s.srv.Handler.ServeHTTP(rec, req)

			assert.Equal(t, tc.expected, rec.Code)
		})
	}
}
