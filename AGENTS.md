# Repository conventions

These instructions apply to all work in this repository. Read them before editing
code. Deliver changes that already follow these conventions; formatting and a
review of the diff are part of the work, not follow-up tasks for the user.

## Readability

- Separate logical steps with blank lines: setup, loading data, validation,
  mutation, persistence, event publication, and return. Keep a call and its
  immediate error check together. Keep tightly related assignments together.
- Leave a blank line after a completed control-flow block before the next
  independent operation. Avoid blank lines immediately inside braces.
- Use one field per line in multi-field struct literals and one statement per
  line. Do not compress setup, error handling, or business operations to save space.
- Prefer names that describe their purpose. Short names are fine for conventional
  receivers and tiny loops; use descriptive names for values spanning several steps.
- Keep functions focused. Extract a helper when it represents a meaningful
  operation, not merely to hide a few lines or satisfy a line-count threshold.
- Comments explain intent, invariants, or tradeoffs. Avoid narrating obvious code.
  Document exported types and functions.

## Go imports and formatting

- Use a single import block with three groups, separated by blank lines:
  standard library, third-party packages, then `planningpoker/...` packages.
- Sort imports lexicographically within each group. Aliased and blank imports
  stay in the group matching their import path.
- Run the repository style check before finishing. Apply formatting only to
  first-party source; never hand-format `vendor/` or generated files.
- Formatting does not replace a readability review. In particular, group test
  setup and assertions by behavior rather than writing an uninterrupted sequence.
- Follow the existing Vue/JavaScript conventions and ESLint configuration when
  changing the frontend; do not introduce another formatting system.

## Domain-driven design and dependency direction

- Business concepts and rules belong in `internal/domain/<context>` (currently
  games, users, state, and events). Use the domain's established vocabulary.
- Entities and aggregates own their invariants and state changes. Domain services
  coordinate use cases through domain contracts; adapters must not duplicate rules.
- Domain packages must not import `internal/infra`, `cmd`, database clients,
  `database/sql`, or HTTP/WebSocket frameworks. Infrastructure depends on domain
  interfaces and types, not the other way around.
- Define ports/interfaces in the domain package that consumes them. Implement
  persistence, transport, authentication adapters, and event delivery in
  `internal/infra`. Assemble concrete dependencies in `cmd/poker`.
- Keep serialization DTOs, SQL, migrations, and storage configuration in
  infrastructure. Use domain constructors and methods for business operations;
  `NewRaw` constructors are for repository hydration only.
- Preserve aggregate boundaries. Make read-modify-write operations atomic using
  repository transactions/locks, and publish domain events only after persistence
  succeeds. Do not introduce cross-context table access to bypass a domain port.
- Extend an existing pattern when it fits. Avoid generic frameworks, extra layers,
  or abstractions without a concrete requirement.

## Errors, resources, and tests

- Handle errors explicitly; add useful context with `%w` when crossing a boundary.
  Keep expected missing-record behavior consistent with repository contracts.
- Do not log credentials or connection strings. Return safe public errors from
  transport adapters while preserving useful diagnostics internally.
- Bound external operations with contexts/timeouts, and close resources. Explain
  intentionally ignored errors when the reason is not evident (such as rollback
  after a committed transaction).
- Test meaningful behavior and failure paths. Persistence changes need real
  PostgreSQL integration tests; concurrent mutations need concurrency coverage.
- Use isolated databases, containers, volumes, and ports for verification. Do not
  reuse or modify another project's services. Do not change production as part of
  a local formatting or implementation task.
- Change dependencies intentionally. Regenerate `vendor/` with Go tooling when
  dependencies change; never patch third-party source by hand.
- Do not weaken checks, blanket-suppress findings, or alter unrelated behavior to
  make a check pass. Explain any check that cannot run and its concrete limitation.

## Required finish checks

Use Go 1.23+ and golangci-lint **v1.61.0** (the version pinned in CI):

```sh
golangci-lint run --config .golangci-style.yml ./...
go vet -mod=vendor ./...
go test -mod=vendor -race ./...
git diff --check
```

Use `golangci-lint run --config .golangci-style.yml --fix ./...` to fix import order
and Go formatting. Resolve remaining spacing findings manually, inspect the diff,
and rerun the check without `--fix`. This focused configuration enforces spacing,
import order, and domain import boundaries; the older `.golangci.yml` retains the
broader legacy lint configuration.

For persistence changes, also run the integration suite with `TEST_DATABASE_URL`
pointing to a disposable database as described in `docs/postgres.md`. For runtime
behavior or container configuration changes, build and smoke-test the Docker
image. Do not rerun unrelated or expensive checks solely for a prose-only edit.
