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
