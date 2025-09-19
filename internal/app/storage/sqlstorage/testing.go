package sqlstorage

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
)

func TestDB(t *testing.T, connectionStr string) (*sql.DB, func(...string)) {
	t.Helper()

	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		t.Fatalf("failed to open db connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("failed to ping db: %v", err)
	}

	return db, func(tables ...string) {
		if len(tables) > 0 {
			if _, err := db.Exec(
				fmt.Sprintf(
					"TRUNCATE %s RESTART IDENTITY CASCADE",
					strings.Join(tables, ", "),
				),
			); err != nil {
				t.Fatalf("failed to truncate tables: %v", err)
			}
		}

		db.Close()
	}
}
