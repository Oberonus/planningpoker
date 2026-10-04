[![codecov](https://codecov.io/gh/Oberonus/planningpoker/branch/master/graph/badge.svg?token=USXNP5KBW0)](https://codecov.io/gh/Oberonus/planningpoker)

# Planning Poker

## What is it?

This is a simple implementation of a Planning Poker tool which is intended to help the Agile 
teams.

According to [wikipedia](https://en.wikipedia.org/wiki/Planning_poker), Planning poker, also called Scrum poker, is a consensus-based, gamified 
technique for estimating, mostly used for timeboxing in Agile principles. In planning poker, members 
of the group make estimates by playing numbered cards face-down to the table, instead of speaking 
them aloud. The cards are revealed, and the estimates are then discussed. By hiding the figures 
in this way, the group can avoid the cognitive bias of anchoring, where the first number spoken 
aloud sets a precedent for subsequent estimates.

## How to run

The service is fully dockerized. Compose runs the app and PostgreSQL locally:
```bash
docker-compose up
```

When the container is up and running, just visit the service at [http://localhost:8080](http://localhost:8080)

## Deploy with Kamal

The live application is at https://poker.deminworks.com.
Deployments run from an operator's computer using `./bin/deploy`; GitHub Actions
only runs verification. The tracked Kamal configuration reads deployment details
from environment variables, so production addresses and resource names stay local.

The deployment SSH user needs Docker access without sudo. On a new server, Kamal
can install Docker when run with root SSH access; otherwise install Docker and
add the deployment user to the `docker` group first, then reconnect.

Install Kamal globally with Ruby 3.2 or newer:

```bash
gem install kamal --version 2.12.0 --no-document
# If using rbenv:
rbenv rehash
```

Create the local deployment settings:

```bash
cp config/deploy.env.example .env.deploy
chmod 600 .env.deploy
```

Edit `.env.deploy` with your SSH host/user, application domain, and shared database
resource names. The wrapper loads this ignored file from the project root and
exports its values before running Kamal. It is sourced as shell code, so keep it
under your control and quote values as needed. If the file is absent, the same
variables can be supplied through the exported environment.

For an existing database, preserve its service, container, and volume names.
Changing them can select a different container or empty volume. Keep private
operator notes in `.kamal/operations.md`, which is also ignored, and back up the
local configuration and secrets privately.

Point the application domain at your server's public IP, with ports 80 and 443
reachable. Keep the exact DNS and SSH details in your private operator notes.

Before the first deployment with PostgreSQL, follow the
[database setup guide](docs/postgres.md#initial-production-setup) to boot the shared
Kamal accessory, provision the app's database/login, and configure local secrets.

Run the first deployment from this directory while connected to Tailscale,
with a local Docker engine running:

```bash
./bin/deploy setup
```

For subsequent deployments and operations:

```bash
./bin/deploy
./bin/deploy app details
./bin/deploy app logs
./bin/deploy rollback <previous-version>
```

Kamal builds the committed Git revision for amd64. Commit deployment and
application changes before running `./bin/deploy`.
It runs a local registry on port 5555 and transfers images through SSH tunnels,
so external registry credentials are unnecessary. The proxy serves HTTPS and
checks `/alive` on application port 8080 before switching traffic.

Games and users are persisted in PostgreSQL. The app requires `DATABASE_URL` and
applies its schema migrations at startup. PostgreSQL is managed separately by
Kamal through `config/database.yml`, so other services can share the database
server with separate databases and credentials. Follow the
[PostgreSQL setup and operations guide](docs/postgres.md) before the first deploy.

## Development

This service is built with Domain Driven Design, CQRS, event based communication, clean code and
hexagonal architecture in mind (yep, a lot of buzzwords here).

Backend and frontend logic live in a single binary, with PostgreSQL persistence.
Go 1.23 or newer is required for backend development.

Languages used:
- Backend - Go
- Frontend - Vue.js

Protocols:
- Plain and simple HTTP
- Socket.io v1.x on top of the websocket protocol

### Overall architecture

The golang project is structured pretty much in a standard way, but a couple of things are worth to be mentioned:
- All the business logic lives in `internal/domain` directory
- All the infra related logic (http, websocket, database, etc...) lives in `internal/infra` directory.
- The frontend lives in `web` directory.

Golang app serves the frontend by itself, so no additional layer (e.g. NGINX) is needed.

#### Architecture diagram: 

<img alt="Architecture" src="docs/architecture.png" />

The event bus remains in memory. Games are stored as JSONB documents and users
in a relational table; memory repositories remain available for unit tests.

### How it works

HTTP layer is used only to serve some basic auth related requests, all the game logic
utilizes the websocket protocol.

The best way to understand how things are working, is to dive deep in the codebase, but I believe 
following diagrams might make this process a bit easier.

### Actions
Any action in the gameplay are processed in the same way. 

Actions are:
- Create a new game
- Update game parameters
- Vote
- Un-vote
- Reveal cards (finish the game)
- Restart the game
- etc...

This sequence diagram describes how the process is working in general:

<img alt="Actions" src="docs/actions.png" />

### Domain events processing

For now there are just a couple of domain events exist in the system:
- Player action (described in the previous section)
- Player changed name

Each event should be propagated to all players in order to reflect changes and display an actual
state of the game. 

<img alt="Game Updated" src="docs/game_updated.png" />

### Frontend development

There is a possibility to run automatic watcher/builder for frontend:
```bash
./mage.sh devFront
```
It will spinup `node:14-alpine` container, mount all frontend codebase and
build right inside the container, avoiding any host dependencies.

### Verification

Repository-wide coding conventions are in [AGENTS.md](AGENTS.md). The style check
enforces import grouping, blank-line spacing, and domain dependency boundaries:

```bash
golangci-lint run --config .golangci-style.yml ./...
```

Use golangci-lint v1.61.0, matching CI.

CI runs Go vet, unit and PostgreSQL integration tests with the race detector,
and a container smoke test. For local checks:

```bash
go vet -mod=vendor ./...
go test -mod=vendor -race ./...
```

See the [PostgreSQL guide](docs/postgres.md#local-development-and-verification)
for integration tests using a disposable database.

## Contribution
Your contribution is very welcomed! Please create pull request or issue.
