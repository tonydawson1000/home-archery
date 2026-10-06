# Home Archery

Household scoring for Tony and Becky: practice sessions, per-arrow scores, totals, history, and personal bests.

v1 stack: Postgres, Go API (OpenAPI), SvelteKit PWA. Design artefacts:

- [Architecture and DDD](docs/architecture.md)
- [PWA wireframes](docs/wireframes.md)
- [Schema](db/migrations/001_create_schema.sql)
- [HTTP contract](spec/openapi.yaml)

Implementation uses TDD (order in the architecture doc). This slice has the Go domain, application use cases against in-memory ports, and an HTTP adapter. The PWA and Postgres adapter come later.

## API (this slice)

Tests do **not** need Compose or a database. Prerequisite: Go 1.23+ ([api/go.mod](api/go.mod)).

Run tests from the module directory (`api/`, where `go.mod` lives — not the repo root):

```bash
cd api
go test ./...
```

Verbose: `go test ./... -v`. One layer: `./internal/domain`, `./internal/application/...`, or `./internal/adapters/http`.

zsh `CORRECT` may rewrite `./...` to `./..`. Answer `n` at the prompt, or add `alias go='nocorrect go'` / `unsetopt CORRECT`.

Expect `ok` for domain, application, memory, and HTTP. `cmd/api` and `adapters/postgres` with no test files is expected.

Optional smoke (in-memory process; restart wipes sessions):

```bash
cd api
go run ./cmd/api
```

Default listen address is `:8080` (`HTTP_ADDR` overrides). Then `GET /api/v1/healthz` and `GET /api/v1/archers`. Seeded Tony and Becky UUIDs match [db/seed/001_archers.sql](db/seed/001_archers.sql).
