package projections_test

import (
	"database/sql"
	"path/filepath"
	"testing"

	. "github.com/dogmatiq/example/database"
	"github.com/dogmatiq/projectionkit/sqlprojection"
	_ "github.com/mattn/go-sqlite3"
)

// openDB returns a new database connection for use in tests.
//
// The database is created in a temporary directory, and is automatically closed
// and deleted when the test finishes.
func openDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", filepath.Join(t.TempDir(), "test.sqlite3"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		db.Close()
	})

	if err := sqlprojection.SQLiteDriver.CreateSchema(t.Context(), db); err != nil {
		t.Fatal(err)
	}

	if err := CreateSchema(t.Context(), db); err != nil {
		t.Fatal(err)
	}

	return db
}
