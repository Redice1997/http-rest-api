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
		connectionString = "host=localhost port=5432 user=postgres password=postgres dbname=api_test sslmode=disable"
	}

	os.Exit(m.Run())
}
