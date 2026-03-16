# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build
go build ./...

# Lint (requires golangci-lint v1.64.8+)
golangci-lint run

# Run tests (no tests exist yet)
go test ./...

# Run a single test
go test -run TestName ./internal/...

# Regenerate protobuf/gRPC code (after editing auth.proto)
protoc --go_out=. --go-grpc_out=. internal/gen/auth.proto
```

## Architecture

This is a **gRPC authentication microservice** with a 3-layer architecture:

```
Controller (gRPC) → Logic → DB / JWTManager / Redis
```

**Layers:**
- `internal/controller/` — gRPC server; implements `Register`, `Login`, `ValidateToken` RPCs defined in `internal/gen/auth.proto`. Generated code lives in `internal/gen/proto/`.
- `internal/logic/` — business logic; orchestrates DB, JWT, and Redis calls.
- `internal/db/` — PostgreSQL client; SQL queries are constants in `queries.go`. Users table: `id`, `first_name`, `middle_name`, `last_name`, `login`, `password`.
- `internal/jwt_manager/` — creates and validates JWT tokens (HS256, claims: `user_id`, `exp`).
- `internal/redis/` — token cache with TTL; methods: `SaveToken`, `GetTokenByUserID`.
- `internal/utils/` — bcrypt password hashing.
- `internal/config/` — reads required env vars; panics on missing values.
- `internal/model/` — shared structs used across layers.

## Required Environment Variables

| Variable | Description |
|---|---|
| `HOST` | gRPC listen address (e.g. `localhost:50051`) |
| `DB_CONNECTION_STRING` | PostgreSQL DSN |
| `SECRET_KEY` | JWT signing secret |
| `TOKEN_DURATION` | Token expiry in seconds |
| `REDIS_SERVERS` | Comma-separated Redis addresses |
| `REDIS_PASSWORD` | Redis password |
| `REDIS_TOKEN_TTL` | Token cache TTL in minutes |

## Docker

Multi-stage Alpine build in `docker/Dockerfile`. Runs as non-root `appuser`. Default exposed port: `8080`. Build args: `GO_VERSION`, `ALPINE_VERSION`, `PORT`.

## Linter

`.golangci.yml` enables `errcheck`, `govet`, `staticcheck` with a 5-minute timeout. CI uses golangci-lint `v1.64.8`.
