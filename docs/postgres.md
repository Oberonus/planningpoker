# PostgreSQL storage and operations

Planning Poker stores game aggregates as JSONB documents and user identities in
ordinary columns. It requires `DATABASE_URL` at startup and fails if PostgreSQL
is unavailable. Schema migrations are embedded in the binary and applied in a
transaction under a database advisory lock, so simultaneous starts are safe.
The application login owns only the `planningpoker` database and can migrate its
own tables; it is not a PostgreSQL administrator.

Game updates use `SELECT ... FOR UPDATE` through commit. This prevents lost votes
across independent connections. Events are published after commit using the
existing in-process bus. `/alive` checks database connectivity before returning
200, allowing Kamal to check readiness before moving traffic.

The current event bus and WebSocket connections still live in one process. This
change provides persistence and safe database writes, but multiple app replicas
would also need event propagation between nodes. Existing in-memory games/users
cannot be imported automatically; the first deployment starts a fresh database.

## Shared database server

`config/database.yml` is the only Kamal configuration that owns PostgreSQL:

- Container: `deminworks-postgres`, on the shared `kamal` Docker network.
- Image: `postgres:18.6`, pinned to a patch version; update it deliberately.
- Data: named Docker volume `deminworks-postgres-data`, mounted at
  `/var/lib/postgresql` as required by the PostgreSQL 18 image.
- Resources: 1 GB container limit, 256 MB shared buffers, 60 connections.
  Planning Poker's connection pool uses at most 10 connections.
- Host port: `127.0.0.1:5432`, for access through an SSH tunnel.

Application deployments and removal do not manage this accessory. Other services
should use the same container hostname, each with its own database and login;
do not copy the accessory definition into their application configs.

## Initial production setup

Create two local files, both excluded from Git, and restrict them to your user:

`.kamal/database-secrets`:

```sh
POSTGRES_PASSWORD=YOUR_RANDOM_ADMIN_PASSWORD
```

`.kamal/secrets`:

```sh
DATABASE_URL=postgres://planningpoker:YOUR_APP_PASSWORD@deminworks-postgres:5432/planningpoker?sslmode=disable
```

Use different random passwords of at least 16 characters. Percent-encode special
characters in the connection URL's password. The provisioning command below
takes the original password, not its URL-encoded representation. `sslmode=disable`
is for the local Docker network on this host.

```sh
chmod 600 .kamal/secrets .kamal/database-secrets
kamal accessory boot postgres -c config/database.yml
kamal accessory exec postgres -c config/database.yml --reuse 'pg_isready -U postgres'
kamal accessory exec postgres -c config/database.yml --interactive --reuse \
  'bash /opt/postgres/provision.sh planningpoker'
```

Enter the app password at the hidden prompt. The script creates a login without
superuser, database creation, role creation, or replication privileges. It also
revokes public access to the new database and its public schema. It refuses to
overwrite an existing login/database. If provisioning fails partway through,
inspect the state in the administrator console before retrying.

Commit the implementation before deploying: Kamal builds the committed Git
revision. Then run:

```sh
kamal deploy
```

The application automatically creates its tables using its own credentials.

## Add another service

```sh
kamal accessory exec postgres -c config/database.yml --interactive --reuse \
  'bash /opt/postgres/provision.sh another_service'
```

Give that service a connection URL with its own login, password, and database:

```text
postgres://another_service:PASSWORD@deminworks-postgres:5432/another_service?sslmode=disable
```

Use a lowercase service name containing only letters, digits, and underscores.
Each future service shares the server's CPU, RAM, disk, and connection capacity.

Apps on other nodes cannot use this Docker hostname. Before adding remote apps,
bind PostgreSQL to the server's Tailscale address, configure access controls, and
use its private address in their URLs; keep it off the public interface.

## Administration and backups

```sh
kamal accessory details postgres -c config/database.yml
kamal accessory logs postgres -c config/database.yml
kamal accessory exec postgres -c config/database.yml --interactive --reuse \
  'psql -U postgres -d postgres'
```

Changing the administrator secret file does not rotate the password of an existing
database. Use `\password postgres` in the console, then update the secret file.
Use `\password planningpoker` and update `DATABASE_URL` to rotate the app password.

Export a database backup through SSH without copying credentials into commands:

```sh
ssh alex@deminworks 'docker exec deminworks-postgres pg_dump -U postgres -Fc planningpoker' \
  > planningpoker.dump
```

Keep backups off the node and protect them as application data. A persistent
volume preserves data across container replacement, but is not a backup. Before
database upgrades, make and test a backup. A patch upgrade can use
`kamal accessory reboot postgres -c config/database.yml` after changing the image
tag; it briefly interrupts all services using this server. A major version upgrade
requires PostgreSQL's upgrade or dump/restore procedure, not just a new image tag.

To verify restoration without overwriting production, create a separate database
from the admin console and restore into it:

```sh
ssh alex@deminworks \
  'docker exec -i deminworks-postgres pg_restore -U postgres --no-owner -d planningpoker_restore' \
  < planningpoker.dump
```

Do not run accessory removal or delete its volume as part of app cleanup.

## Local development and verification

`docker compose up --build` starts the application and a persistent local database,
exposed on host port 54338 to avoid other projects' PostgreSQL instances.
The Compose credentials are only for local development. To run Go outside Docker:

```sh
docker compose up -d postgres
DATABASE_URL='postgres://planningpoker:local-development-only@127.0.0.1:54338/planningpoker?sslmode=disable' \
  go run -mod=vendor ./cmd/poker
```

Go 1.23 or newer is required. Unit tests use memory repositories. Integration tests
require a separate disposable database; they create data and apply migrations:

```sh
go test -mod=vendor -race ./...
TEST_DATABASE_URL='postgres://USER:PASSWORD@127.0.0.1:5432/DISPOSABLE_DB?sslmode=disable' \
  go test -mod=vendor -race -tags=integration ./internal/infra/repository
```

The integration suite covers simultaneous startup/migration, games and identities
surviving reconnects, concurrent votes through separate connection pools, rollback,
and events being emitted after commit.
