```myapp/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── user/
│   ├── order/
│   ├── httpx/                    # renamed from httpserver, per earlier discussion
│   ├── db/
│   │   ├── db.go
│   │   └── migrations/
│   └── config/
│       └── config.go
├── .env.example
├── .env                          # gitignored
├── go.mod
└── go.sum
```

## API Response Contract

### Success responses (`2xx`)

- Content-Type: `application/json`
- Body: bare resource payload (no top-level envelope)

Examples:

- `GET /users` -> `200` with `[{"id":"...","name":"...","email":"..."}]`
- `POST /users` -> `201` with `{"id":"...","name":"...","email":"..."}`

### Error responses (`4xx` / `5xx`)

- Content-Type: `application/problem+json`
- Body: RFC 7807 Problem Details with required fields:
  - `type`, `title`, `status`, `detail`, `instance`
- Extension members:
  - `request_id` (always present)
  - `errors` (validation-only field map)

Problem type URLs:

- `https://api.<your-domain>/problems/validation-error`
- `https://api.<your-domain>/problems/duplicate-resource`
- `https://api.<your-domain>/problems/not-found`
- `https://api.<your-domain>/problems/internal-error`

## Breaking Change Notice

Clients that previously parsed legacy error shapes such as `{"error":"..."}` or top-level `{"errors": {...}}` must migrate to RFC 7807 parsing immediately.

## OAuth2 Auth Endpoints

### Required environment variables

- `OAUTH_CLIENT_ID`
- `OAUTH_CLIENT_SECRET`
- `OAUTH_ACCESS_TOKEN_TTL_MINUTES` (default: `120`)
- `OAUTH_REFRESH_TOKEN_TTL_HOURS` (default: `24`)

### Signup

- `POST /auth/signup`
- JSON body:

```json
{
  "name": "Ada",
  "email": "ada@example.com",
  "password": "Password1"
}
```

### Token issuance (password grant)

- `POST /oauth/token`
- Form-encoded body:

```text
grant_type=password&username=ada@example.com&password=Password1&client_id=<OAUTH_CLIENT_ID>&client_secret=<OAUTH_CLIENT_SECRET>
```

### Authenticated user profile

- `GET /auth/me`
- Header:

```text
Authorization: Bearer <access_token>
```

## Database Schema Workflow (Struct-First)

Schema ownership is split:

- `internal/db/persistence`: canonical persistence schema structs (declarative only, no business hooks)
- `internal/model`: domain models used by handlers/services
- `internal/db/migrations`: reviewed SQL migrations applied in runtime/deploy

Generation boundary:

- `internal/model` is **manual** (not generated). Keep it focused on domain/business-facing shape.
- `internal/db/persistence` is **manual** and is the schema source for Atlas diff.
- `internal/db/query` is **generated** by `gorm/gen` (`make gen-query`).

Relationship guidance for persistence schema:

- Add relationship fields in `internal/db/persistence` when you need foreign key constraints generated and enforced.
- Scalar FK fields alone are not enough for this workflow unless the relationship/constraint is also explicitly represented.
- Keep relationship fields in persistence structs; keep domain relationships in `internal/model` only when the service/handler needs them.

Production schema updates use SQL migrations only. Runtime AutoMigrate is not used.

### Atlas in Docker Compose

Atlas is available as a pinned Docker service in `deploy/docker-compose.yml`.

Useful commands:

- `make atlas-version`
- `make atlas-schema`
- `make atlas-diff name=<migration_name>`
- `make check-schema-drift`

`make atlas-diff` keeps the existing numeric migration naming convention by creating a sequential migration first and writing Atlas diff SQL into the generated `.up.sql` file.

### gorm/gen typed query generation

- Generator entrypoint: `cmd/querygen/main.go`
- Generated output: `internal/db/query`

Run and verify:

- `make gen-query`
- `make verify-gen-query`

`make verify-gen-query` is intended for CI to enforce deterministic generated output.

### Adding a new table/feature (example: `tickets`)

Yes, start from `internal/db/persistence`.

1. Model persistence shape first

- Add `Ticket` (and relations/indexes/FKs) in `internal/db/persistence`.
- Keep this layer declarative only (no business hooks/logic).

2. Register new persistence models in loaders (required)

- Add the new structs to Atlas loader in `cmd/atlasloader/main.go` (`gormschema.New(...).Load(...)`).
- Add the same structs to query generation in `cmd/querygen/main.go` (`g.ApplyBasic(...)`).

If a struct exists in `internal/db/persistence` but is not registered in these two entrypoints, Atlas diff and gorm/gen will ignore it.

3. Generate migration SQL from schema diff

- Ensure local stack is running: `make dev`
- Generate migration: `make atlas-diff name=create_tickets_table`

If there is no schema change, `make atlas-diff` now exits without creating migration files.

4. Review migration files

- Review generated `*.up.sql` in `internal/db/migrations`.
- Review generated `*.down.sql` as well (it is auto-generated from reverse diff); adjust manually if needed.

5. Apply migration

- Apply with existing migration flow: `make migrate`

6. Regenerate typed query layer

- `make gen-query`

7. Implement feature layers

- Repository: persistence-only operations with generated query package (`internal/db/query`).
- Service: domain behavior and `model <-> persistence` conversion.
- Handler/API types/routes: request/response, validation, transport concerns.

8. Validate end-to-end

- `make feature-verify`

`make feature-verify` runs drift checks, deterministic codegen checks, and test suite:

- `make check-schema-drift`
- `make verify-gen-query`
- `make test`

### Troubleshooting: migration file not generated

If `make atlas-diff name=...` does not create a migration for your new table:

- Confirm the model is in `internal/db/persistence`.
- Confirm it is registered in both:
  - `cmd/atlasloader/main.go`
  - `cmd/querygen/main.go`
- Run `make atlas-schema` and verify the table appears in the printed SQL.
- Then run `make atlas-diff name=<migration_name>` again.

### Quick reference (new/updated commands)

- `make test`
- `make feature-verify`
