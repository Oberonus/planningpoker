# Repository conventions

These instructions apply to all work in this repository. Read them before editing
code. Deliver changes that already follow these conventions; formatting and a
review of the diff are part of the work, not follow-up tasks for the user.

## Commits and branches

Every new commit must follow
[Conventional Commits 1.0.0](https://www.conventionalcommits.org/en/v1.0.0/).
Apply these rules whenever creating a branch, commit, PR title, or squash commit.

Use this shared set of lowercase types:

| Type | Purpose |
| --- | --- |
| `feat` | New user-facing functionality |
| `fix` | Correcting a bug |
| `docs` | Documentation changes |
| `refactor` | Restructuring code without changing behavior |
| `test` | Adding or improving tests |
| `perf` | Improving performance |
| `style` | Formatting changes without changing behavior |
| `build` | Dependencies, build tooling, or image construction |
| `ci` | Continuous integration workflows |
| `chore` | Maintenance that does not fit another type |
| `revert` | Reverting an earlier change |

### Commit messages

- Format the subject as `<type>[optional scope][!]: <description>`, for example
  `feat(games): add confidence voting` or `fix(http): handle missing users`.
- Use a short lowercase scope naming the affected area when it adds clarity,
  such as `games`, `users`, `http`, `repository`, or `web`. Omit it for changes
  spanning several areas. Use the same scope for the same area across commits.
- Write the description in imperative form, starting with a lowercase word,
  without a trailing period. Keep the complete subject at most 72 characters.
- Keep each commit focused on one coherent change. Choose its type from the
  actual change, rather than copying the branch prefix automatically.
- For breaking changes, add `!` immediately before `:` and include a
  `BREAKING CHANGE: <description>` footer explaining the required migration.
  Separate optional bodies and footers from the subject with blank lines.
- PR titles and squash commit subjects follow the same rules. When merging,
  use squash or rebase so the resulting commits retain conventional messages.
  Revert messages must also use the `revert:` format; reference the reverted
  commit in the body or footer.

### Branch names

- Name every new work branch `<type>-<short-description>`, using a type from
  the table above that describes the branch's primary purpose.
- Use lowercase words and digits separated by single hyphens. Do not use
  slashes, underscores, spaces, personal prefixes, or a trailing hyphen.
- Keep the description short and specific, for example `feat-confidence-voting`,
  `fix-missing-user`, `test-service-components`, or `docs-git-conventions`.
- The permanent default branch (`master`) is exempt. These rules apply to new
  work; do not rename existing branches or rewrite historical commits merely
  to adopt the convention.

Before creating a branch or commit, check its name or subject against these
rules. Before opening or merging a PR, check its title and resulting commit
message as well.

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
