package sqlstorage_test

import (
	"context"
	"testing"

	"github.com/Redice1997/http-rest-api/internal/app/model"
	"github.com/Redice1997/http-rest-api/internal/app/storage/sqlstorage"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, clear := sqlstorage.TestDB(t, connectionString)
	defer clear("users")
	storage := sqlstorage.New(db)

	err := storage.User().Create(context.Background(), &model.User{
		Email: "user@example.org",
	})

	assert.NoError(t, err)
}
