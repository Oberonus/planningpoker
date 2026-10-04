# PostgreSQL for self-hosting

Planning Poker stores game aggregates as JSONB documents and user identities in
ordinary columns. It requires `DATABASE_URL` at startup and fails if PostgreSQL
is unavailable. Schema migrations are embedded in the binary and applied in a
transaction under a database advisory lock, so simultaneous starts are safe.
Provision a dedicated database and login for the application. That login can
migrate its own tables and does not need PostgreSQL administrator privileges.

Game updates use `SELECT ... FOR UPDATE` through commit. This prevents lost votes
across independent connections. Events are published after commit using the
existing in-process bus. `/alive` checks database connectivity before returning
200, allowing Kamal to check readiness before moving traffic.

The event bus and WebSocket connections live in one process. Use a single
application instance with the supplied configuration; multiple replicas would
also need event propagation between nodes. There is no automatic import from
an in-memory instance into PostgreSQL.

## Default PostgreSQL configuration

The included `config/database.yml` runs PostgreSQL as a separate Kamal accessory
on the same server as the application. Its defaults are:

- Container: `DEPLOY_POSTGRES_CONTAINER`, on the shared `kamal` Docker network.
- Image: `postgres:18.6`, pinned to a patch version; update it deliberately.
- Data: named Docker volume `DEPLOY_POSTGRES_VOLUME`, mounted at
  `/var/lib/postgresql` as required by the PostgreSQL 18 image.
- Resources: 1 GB container limit, 256 MB shared buffers, 60 connections.
  Planning Poker's connection pool uses at most 10 connections.
- Host port: `127.0.0.1:5432`, for access through an SSH tunnel.

Adjust the image, memory settings, connection limits, and port binding in
`config/database.yml` for your server. Application deployments and removal do not
manage this accessory. If you want other applications to share it, give each one
its own database and login. Keep the accessory definition in one configuration
so those applications do not independently manage the same container.

Choose the database service name with `DEPLOY_DATABASE_SERVICE`. All deployment
values come from the ignored `.env.deploy` or exported environment. Use
`./bin/deploy database ...` to load them and select this configuration. Preserve
existing service, container, and volume names when reusing an existing database;
renaming them can select different resources or empty storage.

## Provision your database

Configure `.env.deploy` as described in the
[deployment guide](../README.md#deploy-on-your-infrastructure).
Create the local secrets directory:

```sh
mkdir -p .kamal
```

Put the application's connection URL in `.kamal/secrets` and restrict the file
to your user:

```sh
DATABASE_URL=postgres://planningpoker:YOUR_APP_PASSWORD@YOUR_POSTGRES_CONTAINER:5432/planningpoker?sslmode=disable
```

```sh
chmod 600 .kamal/secrets
```

For the supplied accessory, replace `YOUR_POSTGRES_CONTAINER` with the value of
`DEPLOY_POSTGRES_CONTAINER` from `.env.deploy`. If you already manage PostgreSQL,
use your instance's connection URL, including its database, login, and connection
security settings. The application's login must be able to create and migrate
its own tables. Skip the accessory steps below and continue to
[deploying the application](#deploy-the-application).

### Start the included PostgreSQL accessory

Create `.kamal/database-secrets` with the accessory's administrator password:

```sh
POSTGRES_PASSWORD=YOUR_RANDOM_ADMIN_PASSWORD
```

Use different random passwords of at least 16 characters. Percent-encode special
characters in the connection URL's password. The provisioning command below
takes the original password, not its URL-encoded representation. `sslmode=disable`
is for the local Docker network on this host.

```sh
chmod 600 .kamal/secrets .kamal/database-secrets
./bin/deploy database accessory boot postgres
./bin/deploy database accessory exec postgres --reuse 'pg_isready -U postgres'
./bin/deploy database accessory exec postgres --interactive --reuse \
  'bash /opt/postgres/provision.sh planningpoker'
```

Enter the app password at the hidden prompt. The script creates a login without
superuser, database creation, role creation, or replication privileges. It also
revokes public access to the new database and its public schema. It refuses to
overwrite an existing login/database. If provisioning fails partway through,
inspect the state in the administrator console before retrying.

### Deploy the application

Kamal builds the committed Git revision. Commit your application and
configuration changes, then set up the first application deployment:

```sh
./bin/deploy setup
```

Use `./bin/deploy` for subsequent application deployments. The application
automatically creates its tables using its own credentials.

## Share the database server with another application

If you choose to share the accessory, provision a separate database and login
for each additional application:

```sh
./bin/deploy database accessory exec postgres --interactive --reuse \
  'bash /opt/postgres/provision.sh another_service'
```

Give that service a connection URL with its own login, password, and database:

```text
postgres://another_service:PASSWORD@YOUR_POSTGRES_CONTAINER:5432/another_service?sslmode=disable
```

Use a lowercase service name containing only letters, digits, and underscores.
Each future service shares the server's CPU, RAM, disk, and connection capacity.

The supplied configuration assumes the app and PostgreSQL share one server's
Docker network. Docker container hostnames do not work across servers. For a
database on another server, configure its network access and use an address
reachable from the app through your private network. Keep database access off
the public interface.

## Administration and backups

The commands below apply to the supplied Kamal accessory. For a separately
managed database, use its administration and backup tools instead.

```sh
./bin/deploy database accessory details postgres
./bin/deploy database accessory logs postgres
./bin/deploy database accessory exec postgres --interactive --reuse \
  'psql -U postgres -d postgres'
```

Changing the administrator secret file does not rotate the password of an existing
database. Use `\password postgres` in the console, then update the secret file.
Use `\password planningpoker` and update `DATABASE_URL` to rotate the app password.

Export a database backup through SSH without copying credentials into commands:

```sh
set -a
. ./.env.deploy
set +a
ssh "${DEPLOY_SSH_USER}@${DEPLOY_HOST}" \
  "docker exec ${DEPLOY_POSTGRES_CONTAINER} pg_dump -U postgres -Fc planningpoker" \
  > /path/to/private/backups/planningpoker.dump
```

Keep backups off the node and protect them as application data. A persistent
volume preserves data across container replacement, but is not a backup. Keep
dump files outside the repository. Before database upgrades, make and test a
backup. A patch upgrade can use
`./bin/deploy database accessory reboot postgres` after changing the image
tag; it briefly interrupts all services using this server. A major version upgrade
requires PostgreSQL's upgrade or dump/restore procedure, not just a new image tag.

To test restoration, create a separate database from the admin console and
restore into it without overwriting your application's database:

```sh
# With .env.deploy exported as in the backup example above:
ssh "${DEPLOY_SSH_USER}@${DEPLOY_HOST}" \
  "docker exec -i ${DEPLOY_POSTGRES_CONTAINER} pg_restore -U postgres --no-owner -d planningpoker_restore" \
  < /path/to/private/backups/planningpoker.dump
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
