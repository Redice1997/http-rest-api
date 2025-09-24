package sqlstorage_test

import (
	"os"
	"testing"
)

var (
	connectionString string
)

func TestMain(m *testing.M) {
	connectionString = os.Getenv("TEST_DB_CONNECTION")
	if connectionString == "" {
		connectionString = "host=localhost port=5432 user=api password=password dbname=test_api_db sslmode=disable"
	}

	os.Exit(m.Run())
}
