# Testing

Use Go 1.23+, Node 20+, and the repository's vendored Go dependencies.
Install the existing frontend dependencies with
`npm --prefix web ci --legacy-peer-deps`.

## Fast checks

```sh
go test -mod=vendor -race ./...
npm --prefix web test
```

Domain tests cover the game rules and service failure paths. HTTP adapter tests
check unavailable storage, safe public errors, and missing identities.

## PostgreSQL and service component tests

Set `TEST_DATABASE_URL` to a **disposable** PostgreSQL database as described in
[postgres.md](postgres.md#local-development-and-verification). These tests apply
migrations and persist test users and games; never point them at production or a
shared development database.

```sh
go test -mod=vendor -race -tags=integration ./...
go test -mod=vendor -race -count=1 -v -tags=component ./test/component
```

The component harness starts the real HTTP and Socket.IO adapters, domain
services, event bus, and PostgreSQL repositories on a random local port. Its
client suite uses only public HTTP requests and the frontend's Socket.IO client;
it neither reads database tables nor mocks service internals. The same suite can
test a disposable running Docker service:

```sh
POKER_TEST_URL=http://127.0.0.1:YOUR_TEST_PORT \
  node --test web/tests/component/service.test.js
```

The six scenarios cover identity and authentication, a complete game round,
vote privacy and reveal permissions, reconnect and departure behavior, invalid
requests, and simultaneous votes. Requests and event waits have deadlines;
state listeners attach before mutations so fast broadcasts are not missed.
Each scenario creates its own identities and games.

CI runs the suite both through the Go harness and against the built Docker image.
The harness collects Go coverage across package boundaries; the Docker run also
checks production wiring. Existing PostgreSQL integration tests cover migration
concurrency, reconnect persistence, rollback, and publication after commit.

## Coverage

```sh
go test -mod=vendor -race -tags=integration \
  -coverpkg=./cmd/...,./internal/... -covermode=atomic \
  -coverprofile=coverage-unit.txt ./...
go test -mod=vendor -race -tags=component \
  -coverpkg=./cmd/...,./internal/... -covermode=atomic \
  -coverprofile=coverage-component.txt ./test/component
```

Both reports are uploaded to Codecov, which combines their covered lines.
Coverage targets application packages, including the entrypoint, while excluding
development tooling and test helpers. It does not include JavaScript coverage.
The separate frontend model tests still run in CI.

Before finishing Go changes, also run the pinned golangci-lint v1.61.0 style
check, `go vet -mod=vendor ./...`, and `git diff --check` as required by AGENTS.md.
