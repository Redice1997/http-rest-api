package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	storage "github.com/Redice1997/http-rest-api/internal/app/storage/memorystorage"
	"github.com/stretchr/testify/assert"
)

func TestAPIServer_HandleHello(t *testing.T) {
	s := newServer(
		newLogger("debug"),
		storage.New(),
	)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)

	s.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.NotEmpty(t, rec.Body.String())
	assert.JSONEq(t, `{"message":"Hello, World!"}`, rec.Body.String())
}
