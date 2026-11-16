<div align="center">

# Dogma Example Application

An example [Dogma](https://github.com/dogmatiq/dogma) application, which
implements the features of a simple bank.

[![Build Status](https://img.shields.io/github/actions/workflow/status/dogmatiq/example/ci.yml?style=for-the-badge&branch=main)](https://github.com/dogmatiq/example/actions/workflows/ci.yml)
[![Code Coverage](https://img.shields.io/codecov/c/github/dogmatiq/example/main.svg?style=for-the-badge)](https://codecov.io/github/dogmatiq/example)

</div>

## Running locally

The UI projections use PostgreSQL. Start a local database:

```sh
docker run --rm --name dogma-example-postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=bank \
  -p 5432:5432 postgres:18-alpine
```

Then run the application with a connection string:

```sh
DATABASE_URL='postgres://postgres:password@localhost:5432/bank?sslmode=disable' \
  go run ./cmd/bank
```

The application creates the projection schema on startup and serves the UI at
<http://localhost:8080>. Use `sslmode=disable` only for a local database;
configure TLS as required by your database provider for remote connections.

Only the projections are stored in PostgreSQL at present. The application still
uses the in-memory testkit engine, so use a fresh database for each application
run. Reusing projections across engine restarts requires persistent engine
storage, which will be introduced separately.

## Testing

```sh
go test ./...
```

Projection and UI query tests use testcontainers with PostgreSQL 18 and require
a running Docker daemon. They start a session-scoped container automatically
and create an isolated database for each test; `DATABASE_URL` is not used.
The testcontainers reaper removes the container after the test session ends.
CGO is not required to build or test the application.
