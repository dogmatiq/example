package projections_test

import (
	"testing"

	"github.com/dogmatiq/example/internal/xtesting"
	"github.com/dogmatiq/example/ui/projections"
)

func Test_NewDB_requires_connection_string(t *testing.T) {
	db, err := projections.NewDB(t.Context(), "")
	if err == nil || db != nil {
		t.Fatalf("got %v, %v; want no database and an error", db, err)
	}
}

func Test_CreateSchema_preserves_existing_data(t *testing.T) {
	db := xtesting.NewDatabase(t)
	if _, err := db.ExecContext(t.Context(), `INSERT INTO customers (id, name) VALUES ('C001', 'Anna')`); err != nil {
		t.Fatal(err)
	}
	if err := projections.CreateSchema(t.Context(), db); err != nil {
		t.Fatal(err)
	}
	var name string
	if err := db.QueryRowContext(t.Context(), `SELECT name FROM customers WHERE id = 'C001'`).Scan(&name); err != nil {
		t.Fatal(err)
	}
	if name != "Anna" {
		t.Fatalf("got customer name %q, want Anna", name)
	}
}

func Test_databases_are_isolated(t *testing.T) {
	first := xtesting.NewDatabase(t)
	second := xtesting.NewDatabase(t)
	if _, err := first.ExecContext(t.Context(), `INSERT INTO customers (id, name) VALUES ('C001', 'Anna')`); err != nil {
		t.Fatal(err)
	}
	var count int
	if err := second.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM customers`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("second database contains %d customers, want none", count)
	}
}
