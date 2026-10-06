# Home Archery

Household scoring for Tony and Becky: practice sessions, per-arrow scores, totals, history, and personal bests.

v1 stack: Postgres, Go API (OpenAPI), SvelteKit PWA. Design artefacts:

- [Architecture and DDD](docs/architecture.md)
- [PWA wireframes](docs/wireframes.md)
- [Schema](db/migrations/001_create_schema.sql)
- [HTTP contract](spec/openapi.yaml)

Implementation uses TDD (order in the architecture doc). The Go domain, application use cases, HTTP adapter, and Postgres adapter exist. The PWA comes later.

Containers use **Podman** (`podman compose`), not Docker.

## API (this slice)

Prerequisite: Go 1.25+ ([api/go.mod](api/go.mod)).

### Default tests (no database)

Run from the module directory (`api/`, where `go.mod` lives — not the repo root):

```bash
cd api
go test ./...
```

Without `DATABASE_URL`, domain, application, memory, and HTTP run; Postgres adapter tests **SKIP** (`DATABASE_URL not set`). That is success. Verbose: `go test ./... -v`.

zsh `CORRECT` may rewrite `./...` to `./..`. Answer `n` at the prompt, or add `alias go='nocorrect go'` / `unsetopt CORRECT`.

### Postgres adapter tests

Needs Podman. On macOS, start the machine first. Commands below are from the **repo root** unless noted.

1. Start the engine (ignore “already running”):

```bash
podman machine start
```

2. Start Postgres (do not rely on `--wait`; Homebrew `podman-compose` often fails it):

```bash
podman compose -f deploy/compose/compose.yaml up -d
```

3. Wait until the service is healthy. `podman ps` should show port `5432` and `healthy`. Or:

```bash
podman exec compose_postgres_1 pg_isready -U archery -d home_archery
```

4. From `api/`:

```bash
cd api
DATABASE_URL='postgres://archery:archery@localhost:5432/home_archery?sslmode=disable' go test ./internal/adapters/postgres -v
```

5. Expect `PASS` for archer list/get, save/get/list, arrows and projections, and complete. Seeded Tony and Becky come from [db/seed/001_archers.sql](db/seed/001_archers.sql). Init SQL runs **only on an empty volume**. If archers are missing, recreate the volume (this destroys local DB data):

```bash
podman compose -f deploy/compose/compose.yaml down -v
podman compose -f deploy/compose/compose.yaml up -d
```

If Compose cannot start, check the Podman machine is running. Overlay or storage errors (`readlink … overlay`) are a local Podman issue, not a missing Go test.

### Serve

```bash
cd api
go run ./cmd/api
```

Default listen address is `:8080` (`HTTP_ADDR` overrides). Without `DATABASE_URL` the process is in-memory (restart wipes sessions). With the same `DATABASE_URL` as the adapter tests, the log should say `postgres`. Then `GET /api/v1/healthz` and `GET /api/v1/archers`.
