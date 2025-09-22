package sqlstorage

import (
	"fmt"
	"strings"
	"testing"
)

func NewTestStorage(t *testing.T, connectionStr string) (*Storage, func(...string)) {
	t.Helper()

	storage, err := New(connectionStr)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}

	return storage, func(tables ...string) {
		if len(tables) > 0 {
			if _, err := storage.db.Exec(
				fmt.Sprintf(
					"TRUNCATE %s RESTART IDENTITY CASCADE",
					strings.Join(tables, ", "),
				),
			); err != nil {
				t.Fatalf("failed to truncate tables: %v", err)
			}
		}

		storage.Close()
	}
}
