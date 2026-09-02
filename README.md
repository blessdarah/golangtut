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
