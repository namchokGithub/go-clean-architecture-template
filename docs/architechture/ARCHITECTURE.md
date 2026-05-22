# Backend Architecture Reference

> Production-oriented architecture doc for `st-driver`. Use as a foundation template for new Golang API projects.

---

## Architecture Overview

### Style: Clean Architecture (Hexagonal / Ports & Adapters)

Strict dependency inversion. Domain core has zero knowledge of infrastructure. All external concerns (DB, HTTP, Redis, S3) depend inward.

```
HTTP Request
    → Echo Router (protocol/)
    → Middleware Chain (auth → afterAuth → permission → userAction)
    → Handler (internal/handler/[feature]/)
    → Service Port interface (internal/core/port/)
    → Service impl (internal/core/service/[feature]/)
    → Repository Port interface
    → Repository impl (internal/repository/[feature]/)
    → GORM → PostgreSQL
```

### Request Lifecycle

1. **Global middleware** — injects `Lang`, `RequestID`, `OS`, `Browser` into context
2. **JWT auth** — validates Bearer token, extracts `USER_ID`, `USER_TYPE`, `ROLE_ID`, `TOKEN`
3. **afterAuth** — loads user + role permissions from Redis cache; falls back to DB on miss; sets `USER` and `ROLE_PERMISSION` in context
4. **permission** — checks `ROLE_PERMISSION[menuID]` contains required permission string (currently pass-through, framework ready)
5. **userAction** — logs to `tbl_user_logs` via logrus DB hook
6. **Handler** — bind + validate request, call service, render response
7. **Error handler** — global Echo error handler maps domain errors to HTTP status via error code prefix

### Dependency Flow

```
configs/         → loaded once at startup, injected via Dependencies structs
infrastructure/  → creates *gorm.DB connection
pkg/             → shared technical packages; core may use only transport-neutral/pure helpers
internal/core/   → domain logic; imports only itself and transport-neutral packages
internal/handler → imports service ports; never imports repository
internal/repository → imports domain ports; imports gorm
protocol/        → wires everything together; creates Handler, Service, Repository
```

### Dependency Injection Pattern

All layers use explicit `Dependencies` struct + `NewXxx(deps)` constructors. No global service registries. No reflection-based DI.

```go
type Dependencies struct {
    DB     *gorm.DB
    Config configs.Config
}
func NewRepository(d Dependencies) *Repository { ... }
```

---

## Recommended Reusable Structure

```
project-root/
├── cmd/
│   ├── main.go                  # cobra root execution
│   └── cmds/
│       ├── root.go              # cobra root command
│       └── rest.go              # "serve" subcommand
├── protocol/
│   ├── init.go                  # application struct, bootstrap
│   ├── http.go                  # Echo setup, route registration
│   ├── cron.go                  # cron job registrations
│   ├── mobile_route.go          # mobile-specific route group
│   ├── custom_middleware.go     # afterAuth, setPublishContext, validateType
│   ├── permission_middleware.go # RBAC permission gate
│   ├── user_action_middleware.go# request logging middleware
│   └── external_middleware.go   # external API key middleware
├── configs/
│   ├── config.go                # envconfig struct + Init() + GetConfigs()
│   ├── postgres.go              # DB config struct
│   ├── key.go                   # JWT key config
│   ├── s3.go                    # S3 config struct
│   ├── email.go                 # Email config struct
│   ├── oauth.go                 # OAuth config
│   ├── ad.go                    # Active Directory config
│   └── version.go               # version info
├── infrastructure/
│   └── postgres.go              # *gorm.DB factory
├── internal/
│   ├── core/
│   │   ├── domain/              # structs shared across layers (entities, DTOs, errors)
│   │   ├── port/
│   │   │   ├── auth.go          # auth service + repository contracts
│   │   │   ├── user.go          # user service + repository contracts
│   │   │   ├── order.go         # order service + repository contracts
│   │   │   └── payment.go       # payment service + repository contracts
│   │   ├── service/
│   │   │   ├── service.go       # Service aggregate struct + NewService()
│   │   │   └── [feature]/       # one package per feature
│   │   ├── helper/              # pure business-aware helpers (no HTTP)
│   │   ├── constant/            # typed app constants
│   │   ├── enums/               # typed enum structs with ID/Name/NameEn
│   │   └── transaction/         # DB tx via context (NewContext / FromContext)
│   ├── handler/
│   │   ├── handler.go           # Handler aggregate struct + NewHandler()
│   │   ├── base/                # BaseResponse, BasePageRequest, BaseTransformer
│   │   ├── common/              # response constructors, error codes, error types
│   │   ├── middleware/          # error_handler, api_key_guard, content_type
│   │   ├── validator/           # custom Echo validator (go-playground/validator)
│   │   └── [feature]/           # one package per feature handler
│   └── repository/
│       ├── repository.go        # Repository aggregate struct + NewRepository()
│       └── [feature]/           # one package per feature repository
├── pkg/
│   ├── jwt/                     # JWT middleware + claims model
│   ├── logx/                    # logrus wrapper + DB hook (tbl_logs + tbl_user_logs)
│   ├── redis/                   # Redis client wrapper
│   ├── s3/                      # S3 interface + AWS SDK v2 impl
│   ├── client/                  # reusable HTTP client wrapper
│   └── contexts/                # context key type helpers
├── websocket/
│   └── websocket.go             # WebSocketManager (clients + groups maps)
├── mocks/                       # mockery-generated service mocks for tests
├── assets/
│   ├── email_templates/
│   ├── sql/                     # ad-hoc SQL migration scripts
│   ├── fonts/
│   ├── export/                  # export template files
│   └── import/                  # import template files
├── docs/
│   └── openapi/
│       ├── openapi.yaml         # source of truth
│       ├── README.md            # doc conventions + generation notes
│       ├── components/          # reusable schemas, parameters, responses
│       └── paths/               # endpoint path fragments
├── cicd/                        # CI/CD pipeline configs
├── vendor/                      # vendored dependencies
└── go.mod
```

### Layer Responsibilities

| Layer         | Owns                                                    | Must NOT                                         |
| ------------- | ------------------------------------------------------- | ------------------------------------------------ |
| `domain/`     | Entities, DTOs (base), shared error types               | Import HTTP, GORM, external SDKs                 |
| `port/`       | Service + Repository interfaces                         | Have any impl                                    |
| `service/`    | Business logic, orchestration                           | Import `handler`, `repository` packages directly |
| `repository/` | SQL queries via GORM                                    | Business logic, HTTP concerns                    |
| `handler/`    | HTTP binding, validation, response rendering            | Business logic beyond request parsing            |
| `helper/`     | Pure functions: string ops, math, format, context reads | Side effects, DB calls, HTTP calls               |
| `pkg/`        | Infrastructure wrappers                                 | Domain logic                                     |

### Port File Organization

Use one port file per feature under `internal/core/port/` instead of central `services.go` / `repository.go`.

Why:

- reduce merge conflicts
- improve feature isolation
- improve maintainability
- improve AI-agent navigation/readability

### helpers/ vs utils/ distinction

**Put in `internal/core/helper/`** — business-aware pure functions (e.g., `GetVatAmount`, `MapRolePermission`, `PrepareSearch`, `GetUserID`). May reference `domain` types and `constant` package.

**Put in `pkg/`** — shared technical packages that are project-agnostic. Pure helpers are acceptable; infrastructure adapters are also allowed, but domain/service code should only depend on transport-neutral packages that do not know about HTTP, DB, Redis, or external SDK clients.

**Do NOT create a generic `utils/` package.** Utility functions gravitate toward god-package syndrome. Place them in the layer where they're used.

---

## Reusable Patterns

### Standardized Response Models

**Success:**

```go
type DefaultResponse struct {
    IsError bool   `json:"is_error"` // always false
    Code    string `json:"code"`     // "APP-200"
    Message string `json:"message"`  // "Success" | "Created"
    Data    any    `json:"data"`
}
```

**Error:**

```go
type ErrorResponse struct {
    IsError   bool   `json:"is_error"` // always true
    Code      string `json:"code"`     // "APP-4001"
    Message   string `json:"message"`
    RequestID string `json:"request_id"`
    ErrorType string `json:"error_type"`
    Errors    any    `json:"errors"`   // []SubError or string
}
```

**Error code convention:** `{APP_PREFIX}-{ERRORCODE}` where the first 3 digits of ERRORCODE map to HTTP status (e.g., `4001` → 400, `5001` → 500).

**Use `common.NewXxx()` constructors** — never manually construct error responses in handlers.

### Pagination Pattern

Request (POST `/query` endpoints):

```go
type BasePageRequest struct {
    Page     int    `json:"page"`
    PageSize int    `json:"page_size"`
    SortName string `json:"sort_name"`
    SortBy   string `json:"sort_by"`   // "asc" | "desc"
    Search   string `json:"search"`
    IsAll    bool   `json:"is_all"`    // bypass pagination
}
```

Response wraps in `domain.Paginator` (includes `data`, `page`, `page_size`, `total`, `total_page`).

**Use POST for all paginated queries** — not GET with query params. This enables complex filter bodies.

### Auth Pattern

1. JWT token in `Authorization: Bearer <token>`
2. `jwt.AuthMiddleware` validates signature, populates context keys: `USER_ID`, `USER_TYPE`, `ROLE_ID`, `EMAIL`, `TOKEN`
3. `afterAuthMiddleware` fetches full user + role permissions; caches with key `{token}_user` and `{token}_rolePermission`
4. Handlers read user from context via `helper.GetUserID(ctx)` or `ctx.Value(domain.USER)`

**External routes** use API key middleware (`X-Api-Key` header checked against `X_API_KEY` env).

### Repository/Service Conventions

**Repository:**

- Receives `*gorm.DB` via `Dependencies.DB`
- Checks `transaction.FromContext(ctx)` first; uses tx if present, falls back to `r.db`
- Method names: `Create`, `FindByID`, `FindPaging`, `Update`, `Delete`, `FindAll`
- Returns domain types only, not GORM models

**Service:**

- Receives port interfaces (not concrete repository structs)
- Orchestrates cross-repo operations within a DB transaction
- Semaphore (`semaphore.NewWeighted(1)`) for mutex-like operations on shared counters (budget, coupon quota, etc.)

### Transaction Handling

```go
err := transaction.WithTx(ctx, db, func(ctx context.Context) error {
    // service logic
    return nil
})
```

Prefer `transaction.WithTx(...)` in service layer. Keep tx propagation in context, let repositories auto-detect tx from context. Centralized lifecycle gives automatic rollback, panic safety, less duplicated commit/rollback code.

Pass context through all layers — never pass `*gorm.DB` as a service parameter.

### API Documentation Strategy

- `docs/openapi/openapi.yaml` is source of truth
- Use Swagger/OpenAPI-first workflow compatible with Echo handlers
- Keep reusable schemas, parameters, and common response components under `docs/openapi/components/`
- Treat API contract as development input, not post-implementation artifact
- Every new endpoint must update OpenAPI contract
- Reuse common response schemas to keep success/error models consistent

### Error Handling Conventions

**Domain errors** (`domain.Error`):

```go
domain.NewError(domain.ErrorCodeForbidden)        // with []SubError
domain.NewErrorString(domain.ErrorCodeForbidden)  // with string message
```

**Handler-layer errors** (`common.Error`):

```go
common.NewBadRequestResponse([]common.SubError{{Field: "email", Message: "invalid"}})
common.NewUnAuthorized()
common.NewForbidden()
common.NewInternalServerError()
```

Global error handler in `protocol/http.go` → `middleware.ErrorHandler()` handles both types, unwraps `errors.withStack`, maps to HTTP status.

**Never `panic` in handlers.** Echo's Recover middleware catches panics, but error types must flow through `return err`.

### Env/Config Conventions

- `.env` file for local dev; OS env in CI/prod
- `ENV_FILE_PATH` env var overrides `.env` path
- All config via struct tags: `envconfig:"ENV_VAR_NAME"`
- Config loaded once in `cmd/cmds/rest.go`, passed via `Dependencies` — never call `configs.GetConfigs()` inside service/repository (only handler/protocol layer acceptable)
- Split config by concern: `appConfigs`, `PostgresConfig`, `S3Config`, `KeyConfig`, `EmailConfig`, `OAuthConfig`

---

## Bootstrap Plan

### Step 1: Project Init

```bash
mkdir new-project && cd new-project
go mod init gitlab.example.com/org/new-project
go mod vendor  # or use Go modules without vendor
```

### Step 2: Directory Scaffold

Create directories in this order:

```
configs/ infrastructure/ pkg/ internal/core/domain/ internal/core/port/
internal/core/service/ internal/core/helper/ internal/core/constant/ internal/core/enums/
internal/core/transaction/ internal/handler/base/ internal/handler/common/
internal/handler/middleware/ internal/handler/validator/ internal/repository/
protocol/ cmd/cmds/ websocket/ mocks/ assets/ docs/openapi/
```

### Step 3: Foundation Files (in order)

1. `configs/config.go` — config struct with envconfig
2. `infrastructure/postgres.go` — DB factory
3. `pkg/logx/` — logger init
4. `pkg/jwt/` — JWT middleware
5. `pkg/redis/` — Redis wrapper (if needed)
6. `internal/core/domain/error.go` — error types
7. `internal/handler/common/common.go` — response constructors
8. `internal/handler/middleware/error_handler.go` — global error handler
9. `internal/handler/validator/validator.go` — custom validator
10. `internal/core/transaction/transaction.go` — tx context helpers
11. `protocol/init.go` — application bootstrap struct
12. `cmd/cmds/root.go` + `rest.go` — cobra CLI
13. `cmd/main.go`

### Step 4: First Feature (use as template)

Pick one simple CRUD feature (e.g., `configurations`). Implement in order:

```
domain entity → port interface → repository → service → handler → route registration
```

### Step 5: Cross-cutting Concerns

- Add `websocket/websocket.go` if realtime needed
- Add `protocol/cron.go` for background jobs
- Add Prometheus metrics via `echo-contrib/prometheus`

### Migration/Bootstrap Checklist

- [ ] `.env` file with all required vars
- [ ] PostgreSQL DB created, connection string valid
- [ ] Redis available (if caching enabled)
- [ ] S3 bucket created with proper IAM policy
- [ ] `APP_PREFIX` set (used in error code formatting)
- [ ] `APP_ENV` set (`local` disables crons)
- [ ] `LOG_LEVEL` set (`info` recommended for prod)
- [ ] JWT keys configured in `KeyConfig`
- [ ] External API keys set if third-party integration needed

### Recommended Base Packages

| Purpose        | Package                                                             |
| -------------- | ------------------------------------------------------------------- |
| HTTP framework | `github.com/labstack/echo/v4`                                       |
| ORM            | `gorm.io/gorm` + `gorm.io/driver/postgres`                          |
| Config         | `github.com/kelseyhightower/envconfig` + `github.com/joho/godotenv` |
| JWT            | `github.com/golang-jwt/jwt/v5`                                      |
| Logger         | `github.com/sirupsen/logrus`                                        |
| Validator      | `github.com/go-playground/validator/v10`                            |
| UUID           | `github.com/google/uuid`                                            |
| Errors         | `github.com/pkg/errors`                                             |
| CLI            | `github.com/spf13/cobra`                                            |
| Cron           | `github.com/robfig/cron`                                            |
| WebSocket      | `github.com/gorilla/websocket`                                      |
| Excel          | `github.com/xuri/excelize/v2`                                       |
| S3             | `github.com/aws/aws-sdk-go-v2`                                      |
| Pointer utils  | `github.com/openlyinc/pointy`                                       |
| Decimal math   | `github.com/shopspring/decimal`                                     |
| Semaphore      | `golang.org/x/sync/semaphore`                                       |
| Metrics        | `github.com/labstack/echo-contrib/prometheus`                       |

### Docker/Dev Setup

```dockerfile
# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY vendor/ vendor/
COPY . .
RUN go build -mod=vendor -o server ./cmd/main.go

# Runtime stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/server .
COPY assets/ assets/
ENV_FILE_PATH=/app/.env
CMD ["./server", "serve"]
```

**Local dev**: Use `docker-compose` for Postgres + Redis. Run `go run ./cmd/main.go serve` directly.

**Env var for local**: Set `ENV_FILE_PATH=.env.local` and `APP_ENV=local` (disables crons).

---

## AI Agent Guidance

### How to Work with This Architecture

**Finding the right file:**

- Business logic → `internal/core/service/[feature]/`
- DB queries → `internal/repository/[feature]/`
- HTTP request/response → `internal/handler/[feature]/`
- Interface contracts → `internal/core/port/[feature].go`
- Shared structs → `internal/core/domain/`
- App-wide constants → `internal/core/constant/`
- Typed enums → `internal/core/enums/`

**Adding a new feature:** Always follow this sequence:

1. Add domain struct to `internal/core/domain/`
2. Add repository interface to `internal/core/port/[feature].go`
3. Add service interface to `internal/core/port/[feature].go`
4. Implement repository in `internal/repository/[feature]/`
5. Register repository in `internal/repository/repository.go`
6. Implement service in `internal/core/service/[feature]/`
7. Register service in `internal/core/service/service.go`
8. Add handler interface to `internal/handler/handler.go`
9. Implement handler in `internal/handler/[feature]/`
10. Register handler in `internal/handler/handler.go` `NewHandler()`
11. Register routes in `protocol/http.go`

**Context is the backbone.** Pass `ctx context.Context` as first arg to every service and repository method. Extract user, lang, request ID from context — never from global state.

### Conventions to Preserve

- **Aggregate structs**: `Handler`, `Service`, `Repository` are aggregate structs that hold all feature instances. Do not break this pattern by adding singleton accessors.
- **Port interfaces define contracts**: Never let a handler import a concrete repository or service struct. Always go through the port interface.
- **Dependencies struct per constructor**: Every `NewXxx()` takes a dedicated `Dependencies` struct, not variadic arguments.
- **POST for queries**: Paginated list endpoints use `POST /query` with a JSON body, not `GET` with query params.
- **Error code format**: `{PREFIX}-{CODE}` where CODE first 3 digits = HTTP status. Never return raw HTTP status without this wrapper.
- **Cron skip on local**: `if app.conf.App.ENV == "local" { return }` in `cronStart`. Preserve this guard.
- **Thai + English**: Domain structs often have `Name`/`NameEn` pairs. Preserve bilingual field patterns.
- **Context keys are typed**: Use `domain.ContextKey` typed constants for `context.WithValue`. Never use raw strings.

### Anti-Patterns to Avoid

- **No global service variables** — all deps injected, never `var globalService *MyService`
- **No business logic in handlers** — handlers bind, validate, call service, render; nothing else
- **No DB calls in handlers** — must go through service → repository chain
- **No cross-feature repository imports** — services may call other services; repositories should not call other repositories
- **No raw SQL in services** — all SQL in repository layer
- **No `configs.GetConfigs()` in service/repository** — pass config values via Dependencies at startup
- **No ignoring context cancellation** — cron jobs use `context.WithTimeout`; respect cancellation
- **No fat helpers** — helpers must be pure functions. If a helper needs a DB or HTTP call, it belongs in a service

### Documentation Update Rules

When adding/changing features:

1. Update feature port file in `internal/core/port/[feature].go`
2. Update `handler.go` with new handler interface methods
3. Update `repository/repository.go` and `service/service.go` aggregate structs
4. Register new routes in `protocol/http.go` with appropriate middleware chain
5. Update `docs/openapi/openapi.yaml` and related `components/` or `paths/`
6. If adding env vars, document them in `configs/config.go` struct tags and update `.env.example`
7. If adding cron job, add to `protocol/cron.go` following the timeout + goroutine pattern

---

## Key Technical Notes

### Semaphore Usage

`semaphore.NewWeighted(1)` acts as a mutex for operations that must not run concurrently (coupon redemption, budget deduction). Shared single semaphore instance across services that touch shared counters.

### Logging Strategy

- `logx.GetLog()` returns global logrus entry
- DB hook auto-persists all logs to `tbl_logs`; user-contextual logs to `tbl_user_logs`
- Sensitive paths (`/login`, `/refresh-token`) have body scrubbed before persisting
- Multipart file upload bodies are never stored

### WebSocket Architecture

`WebSocketManager` holds:

- `clients map[string]*Client` — userID → connection
- `groups map[string]map[string]bool` — groupName → set of userIDs

Notifications broadcast via `SendToUsers([]string, data)`. Product updates broadcast to `"product"` group via `RefreshProduct(id)`.

### External API Key Middleware

One API key:

- `X_API_KEY` — generic external service calls

Check `internal/handler/middleware/api_key_guard.go` for validation logic. Key passed in `X-Api-Key` header.
