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

## Commands and the UI

Runkit accepts commands asynchronously. The deposit, withdrawal, and transfer
handlers use `WithEventObserver()` to wait for outcome events, so the UI can
show a result such as a declined withdrawal. These waits are bounded by the
request deadline; a timeout does not cancel a command already accepted by the
engine.

An observed event does not mean the UI projections have caught up. Account
creation uses `waitForAccount()` before redirecting so the next page can read
the customer and account. This demonstrates one way to support a synchronous
UI; a UI that shows pending work and updates as projections arrive would avoid
holding the request open.

## Container image

Run `make docker` to build `ghcr.io/dogmatiq/example:dev` locally. The existing
CI workflow publishes version-tagged images to GHCR when a `v*` Git tag is
pushed, for example `ghcr.io/dogmatiq/example:v0.1.0`.

The image listens on port 8080 and requires `DOGMA_EXAMPLE_ENGINE_DSN` and
`DOGMA_EXAMPLE_PROJECTION_DSN`. Both can point to the same PostgreSQL 18 database.

Railway can deploy the published image and follow new versions using
[image auto updates](https://docs.railway.com/deployments/image-auto-updates).
Make the GHCR package public after its first publication, configure the service's
DSNs, target port 8080 and health check path `/`, and disable sleeping so scheduled
operations continue to run.
