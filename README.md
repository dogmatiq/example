<div align="center">

# Dogma Example Application

An example [Dogma](https://github.com/dogmatiq/dogma) application, which
implements the features of a simple bank.

[![Build Status](https://img.shields.io/github/actions/workflow/status/dogmatiq/example/ci.yml?style=for-the-badge&branch=main)](https://github.com/dogmatiq/example/actions/workflows/ci.yml)
[![Code Coverage](https://img.shields.io/codecov/c/github/dogmatiq/example/main.svg?style=for-the-badge)](https://codecov.io/github/dogmatiq/example)

</div>

## Running locally

Set `DOGMA_EXAMPLE_PROJECTION_DSN` to a PostgreSQL connection string:

```sh
DOGMA_EXAMPLE_PROJECTION_DSN='postgres://postgres:password@localhost:5432/bank?sslmode=disable' go run ./cmd/bank
```

The application creates its projection schema on startup. Use `sslmode=disable`
only for local databases. The engine still uses in-memory testkit storage, so
use a fresh projection database for each run until the runkit migration.

Run `go test ./...` with Docker available. Projection tests use testcontainers
with PostgreSQL 18 and an isolated database per test. CGO is not required.
