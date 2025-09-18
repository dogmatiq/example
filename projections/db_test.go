package projections_test

import (
	"context"
	"database/sql"

	. "github.com/dogmatiq/example/database"
	"github.com/dogmatiq/projectionkit/sqlprojection"
	"github.com/dogmatiq/sqltest"
)

func openDB(ctx context.Context) (*sqltest.Database, *sql.DB) {
	database, err := sqltest.NewDatabase(ctx, sqltest.SQLite3Driver, sqltest.SQLite)
	if err != nil {
		panic(err)
	}

	db, err := database.Open()
	if err != nil {
		panic(err)
	}

	if err := sqlprojection.SQLiteDriver.CreateSchema(ctx, db); err != nil {
		database.Close()
		panic(err)
	}

	if err := CreateSchema(ctx, db); err != nil {
		database.Close()
		panic(err)
	}

	return database, db
}
