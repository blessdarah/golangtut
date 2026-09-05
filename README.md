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
