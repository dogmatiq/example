<div align="center">

# Dogma Example Application

An example [Dogma](https://github.com/dogmatiq/dogma) application, which
implements the features of a simple bank.

[![Build Status](https://img.shields.io/github/actions/workflow/status/dogmatiq/example/ci.yml?style=for-the-badge&branch=main)](https://github.com/dogmatiq/example/actions/workflows/ci.yml)
[![Code Coverage](https://img.shields.io/codecov/c/github/dogmatiq/example/main.svg?style=for-the-badge)](https://codecov.io/github/dogmatiq/example)

</div>

## Running locally

Start PostgreSQL:

```sh
docker run --rm -d --name dogma-example-postgres \
  -e POSTGRES_PASSWORD=password -e POSTGRES_DB=bank \
  -p 5432:5432 postgres:18-alpine
```

Once PostgreSQL is ready, run the example:

```sh
DOGMA_EXAMPLE_ENGINE_DSN='postgres://postgres:password@localhost:5432/bank?sslmode=disable' \
DOGMA_EXAMPLE_PROJECTION_DSN='postgres://postgres:password@localhost:5432/bank?sslmode=disable' \
  go run ./cmd/bank
```

The application creates the engine and projection schemas on startup and
preserves state across restarts. Both DSNs can use the same PostgreSQL 18
database, as above, or separate databases. Use `sslmode=disable` only for local
databases.

Run `go test ./...` with Docker available. Projection tests use testcontainers
with PostgreSQL 18 and an isolated database per test. CGO is not required.
