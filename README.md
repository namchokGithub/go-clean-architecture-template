# [project_name] Backend API

Go backend service for [project_name] platform. Built with Echo, GORM, PostgreSQL, and clean architecture (hexagonal / ports & adapters).

## Architecture

```
HTTP → Echo Router → Middleware → Handler → Service (port) → Repository (port) → GORM → PostgreSQL
```

Dependency rule: everything depends inward. Domain has zero knowledge of HTTP or DB. See [ARCHITECTURE.md](docs/architechture/ARCHITECTURE.md) for full reference.

## Engineering Docs

- Architecture: [ARCHITECTURE.md](docs/architechture/ARCHITECTURE.md)
- Agent rules: [AGENTS.md](docs/AGENTS.md)
- Project context: [CONTEXT.md](docs/CONTEXT.md)
- API guidelines: [API_GUIDELINES.md](docs/engineering/API_GUIDELINES.md)
- Testing strategy: [TESTING.md](docs/engineering/TESTING.md)

## Quick Start

### Prerequisites

- Go 1.24+
- Docker + Docker Compose

### 1. Clone and copy env

```bash
cp .env.example .env
# edit .env with your local values
```

### 2. Start dependencies

```bash
docker-compose up -d postgres redis
```

### 3. Run the server

```bash
go run ./cmd/main.go serve
```

Server listens on `APP_PORT` (default `8080`).

### 4. Health check

```bash
curl http://localhost:8080/health
```

## Docker (full stack)

```bash
docker-compose up --build
```

## API Documentation

OpenAPI contract lives in `docs/openapi/openapi.yaml`.

- source of truth for request/response contract
- update contract for every new endpoint
- keep reusable schemas in `docs/openapi/components/`
- keep path fragments in `docs/openapi/paths/`

## Environment Variables

| Variable        | Default     | Description                              |
| --------------- | ----------- | ---------------------------------------- |
| `APP_ENV`       | `local`     | Environment (`local` disables cron jobs) |
| `APP_PREFIX`    | `SDD`       | Prefix for error codes (e.g.`SDD-4001`)  |
| `APP_PORT`      | `8080`      | HTTP listen port                         |
| `LOG_LEVEL`     | `info`      | Logrus level                             |
| `POSTGRES_HOST` | `localhost` | PostgreSQL host                          |
| `POSTGRES_DB`   | `database`  | Database name                            |
| `JWT_SECRET`    | —           | HMAC secret for JWT validation           |
| `X_API_KEY`     | —           | API key for external routes              |
| `REDIS_ADDR`    | —           | Redis address (`host:port`)              |

See `.env.example` for the full list.

## Adding a Feature

Follow this sequence every time (enforced by architecture):

1. `internal/core/domain/` — add entity struct
2. `internal/core/port/<feature>.go` — add repository interface
3. `internal/core/port/<feature>.go` — add service interface
4. `internal/repository/<feature>/` — implement repository
5. `internal/repository/repository.go` — register in aggregate
6. `internal/core/service/<feature>/` — implement service
7. `internal/core/service/service.go` — register in aggregate
8. `internal/handler/<feature>/` — implement handler
9. `internal/handler/handler.go` — register in aggregate + routes
10. `protocol/http.go` — add route to correct group

## Response Format

**Success:**

```json
{ "is_error": false, "code": "SDD-200", "message": "Success", "data": {} }
```

**Error:**

```json
{
  "is_error": true,
  "code": "SDD-4001",
  "message": "Bad Request",
  "request_id": "...",
  "error_type": "BAD_REQUEST",
  "errors": []
}
```

Error code convention: `{PREFIX}-{CODE}` where first 3 digits of CODE = HTTP status.

## Pagination

All list endpoints use `POST /query` with this body:

```json
{
  "page": 1,
  "page_size": 20,
  "sort_name": "created_at",
  "sort_by": "desc",
  "search": "",
  "is_all": false
}
```

## Running Tests

```bash
go test ./...
```

Testing direction:

- unit tests stay near features as `*_test.go`
- generated mocks live in `tests/mocks/`
- integration tests live in `tests/integration/`
- generate mocks from `internal/core/port` with `mockery`

## Build

```bash
go build -o server ./cmd/main.go
./server serve
```
