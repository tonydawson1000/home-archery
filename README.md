# Home Archery

Household scoring for Tony and Becky: practice sessions, per-arrow scores, totals, history, and personal bests.

v1 stack: Postgres, Go API (OpenAPI), SvelteKit PWA. Design artefacts:

- [Architecture and DDD](docs/architecture.md)
- [PWA wireframes](docs/wireframes.md)
- [Schema](db/migrations/001_create_schema.sql)
- [HTTP contract](spec/openapi.yaml)

Implementation uses TDD, starting in `api/internal/domain`. The API and PWA are not implemented yet.
