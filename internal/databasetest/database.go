// Package databasetest provides PostgreSQL databases for tests.
package databasetest

import (
	"database/sql"
	"fmt"
	"net/url"
	"sync"
	"testing"

	"github.com/dogmatiq/example/ui/projections"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// NewDatabase returns an isolated PostgreSQL database with the projection
// schema applied. The connection pool is closed when the test ends.
func NewDatabase(t testing.TB) *sql.DB {
	t.Helper()

	container := getContainer(t)
	dsn, err := container.ConnectionString(t.Context(), "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	admin, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { admin.Close() })

	name := "example-" + uuid.NewString()
	if _, err := admin.ExecContext(t.Context(), `CREATE DATABASE "`+name+`"`); err != nil {
		t.Fatalf("unable to create test database: %s", err)
	}

	u, err := url.Parse(dsn)
	if err != nil {
		t.Fatal(err)
	}
	u.Path = "/" + name

	db, err := projections.NewDB(t.Context(), u.String())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	return db
}

var containerState struct {
	mutex     sync.Mutex
	container *postgres.PostgresContainer
}

// getContainer caches the container to avoid a Docker round-trip per test.
func getContainer(t testing.TB) *postgres.PostgresContainer {
	t.Helper()

	containerState.mutex.Lock()
	defer containerState.mutex.Unlock()

	if containerState.container != nil {
		return containerState.container
	}

	container, err := postgres.Run(
		t.Context(),
		"postgres:18-alpine",
		// Reuse is scoped to testcontainers' session so a different session's
		// Ryuk reaper cannot remove a container that is still in use.
		testcontainers.WithReuseByName(fmt.Sprintf("dogmatiq-example-%s", testcontainers.SessionID())),
		// The image starts a temporary server during initialization. Wait for
		// both the final server's log message and its listening port.
		postgres.BasicWaitStrategies(),
		postgres.WithUsername("dogmatiq-example-test"),
		postgres.WithPassword("password"),
	)
	if err != nil {
		t.Fatalf("unable to start PostgreSQL container: %s", err)
	}

	// Ryuk owns the container lifecycle; an individual test must not stop it.
	containerState.container = container
	return container
}
