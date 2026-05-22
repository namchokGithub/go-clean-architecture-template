# Golang Backend Scaffold Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use `superpowers:subagent-driven-development` (recommended) or `superpowers:executing-plans` to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Scaffold a production-ready Go backend API for [project_name] project following clean architecture (hexagonal/ports & adapters) defined in `docs/architechture/ARCHITECTURE.md`.

**Architecture:** Strict dependency inversion — domain core has zero knowledge of infrastructure. All external concerns (DB, HTTP, Redis) depend inward through port interfaces. Every layer uses explicit `Dependencies` structs with `NewXxx(deps)` constructors. No global state, no reflection-based DI.

**Tech Stack:** Go 1.24, Echo v4, GORM + PostgreSQL, golang-jwt/v5, logrus, go-playground/validator/v10, envconfig + godotenv, go-redis/v9, Cobra CLI

---

## File Map

```
go.mod
.env.example
.gitignore
Dockerfile
docker-compose.yml
README.md

cmd/main.go
cmd/cmds/root.go
cmd/cmds/rest.go

configs/config.go
configs/postgres.go
configs/key.go
configs/s3.go
configs/email.go
configs/oauth.go
configs/version.go

infrastructure/postgres.go

pkg/contexts/contexts.go
pkg/logx/logx.go
pkg/jwt/jwt.go
pkg/redis/redis.go
pkg/client/client.go
pkg/s3/s3.go

internal/core/domain/context.go
internal/core/domain/error.go
internal/core/domain/paginator.go
internal/core/domain/health.go

internal/core/constant/constant.go
internal/core/enums/enums.go
internal/core/transaction/transaction.go
internal/core/helper/context.go

internal/core/port/health.go

internal/core/service/service.go
internal/core/service/health/health.go
internal/core/service/health/health_test.go

internal/handler/handler.go
internal/handler/base/response.go
internal/handler/base/page.go
internal/handler/common/common.go
internal/handler/common/error.go
internal/handler/middleware/error_handler.go
internal/handler/middleware/api_key_guard.go
internal/handler/middleware/content_type.go
internal/handler/validator/validator.go
internal/handler/health/health.go

internal/repository/repository.go

protocol/init.go
protocol/http.go
protocol/cron.go
protocol/custom_middleware.go
protocol/permission_middleware.go
protocol/user_action_middleware.go

websocket/websocket.go

mocks/.gitkeep
assets/.gitkeep
docs/AGENTS.md
docs/CONTEXT.md
docs/engineering/API_GUIDELINES.md
docs/engineering/TESTING.md
docs/openapi/openapi.yaml
docs/openapi/README.md
docs/openapi/components/.gitkeep
docs/openapi/paths/.gitkeep
tests/mocks/.gitkeep
tests/integration/.gitkeep
```

---

## Task 1: Project Init + Directory Scaffold

**Files:**

- Create: `go.mod`
- Create: `.gitignore`
- Create: all directories

- [ ] **Step 1: Initialize go module**

```bash
cd /path/to/backend_[project_name]
go mod init gitlab.company.com/[project_name]/backend
```

Expected output: `go: creating new go.mod: module gitlab.company.com/[project_name]/backend`

- [ ] **Step 2: Create directory structure**

```bash
mkdir -p cmd/cmds \
  configs \
  infrastructure \
  pkg/contexts pkg/logx pkg/jwt pkg/redis pkg/client pkg/s3 \
  internal/core/domain \
  internal/core/constant \
  internal/core/enums \
  internal/core/transaction \
  internal/core/helper \
  internal/core/port \
  internal/core/service/health \
  internal/handler/base \
  internal/handler/common \
  internal/handler/middleware \
  internal/handler/validator \
  internal/handler/health \
  internal/repository \
  protocol \
  websocket \
  mocks \
  assets \
  docs/engineering \
  docs/openapi/components \
  docs/openapi/paths \
  tests/mocks \
  tests/integration \
  docs/superpowers/plans
```

- [ ] **Step 3: Create .gitignore**

```gitignore
# Binaries
server
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test
*.test
*.out
coverage.html

# Env
.env
.env.local
.env.*.local

# OS
.DS_Store
Thumbs.db

# IDE
.idea/
.vscode/
*.swp
*.swo

# Vendor (keep — project uses vendoring)
# vendor/

# Logs
*.log
```

- [ ] **Step 4: Create placeholder files so empty directories are tracked**

```bash
touch mocks/.gitkeep assets/.gitkeep docs/openapi/components/.gitkeep docs/openapi/paths/.gitkeep docs/openapi/openapi.yaml docs/openapi/README.md tests/mocks/.gitkeep tests/integration/.gitkeep
```

---

## Task 1.1: AI Workspace Foundation Docs

**Files:**

- Create: `docs/AGENTS.md`
- Create: `docs/CONTEXT.md`
- Create: `docs/engineering/API_GUIDELINES.md`
- Create: `docs/engineering/TESTING.md`

- [ ] **Step 1: Write docs/AGENTS.md**

Include concise agent rules:

- preserve `Handler -> Service -> Repository`
- no business logic in handlers
- repositories must not call other repositories
- use `context.Context` as first argument
- preserve explicit constructor injection
- use `transaction.WithTx(...)` for service-layer orchestration
- use Swagger/OpenAPI as contract source of truth
- avoid generic `utils/` package

- [ ] **Step 2: Write docs/CONTEXT.md**

Include lightweight project context:

- project overview
- current phase = scaffold / initialization
- business/domain placeholders as `TODO`
- external integration placeholders
- architecture summary

- [ ] **Step 3: Write docs/engineering/API_GUIDELINES.md**

Include concise API guidance:

- OpenAPI-first
- thin Echo handlers
- common response shape reuse
- service-layer transaction boundary
- repository owns SQL/GORM access

- [ ] **Step 4: Write docs/engineering/TESTING.md**

Include concise testing guidance:

- `mockery` generates mocks from `internal/core/port`
- generated mocks live in `tests/mocks/`
- do not mock concrete implementations
- do not edit generated mocks manually
- unit tests stay close to features
- integration tests live in `tests/integration/`

---

## Task 2: Go Dependencies

**Files:**

- Modify: `go.mod`
- Create: `go.sum`
- Create: `vendor/`

- [ ] **Step 1: Add all required dependencies**

```bash
go get github.com/labstack/echo/v4
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/kelseyhightower/envconfig
go get github.com/joho/godotenv
go get github.com/golang-jwt/jwt/v5
go get github.com/sirupsen/logrus
go get github.com/go-playground/validator/v10
go get github.com/google/uuid
go get github.com/pkg/errors
go get github.com/spf13/cobra
go get github.com/robfig/cron/v3
go get github.com/gorilla/websocket
go get github.com/openlyinc/pointy
go get github.com/shopspring/decimal
go get golang.org/x/sync
go get github.com/redis/go-redis/v9
go get github.com/aws/aws-sdk-go-v2/aws
go get github.com/aws/aws-sdk-go-v2/config
go get github.com/aws/aws-sdk-go-v2/service/s3
go get github.com/labstack/echo-contrib/prometheus
go get github.com/xuri/excelize/v2
```

- [ ] **Step 2: Tidy and vendor**

```bash
go mod tidy
go mod vendor
```

Expected: `vendor/` directory created with all dependencies. No errors.

---

## Task 3: Configs Layer

**Files:**

- Create: `configs/config.go`
- Create: `configs/postgres.go`
- Create: `configs/key.go`
- Create: `configs/s3.go`
- Create: `configs/email.go`
- Create: `configs/oauth.go`
- Create: `configs/version.go`

- [ ] **Step 1: Write configs/version.go**

```go
package configs

// These are set via -ldflags at build time.
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)
```

- [ ] **Step 2: Write configs/postgres.go**

```go
package configs

type PostgresConfig struct {
	Host         string `envconfig:"POSTGRES_HOST" default:"localhost"`
	Port         string `envconfig:"POSTGRES_PORT" default:"5432"`
	User         string `envconfig:"POSTGRES_USER" default:"postgres"`
	Password     string `envconfig:"POSTGRES_PASSWORD"`
	DBName       string `envconfig:"POSTGRES_DB" default:"database"`
	SSLMode      string `envconfig:"POSTGRES_SSL_MODE" default:"disable"`
	MaxOpenConns int    `envconfig:"POSTGRES_MAX_OPEN_CONNS" default:"10"`
	MaxIdleConns int    `envconfig:"POSTGRES_MAX_IDLE_CONNS" default:"5"`
}
```

- [ ] **Step 3: Write configs/key.go**

```go
package configs

type KeyConfig struct {
	JWTSecret string `envconfig:"JWT_SECRET"`
}
```

- [ ] **Step 4: Write configs/s3.go**

```go
package configs

type S3Config struct {
	Bucket    string `envconfig:"S3_BUCKET"`
	Region    string `envconfig:"S3_REGION" default:"ap-southeast-1"`
	AccessKey string `envconfig:"S3_ACCESS_KEY"`
	SecretKey string `envconfig:"S3_SECRET_KEY"`
	Endpoint  string `envconfig:"S3_ENDPOINT"`
}
```

- [ ] **Step 5: Write configs/email.go**

```go
package configs

type EmailConfig struct {
	Host     string `envconfig:"EMAIL_HOST"`
	Port     int    `envconfig:"EMAIL_PORT" default:"587"`
	User     string `envconfig:"EMAIL_USER"`
	Password string `envconfig:"EMAIL_PASSWORD"`
	From     string `envconfig:"EMAIL_FROM"`
}
```

- [ ] **Step 6: Write configs/oauth.go**

```go
package configs

type OAuthConfig struct {
	GoogleClientID     string `envconfig:"OAUTH_GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `envconfig:"OAUTH_GOOGLE_CLIENT_SECRET"`
	GoogleCallbackURL  string `envconfig:"OAUTH_GOOGLE_CALLBACK_URL"`
}
```

- [ ] **Step 7: Write configs/config.go**

```go
package configs

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
	Key      KeyConfig
	S3       S3Config
	Email    EmailConfig
	OAuth    OAuthConfig
}

type AppConfig struct {
	ENV      string `envconfig:"APP_ENV" default:"local"`
	Prefix   string `envconfig:"APP_PREFIX" default:"SDD"`
	Port     string `envconfig:"APP_PORT" default:"8080"`
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`
	XAPIKey  string `envconfig:"X_API_KEY"`
	Version  string `envconfig:"APP_VERSION" default:"1.0.0"`
}

// Init loads the .env file. Path is overridden by ENV_FILE_PATH env var.
func Init() {
	path := os.Getenv("ENV_FILE_PATH")
	if path == "" {
		path = ".env"
	}
	_ = godotenv.Load(path)
}

// GetConfigs loads and returns all config. Call once at startup; inject via Dependencies.
func GetConfigs() Config {
	Init()
	var cfg Config
	mustProcess(&cfg.App)
	mustProcess(&cfg.Postgres)
	mustProcess(&cfg.Key)
	mustProcess(&cfg.S3)
	mustProcess(&cfg.Email)
	mustProcess(&cfg.OAuth)
	return cfg
}

func mustProcess(spec interface{}) {
	if err := envconfig.Process("", spec); err != nil {
		panic(err)
	}
}
```

- [ ] **Step 8: Verify compilation**

```bash
go build ./configs/...
```

Expected: no output, exit 0.

---

## Task 4: Infrastructure

**Files:**

- Create: `infrastructure/postgres.go`

- [ ] **Step 1: Write infrastructure/postgres.go**

```go
package infrastructure

import (
	"fmt"

	"gitlab.company.com/[project_name]/backend/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewPostgresDB opens a GORM DB connection with the given config.
func NewPostgresDB(cfg configs.PostgresConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Bangkok",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	return db, nil
}
```

- [ ] **Step 2: Verify compilation**

```bash
go build ./infrastructure/...
```

---

## Task 5: pkg Layer

**Files:**

- Create: `pkg/contexts/contexts.go`
- Create: `pkg/logx/logx.go`
- Create: `pkg/jwt/jwt.go`
- Create: `pkg/redis/redis.go`
- Create: `pkg/client/client.go`
- Create: `pkg/s3/s3.go`

- [ ] **Step 1: Write pkg/contexts/contexts.go**

```go
package contexts

import "context"

// Key is the typed context key to prevent collisions with other packages.
type Key string

// Set stores a value under a typed key.
func Set(ctx context.Context, key Key, val any) context.Context {
	return context.WithValue(ctx, key, val)
}

// Get retrieves a value and casts it to T. Returns zero value + false if missing or wrong type.
func Get[T any](ctx context.Context, key Key) (T, bool) {
	v, ok := ctx.Value(key).(T)
	return v, ok
}
```

- [ ] **Step 2: Write pkg/logx/logx.go**

```go
package logx

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	once sync.Once
	log  *logrus.Entry
)

// Init configures the global logger. Safe to call multiple times; only first call takes effect.
func Init(level string) {
	once.Do(func() {
		l := logrus.New()
		l.SetOutput(os.Stdout)
		l.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z",
		})
		lvl, err := logrus.ParseLevel(level)
		if err != nil {
			lvl = logrus.InfoLevel
		}
		l.SetLevel(lvl)
		log = logrus.NewEntry(l)
	})
}

// GetLog returns the global logrus entry, initializing with "info" level if not yet initialized.
func GetLog() *logrus.Entry {
	if log == nil {
		Init("info")
	}
	return log
}
```

- [ ] **Step 3: Write pkg/jwt/jwt.go**

```go
package jwt

import (
	"context"
	"fmt"
	"strings"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gitlab.company.com/[project_name]/backend/internal/core/domain"
	"gitlab.company.com/[project_name]/backend/internal/handler/common"
)

// Claims is the JWT payload stored in the token.
type Claims struct {
	UserID   int64  `json:"user_id"`
	UserType string `json:"user_type"`
	RoleID   int64  `json:"role_id"`
	Email    string `json:"email"`
	gojwt.RegisteredClaims
}

// AuthMiddleware validates the Bearer JWT and injects claims into the request context.
func AuthMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				return common.NewUnAuthorized()
			}
			tokenStr := strings.TrimPrefix(auth, "Bearer ")

			claims := &Claims{}
			token, err := gojwt.ParseWithClaims(tokenStr, claims, func(t *gojwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*gojwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				return common.NewUnAuthorized()
			}

			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, domain.USER_ID, claims.UserID)
			ctx = context.WithValue(ctx, domain.USER_TYPE, claims.UserType)
			ctx = context.WithValue(ctx, domain.ROLE_ID, claims.RoleID)
			ctx = context.WithValue(ctx, domain.EMAIL, claims.Email)
			ctx = context.WithValue(ctx, domain.TOKEN, tokenStr)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
```

- [ ] **Step 4: Write pkg/redis/redis.go**

```go
package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client wraps go-redis for typed access.
type Client struct {
	rdb *redis.Client
}

// Config holds Redis connection settings.
type Config struct {
	Addr     string
	Password string
	DB       int
}

// New creates a Redis client and pings to verify connectivity.
func New(cfg Config) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &Client{rdb: rdb}, nil
}

func (c *Client) Set(ctx context.Context, key string, val any, ttl time.Duration) error {
	return c.rdb.Set(ctx, key, val, ttl).Err()
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.rdb.Get(ctx, key).Result()
}

func (c *Client) Del(ctx context.Context, keys ...string) error {
	return c.rdb.Del(ctx, keys...).Err()
}

func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.rdb.Exists(ctx, key).Result()
	return n > 0, err
}
```

- [ ] **Step 5: Write pkg/client/client.go**

```go
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is a reusable HTTP client wrapper.
type Client struct {
	base    string
	headers map[string]string
	http    *http.Client
}

// New creates a Client pointing at baseURL with a default 30s timeout.
func New(baseURL string) *Client {
	return &Client{
		base:    baseURL,
		headers: make(map[string]string),
		http:    &http.Client{Timeout: 30 * time.Second},
	}
}

// WithHeader adds a persistent header to all requests.
func (c *Client) WithHeader(key, value string) *Client {
	c.headers[key] = value
	return c
}

// Post sends a JSON POST request and decodes the response into out.
func (c *Client) Post(ctx context.Context, path string, body any, out any) error {
	return c.do(ctx, http.MethodPost, path, body, out)
}

// Get sends a GET request and decodes the response into out.
func (c *Client) Get(ctx context.Context, path string, out any) error {
	return c.do(ctx, http.MethodGet, path, nil, out)
}

func (c *Client) do(ctx context.Context, method, path string, body any, out any) error {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.base+path, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("upstream returned %d", resp.StatusCode)
	}
	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}
```

- [ ] **Step 6: Write pkg/s3/s3.go**

```go
package s3

import (
	"context"
	"io"
)

// Storage defines the S3 operations used by the application.
type Storage interface {
	Upload(ctx context.Context, key string, r io.Reader, contentType string) (string, error)
	Delete(ctx context.Context, key string) error
	GetURL(key string) string
}
```

- [ ] **Step 7: Verify pkg compilation**

```bash
go build ./pkg/...
```

Note: `pkg/jwt` imports `internal/handler/common` and `internal/core/domain` — those packages are written in later tasks. If this step fails with import errors, defer this verification to after Task 8.

---

## Task 6: Core Domain Types

**Files:**

- Create: `internal/core/domain/context.go`
- Create: `internal/core/domain/error.go`
- Create: `internal/core/domain/paginator.go`
- Create: `internal/core/domain/health.go`

- [ ] **Step 1: Write internal/core/domain/context.go**

```go
package domain

// ContextKey is the typed key for context values to prevent collisions.
type ContextKey string

const (
	Lang           ContextKey = "Lang"
	RequestID      ContextKey = "RequestID"
	OS             ContextKey = "OS"
	Browser        ContextKey = "Browser"
	USER_ID        ContextKey = "USER_ID"
	USER_TYPE      ContextKey = "USER_TYPE"
	ROLE_ID        ContextKey = "ROLE_ID"
	TOKEN          ContextKey = "TOKEN"
	EMAIL          ContextKey = "EMAIL"
	USER           ContextKey = "USER"
	ROLE_PERMISSION ContextKey = "ROLE_PERMISSION"
	MENU_ID        ContextKey = "MENU_ID"
)
```

- [ ] **Step 2: Write internal/core/domain/error.go**

```go
package domain

import "fmt"

// ErrorCode maps to HTTP status via its first 3 digits (e.g., "4001" → HTTP 400).
type ErrorCode string

const (
	ErrorCodeBadRequest   ErrorCode = "4001"
	ErrorCodeUnauthorized ErrorCode = "4011"
	ErrorCodeForbidden    ErrorCode = "4031"
	ErrorCodeNotFound     ErrorCode = "4041"
	ErrorCodeConflict     ErrorCode = "4091"
	ErrorCodeInternal     ErrorCode = "5001"
)

// SubError represents a single field-level validation error.
type SubError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error is the domain-layer error type. Handlers map it to HTTP responses.
type Error struct {
	Code    ErrorCode
	Message string
	Errors  []SubError
}

func (e *Error) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewError creates a domain error with field-level sub-errors.
func NewError(code ErrorCode, errs []SubError) *Error {
	return &Error{Code: code, Errors: errs}
}

// NewErrorString creates a domain error with a plain message.
func NewErrorString(code ErrorCode, message string) *Error {
	return &Error{Code: code, Message: message}
}
```

- [ ] **Step 3: Write internal/core/domain/paginator.go**

```go
package domain

// Paginator wraps a paginated result set with metadata.
type Paginator struct {
	Data      any   `json:"data"`
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

// NewPaginator constructs a Paginator, calculating total pages from total + pageSize.
func NewPaginator(data any, page, pageSize int, total int64) Paginator {
	totalPage := int(total) / pageSize
	if total%int64(pageSize) != 0 {
		totalPage++
	}
	if pageSize == 0 {
		totalPage = 1
	}
	return Paginator{
		Data:      data,
		Page:      page,
		PageSize:  pageSize,
		Total:     total,
		TotalPage: totalPage,
	}
}
```

- [ ] **Step 4: Write internal/core/domain/health.go**

```go
package domain

import "time"

// HealthStatus is the response model for the health check endpoint.
type HealthStatus struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}
```

- [ ] **Step 5: Verify compilation**

```bash
go build ./internal/core/domain/...
```

---

## Task 7: Core Constants, Enums, Transaction, Helper

**Files:**

- Create: `internal/core/constant/constant.go`
- Create: `internal/core/enums/enums.go`
- Create: `internal/core/transaction/transaction.go`
- Create: `internal/core/helper/context.go`

- [ ] **Step 1: Write internal/core/constant/constant.go**

```go
package constant

const (
	SortAsc  = "asc"
	SortDesc = "desc"

	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)
```

- [ ] **Step 2: Write internal/core/enums/enums.go**

```go
package enums

// Enum is the base type for typed enumerations with bilingual labels.
type Enum struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`    // Thai
	NameEn string `json:"name_en"` // English
}

// UserType enumerates valid user type values.
var UserType = struct {
	Admin  Enum
	Driver Enum
}{
	Admin:  Enum{ID: 1, Name: "แอดมิน", NameEn: "Admin"},
	Driver: Enum{ID: 2, Name: "คนขับ", NameEn: "Driver"},
}
```

- [ ] **Step 3: Write internal/core/transaction/transaction.go**

```go
package transaction

import (
	"context"

	"gorm.io/gorm"
)

type txKey struct{}

// NewContext embeds a GORM transaction into context.
func NewContext(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// FromContext extracts the transaction from context. Returns (nil, false) if none.
func FromContext(ctx context.Context) (*gorm.DB, bool) {
	tx, ok := ctx.Value(txKey{}).(*gorm.DB)
	return tx, ok && tx != nil
}

// WithTx centralizes tx lifecycle, panic safety, and rollback handling.
func WithTx(ctx context.Context, db *gorm.DB, fn func(ctx context.Context) error) (err error) {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		}
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = fn(NewContext(ctx, tx)); err != nil {
		return err
	}

	return tx.Commit().Error
}
```

- [ ] **Step 4: Write internal/core/helper/context.go**

```go
package helper

import (
	"context"

	"gitlab.company.com/[project_name]/backend/internal/core/domain"
)

// GetUserID extracts USER_ID from context. Returns 0 if missing.
func GetUserID(ctx context.Context) int64 {
	v, _ := ctx.Value(domain.USER_ID).(int64)
	return v
}

// GetUserType extracts USER_TYPE from context.
func GetUserType(ctx context.Context) string {
	v, _ := ctx.Value(domain.USER_TYPE).(string)
	return v
}

// GetRoleID extracts ROLE_ID from context.
func GetRoleID(ctx context.Context) int64 {
	v, _ := ctx.Value(domain.ROLE_ID).(int64)
	return v
}

// GetLang extracts the Lang from context. Returns "th" as default.
func GetLang(ctx context.Context) string {
	v, _ := ctx.Value(domain.Lang).(string)
	if v == "" {
		return "th"
	}
	return v
}

// GetToken extracts the raw JWT token string from context.
func GetToken(ctx context.Context) string {
	v, _ := ctx.Value(domain.TOKEN).(string)
	return v
}
```

- [ ] **Step 5: Verify compilation**

```bash
go build ./internal/core/...
```

---

## Task 8: Core Port Interfaces

**Files:**

- Create: `internal/core/port/health.go`

- [ ] **Step 1: Write internal/core/port/health.go**

```go
package port

import (
	"context"

	"gitlab.company.com/[project_name]/backend/internal/core/domain"
)

// HealthService defines the contract for application health checks.
type HealthService interface {
	GetStatus(ctx context.Context) (domain.HealthStatus, error)
}
type HealthRepository interface {
	GetStatus(ctx context.Context) (domain.HealthStatus, error)
}
```

- [ ] **Step 2: Verify compilation**

```bash
go build ./internal/core/port/...
```

---

## Task 9: Handler Framework — Base + Common

**Files:**

- Create: `internal/handler/base/response.go`
- Create: `internal/handler/base/page.go`
- Create: `internal/handler/common/common.go`
- Create: `internal/handler/common/error.go`

- [ ] **Step 1: Write internal/handler/base/response.go**

```go
package base

// DefaultResponse is the standard success response envelope.
type DefaultResponse struct {
	IsError bool   `json:"is_error"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
```

- [ ] **Step 2: Write internal/handler/base/page.go**

```go
package base

// BasePageRequest is the standard body for all paginated POST /query endpoints.
type BasePageRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	SortName string `json:"sort_name"`
	SortBy   string `json:"sort_by"` // "asc" | "desc"
	Search   string `json:"search"`
	IsAll    bool   `json:"is_all"` // bypass pagination when true
}

// Normalize applies defaults for missing pagination fields.
func (r *BasePageRequest) Normalize() {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.PageSize <= 0 || r.PageSize > 100 {
		r.PageSize = 20
	}
	if r.SortBy != "asc" && r.SortBy != "desc" {
		r.SortBy = "desc"
	}
}
```

- [ ] **Step 3: Write internal/handler/common/error.go**

```go
package common

// ErrorResponse is the standard error response envelope.
type ErrorResponse struct {
	IsError   bool   `json:"is_error"`
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	ErrorType string `json:"error_type"`
	Errors    any    `json:"errors"`
}

// Error is the handler-layer error type returned from handler functions.
// Echo's global error handler maps it to an HTTP response.
type Error struct {
	HTTPStatus int
	Code       string
	Message    string
	ErrorType  string
	Errors     any
}

func (e *Error) Error() string { return e.Message }
```

- [ ] **Step 4: Write internal/handler/common/common.go**

```go
package common

import (
	"net/http"

	"gitlab.company.com/[project_name]/backend/internal/handler/base"
)

var appPrefix = "APP"

// SetAppPrefix sets the prefix used in all response codes (e.g., "SDD").
// Call once during protocol initialization.
func SetAppPrefix(prefix string) {
	appPrefix = prefix
}

// GetAppPrefix returns the current app prefix.
func GetAppPrefix() string {
	return appPrefix
}

func code(suffix string) string {
	return appPrefix + "-" + suffix
}

// Success responses

func NewSuccessResponse(data any) base.DefaultResponse {
	return base.DefaultResponse{IsError: false, Code: code("200"), Message: "Success", Data: data}
}

func NewCreatedResponse(data any) base.DefaultResponse {
	return base.DefaultResponse{IsError: false, Code: code("200"), Message: "Created", Data: data}
}

// Error responses — each returns *Error which Echo's error handler renders.

func NewBadRequest(errs any) *Error {
	return &Error{HTTPStatus: http.StatusBadRequest, Code: code("4001"), Message: "Bad Request", ErrorType: "BAD_REQUEST", Errors: errs}
}

func NewUnAuthorized() *Error {
	return &Error{HTTPStatus: http.StatusUnauthorized, Code: code("4011"), Message: "Unauthorized", ErrorType: "UNAUTHORIZED"}
}

func NewForbidden() *Error {
	return &Error{HTTPStatus: http.StatusForbidden, Code: code("4031"), Message: "Forbidden", ErrorType: "FORBIDDEN"}
}

func NewNotFound() *Error {
	return &Error{HTTPStatus: http.StatusNotFound, Code: code("4041"), Message: "Not Found", ErrorType: "NOT_FOUND"}
}

func NewConflict() *Error {
	return &Error{HTTPStatus: http.StatusConflict, Code: code("4091"), Message: "Conflict", ErrorType: "CONFLICT"}
}

func NewInternalServerError() *Error {
	return &Error{HTTPStatus: http.StatusInternalServerError, Code: code("5001"), Message: "Internal Server Error", ErrorType: "INTERNAL_SERVER_ERROR"}
}

// SubError is used in validation error responses.
type SubError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
```

- [ ] **Step 5: Verify compilation**

```bash
go build ./internal/handler/base/... ./internal/handler/common/...
```

---

## Task 10: Handler Framework — Middleware + Validator

**Files:**

- Create: `internal/handler/middleware/error_handler.go`
- Create: `internal/handler/middleware/api_key_guard.go`
- Create: `internal/handler/middleware/content_type.go`
- Create: `internal/handler/validator/validator.go`

- [ ] **Step 1: Write internal/handler/middleware/error_handler.go**

```go
package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.company.com/[project_name]/backend/internal/core/domain"
	"gitlab.company.com/[project_name]/backend/internal/handler/common"
)

// ErrorHandler returns an Echo HTTPErrorHandler that maps domain and handler errors to JSON.
func ErrorHandler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		reqID := c.Response().Header().Get(echo.HeaderXRequestID)
		var httpStatus int
		var resp common.ErrorResponse

		switch e := err.(type) {
		case *common.Error:
			httpStatus = e.HTTPStatus
			resp = common.ErrorResponse{
				IsError:   true,
				Code:      e.Code,
				Message:   e.Message,
				RequestID: reqID,
				ErrorType: e.ErrorType,
				Errors:    e.Errors,
			}
		case *domain.Error:
			httpStatus = domainHTTPStatus(e.Code)
			resp = common.ErrorResponse{
				IsError:   true,
				Code:      common.GetAppPrefix() + "-" + string(e.Code),
				Message:   e.Message,
				RequestID: reqID,
				ErrorType: "DOMAIN_ERROR",
				Errors:    e.Errors,
			}
		case *echo.HTTPError:
			httpStatus = e.Code
			resp = common.ErrorResponse{
				IsError:   true,
				Code:      common.GetAppPrefix() + "-" + echoCode(e.Code),
				Message:   http.StatusText(e.Code),
				RequestID: reqID,
			}
		default:
			httpStatus = http.StatusInternalServerError
			resp = common.ErrorResponse{
				IsError:   true,
				Code:      common.GetAppPrefix() + "-5001",
				Message:   "Internal Server Error",
				RequestID: reqID,
			}
		}

		_ = c.JSON(httpStatus, resp)
	}
}

func domainHTTPStatus(code domain.ErrorCode) int {
	switch code {
	case domain.ErrorCodeBadRequest:
		return http.StatusBadRequest
	case domain.ErrorCodeUnauthorized:
		return http.StatusUnauthorized
	case domain.ErrorCodeForbidden:
		return http.StatusForbidden
	case domain.ErrorCodeNotFound:
		return http.StatusNotFound
	case domain.ErrorCodeConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func echoCode(status int) string {
	switch status {
	case 400:
		return "4001"
	case 401:
		return "4011"
	case 403:
		return "4031"
	case 404:
		return "4041"
	case 409:
		return "4091"
	default:
		return "5001"
	}
}
```

- [ ] **Step 2: Write internal/handler/middleware/api_key_guard.go**

```go
package middleware

import (
	"github.com/labstack/echo/v4"
	"gitlab.company.com/[project_name]/backend/internal/handler/common"
)

// APIKeyGuard rejects requests that do not supply the correct X-Api-Key header.
func APIKeyGuard(expectedKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Header.Get("X-Api-Key") != expectedKey {
				return common.NewUnAuthorized()
			}
			return next(c)
		}
	}
}
```

- [ ] **Step 3: Write internal/handler/middleware/content_type.go**

```go
package middleware

import (
	"github.com/labstack/echo/v4"
	"strings"
)

// EnforceJSONContentType rejects non-JSON POST/PUT/PATCH requests.
func EnforceJSONContentType() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			method := c.Request().Method
			if method == "POST" || method == "PUT" || method == "PATCH" {
				ct := c.Request().Header.Get("Content-Type")
				if !strings.HasPrefix(ct, "application/json") && !strings.HasPrefix(ct, "multipart/form-data") {
					c.Request().Header.Set("Content-Type", "application/json")
				}
			}
			return next(c)
		}
	}
}
```

- [ ] **Step 4: Write internal/handler/validator/validator.go**

```go
package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gitlab.company.com/[project_name]/backend/internal/handler/common"
)

// CustomValidator adapts go-playground/validator to Echo's Validator interface.
type CustomValidator struct {
	v *validator.Validate
}

// New creates a CustomValidator and registers any custom validation rules.
func New() *CustomValidator {
	v := validator.New()
	return &CustomValidator{v: v}
}

// Validate implements echo.Validator. Returns common.Error on failure so the error
// handler maps it to HTTP 400 with field-level sub-errors.
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.v.Struct(i); err != nil {
		var errs []common.SubError
		for _, e := range err.(validator.ValidationErrors) {
			errs = append(errs, common.SubError{
				Field:   e.Field(),
				Message: e.Tag(),
			})
		}
		return common.NewBadRequest(errs)
	}
	return nil
}

// Register adds this validator to the Echo instance.
func Register(e *echo.Echo) {
	e.Validator = New()
}
```

- [ ] **Step 5: Verify compilation**

```bash
go build ./internal/handler/...
```

---

## Task 11: Aggregate Structs

**Files:**

- Create: `internal/repository/repository.go`
- Create: `internal/core/service/service.go`
- Create: `internal/handler/handler.go`

- [ ] **Step 1: Write internal/repository/repository.go**

```go
package repository

import (
	"gorm.io/gorm"
)

// Dependencies holds all infrastructure deps for repository constructors.
type Dependencies struct {
	DB *gorm.DB
}

// Repository is the aggregate that holds all feature repositories.
// Add a field for each feature repository as they are created.
type Repository struct {
	// Example: Health *health.Repository
}

// New creates the Repository aggregate. Register feature repositories here.
func New(d Dependencies) *Repository {
	return &Repository{}
}
```

- [ ] **Step 2: Write internal/core/service/service.go**

```go
package service

import (
	"gitlab.company.com/[project_name]/backend/configs"
	"gitlab.company.com/[project_name]/backend/internal/core/port"
	"gitlab.company.com/[project_name]/backend/internal/core/service/health"
)

// Dependencies holds all service-layer dependencies.
type Dependencies struct {
	Config configs.Config
	Repo   interface{} // typed repos injected by protocol/init.go
}

// Service is the aggregate that holds all feature services.
type Service struct {
	Health port.HealthService
}

// New creates the Service aggregate. Register feature services here.
func New(d Dependencies) *Service {
	return &Service{
		Health: health.New(health.Dependencies{
			Version: d.Config.App.Version,
		}),
	}
}
```

- [ ] **Step 3: Write internal/handler/handler.go**

```go
package handler

import (
	"github.com/labstack/echo/v4"
	"gitlab.company.com/[project_name]/backend/internal/core/port"
	internalhealth "gitlab.company.com/[project_name]/backend/internal/handler/health"
)

// Handler is the aggregate that holds all feature handlers.
type Handler struct {
	Health *internalhealth.Handler
}

// Dependencies holds port interfaces required by handlers.
type Dependencies struct {
	HealthService port.HealthService
}

// New creates the Handler aggregate.
func New(d Dependencies) *Handler {
	return &Handler{
		Health: internalhealth.New(internalhealth.Dependencies{
			Service: d.HealthService,
		}),
	}
}

// RegisterRoutes wires all handler routes onto the Echo instance.
// Pass the auth-gated group as authGroup; public routes go on e directly.
func (h *Handler) RegisterRoutes(e *echo.Echo, authGroup *echo.Group) {
	h.Health.RegisterRoutes(e)
}
```

- [ ] **Step 4: Verify compilation (will fail until health feature is implemented in Task 12)**

Wait for Task 12 before running `go build ./...`.

---

## Task 12: Health Feature — Service, Handler, Test

**Files:**

- Create: `internal/core/service/health/health.go`
- Create: `internal/core/service/health/health_test.go`
- Create: `internal/handler/health/health.go`

- [ ] **Step 1: Write internal/core/service/health/health.go**

```go
package health

import (
	"context"
	"time"

	"gitlab.company.com/[project_name]/backend/internal/core/domain"
)

// Dependencies holds construction-time config for the health service.
type Dependencies struct {
	Version string
}

// service implements port.HealthService.
type service struct {
	version string
}

// New creates a health service.
func New(d Dependencies) *service {
	return &service{version: d.Version}
}

func (s *service) GetStatus(_ context.Context) (domain.HealthStatus, error) {
	return domain.HealthStatus{
		Status:    "ok",
		Timestamp: time.Now().UTC(),
		Version:   s.version,
	}, nil
}
```

- [ ] **Step 2: Write failing test**

```go
// internal/core/service/health/health_test.go
package health_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.company.com/[project_name]/backend/internal/core/service/health"
)

func TestGetStatus_ReturnsOKWithVersion(t *testing.T) {
	svc := health.New(health.Dependencies{Version: "1.2.3"})

	status, err := svc.GetStatus(context.Background())

	require.NoError(t, err)
	assert.Equal(t, "ok", status.Status)
	assert.Equal(t, "1.2.3", status.Version)
	assert.False(t, status.Timestamp.IsZero())
}
```

- [ ] **Step 3: Add testify dependency**

```bash
go get github.com/stretchr/testify
go mod tidy
```

- [ ] **Step 4: Run test — expect FAIL (service not yet confirmed compilable)**

```bash
go test ./internal/core/service/health/... -v
```

Expected: PASS (service already written above). If FAIL, check compilation errors.

- [ ] **Step 5: Write internal/handler/health/health.go**

```go
package health

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.company.com/[project_name]/backend/internal/core/port"
	"gitlab.company.com/[project_name]/backend/internal/handler/common"
)

// Handler handles HTTP requests for health checks.
type Handler struct {
	svc port.HealthService
}

// Dependencies holds the port interfaces needed by this handler.
type Dependencies struct {
	Service port.HealthService
}

// New creates a health Handler.
func New(d Dependencies) *Handler {
	return &Handler{svc: d.Service}
}

// RegisterRoutes registers public routes on the Echo root instance.
func (h *Handler) RegisterRoutes(e *echo.Echo) {
	e.GET("/health", h.GetStatus)
}

// GetStatus handles GET /health.
func (h *Handler) GetStatus(c echo.Context) error {
	status, err := h.svc.GetStatus(c.Request().Context())
	if err != nil {
		return common.NewInternalServerError()
	}
	return c.JSON(http.StatusOK, common.NewSuccessResponse(status))
}
```

- [ ] **Step 6: Verify full compilation**

```bash
go build ./...
```

Expected: exit 0, no errors.

- [ ] **Step 7: Run all tests**

```bash
go test ./... -v
```

---

## Task 13: Protocol Layer

**Files:**

- Create: `protocol/init.go`
- Create: `protocol/http.go`
- Create: `protocol/cron.go`
- Create: `protocol/custom_middleware.go`
- Create: `protocol/permission_middleware.go`
- Create: `protocol/user_action_middleware.go`

- [ ] **Step 1: Write protocol/init.go**

```go
package protocol

import (
	"gitlab.company.com/[project_name]/backend/configs"
	"gitlab.company.com/[project_name]/backend/infrastructure"
	"gitlab.company.com/[project_name]/backend/internal/core/service"
	"gitlab.company.com/[project_name]/backend/internal/handler"
	"gitlab.company.com/[project_name]/backend/internal/handler/common"
	"gitlab.company.com/[project_name]/backend/internal/repository"
	"gitlab.company.com/[project_name]/backend/pkg/logx"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// App is the root application struct. All deps flow through here.
type App struct {
	conf    configs.Config
	echo    *echo.Echo
	db      *gorm.DB
	handler *handler.Handler
	service *service.Service
}

// New bootstraps the entire application from config.
func New(conf configs.Config) (*App, error) {
	logx.Init(conf.App.LogLevel)
	log := logx.GetLog()

	db, err := infrastructure.NewPostgresDB(conf.Postgres)
	if err != nil {
		log.WithError(err).Error("failed to connect to postgres")
		return nil, err
	}
	log.Info("postgres connected")

	repo := repository.New(repository.Dependencies{DB: db})
	_ = repo // suppress unused warning until repos are registered

	svc := service.New(service.Dependencies{Config: conf})

	common.SetAppPrefix(conf.App.Prefix)

	h := handler.New(handler.Dependencies{
		HealthService: svc.Health,
	})

	app := &App{
		conf:    conf,
		db:      db,
		handler: h,
		service: svc,
	}

	app.echo = app.newEcho()
	return app, nil
}

// Start begins listening for HTTP requests.
func (a *App) Start() error {
	return a.echo.Start(":" + a.conf.App.Port)
}
```

- [ ] **Step 2: Write protocol/http.go**

```go
package protocol

import (
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"gitlab.company.com/[project_name]/backend/internal/handler/middleware"
	"gitlab.company.com/[project_name]/backend/internal/handler/validator"
	pkgjwt "gitlab.company.com/[project_name]/backend/pkg/jwt"
)

func (a *App) newEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HTTPErrorHandler = middleware.ErrorHandler()
	validator.Register(e)

	// Global middleware
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.RequestID())
	e.Use(echomiddleware.CORS())
	e.Use(a.globalContextMiddleware())

	// Auth-gated group
	authGroup := e.Group("")
	authGroup.Use(pkgjwt.AuthMiddleware(a.conf.Key.JWTSecret))
	authGroup.Use(a.afterAuthMiddleware())
	authGroup.Use(a.permissionMiddleware())
	authGroup.Use(a.userActionMiddleware())

	// Register routes (public + auth-gated)
	a.handler.RegisterRoutes(e, authGroup)

	return e
}
```

- [ ] **Step 3: Write protocol/cron.go**

```go
package protocol

import (
	"github.com/robfig/cron/v3"
	"gitlab.company.com/[project_name]/backend/pkg/logx"
)

// StartCron registers and starts background cron jobs.
// No-ops when APP_ENV=local to prevent accidental job execution during development.
func (a *App) StartCron() {
	if a.conf.App.ENV == "local" {
		return
	}
	c := cron.New()
	log := logx.GetLog()

	// Register jobs here:
	// c.AddFunc("@every 1m", func() { a.service.SomeJob(...) })

	c.Start()
	log.Info("cron started")
}
```

- [ ] **Step 4: Write protocol/custom_middleware.go**

```go
package protocol

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"gitlab.company.com/[project_name]/backend/internal/core/domain"
	"gitlab.company.com/[project_name]/backend/internal/handler/common"
	"gitlab.company.com/[project_name]/backend/pkg/logx"
)

// globalContextMiddleware injects Lang, OS, Browser into the request context.
func (a *App) globalContextMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			lang := c.Request().Header.Get("Accept-Language")
			if lang == "" {
				lang = "th"
			}
			ctx = context.WithValue(ctx, domain.Lang, lang)
			ctx = context.WithValue(ctx, domain.OS, c.Request().Header.Get("X-OS"))
			ctx = context.WithValue(ctx, domain.Browser, c.Request().Header.Get("User-Agent"))
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

// afterAuthMiddleware loads the authenticated user and role permissions.
// Checks Redis cache first; falls back to DB on cache miss.
// For routes that need full user context (all auth-gated routes).
func (a *App) afterAuthMiddleware() echo.MiddlewareFunc {
	log := logx.GetLog()
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			token, _ := ctx.Value(domain.TOKEN).(string)
			if token == "" {
				return common.NewUnAuthorized()
			}

			// TODO: load user from Redis cache (`{token}_user`)
			// TODO: fall back to DB and re-cache on miss
			// TODO: load role permissions from Redis (`{token}_rolePermission`)
			// TODO: inject domain.USER and domain.ROLE_PERMISSION into ctx

			log.WithField("token_prefix", safePrefix(token)).Debug("afterAuth")
			return next(c)
		}
	}
}

func safePrefix(s string) string {
	if len(s) > 8 {
		return s[:8] + "..."
	}
	return s
}

// keepAliveTimeout is used to limit afterAuth DB fallback queries.
const keepAliveTimeout = 5 * time.Second

var _ = keepAliveTimeout // prevent unused error until used
```

- [ ] **Step 5: Write protocol/permission_middleware.go**

```go
package protocol

import "github.com/labstack/echo/v4"

// permissionMiddleware checks ROLE_PERMISSION[menuID] for required access.
// Currently a pass-through; populate when RBAC is fully wired.
func (a *App) permissionMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// TODO: extract MENU_ID from route, check ROLE_PERMISSION from context
			return next(c)
		}
	}
}
```

- [ ] **Step 6: Write protocol/user_action_middleware.go**

```go
package protocol

import (
	"github.com/labstack/echo/v4"
	"gitlab.company.com/[project_name]/backend/internal/core/domain"
	"gitlab.company.com/[project_name]/backend/pkg/logx"
)

// userActionMiddleware logs every authenticated request to tbl_user_logs.
func (a *App) userActionMiddleware() echo.MiddlewareFunc {
	log := logx.GetLog()
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			userID, _ := ctx.Value(domain.USER_ID).(int64)
			log.WithFields(map[string]interface{}{
				"user_id": userID,
				"method":  c.Request().Method,
				"path":    c.Path(),
			}).Info("user_action")
			// TODO: persist to tbl_user_logs via logx DB hook
			return next(c)
		}
	}
}
```

- [ ] **Step 7: Verify protocol compilation**

```bash
go build ./protocol/...
```

---

## Task 14: WebSocket Stub

**Files:**

- Create: `websocket/websocket.go`

- [ ] **Step 1: Write websocket/websocket.go**

```go
package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Client holds a WebSocket connection for a single user.
type Client struct {
	UserID string
	Conn   *websocket.Conn
}

// Manager tracks connected clients and named groups.
type Manager struct {
	mu      sync.RWMutex
	clients map[string]*Client            // userID → client
	groups  map[string]map[string]bool    // groupName → set of userIDs
}

// NewManager creates an empty WebSocket manager.
func NewManager() *Manager {
	return &Manager{
		clients: make(map[string]*Client),
		groups:  make(map[string]map[string]bool),
	}
}

// Register adds a client to the manager.
func (m *Manager) Register(c *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clients[c.UserID] = c
}

// Unregister removes a client and closes the connection.
func (m *Manager) Unregister(userID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if c, ok := m.clients[userID]; ok {
		c.Conn.Close()
		delete(m.clients, userID)
	}
}

// JoinGroup adds a user to a named group.
func (m *Manager) JoinGroup(groupName, userID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.groups[groupName] == nil {
		m.groups[groupName] = make(map[string]bool)
	}
	m.groups[groupName][userID] = true
}

// SendToUsers broadcasts data to a list of users.
func (m *Manager) SendToUsers(userIDs []string, data []byte) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, id := range userIDs {
		if c, ok := m.clients[id]; ok {
			_ = c.Conn.WriteMessage(websocket.TextMessage, data)
		}
	}
}

// RefreshProduct broadcasts a product refresh event to the "product" group.
func (m *Manager) RefreshProduct(productID int64) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	group, ok := m.groups["product"]
	if !ok {
		return
	}
	for userID := range group {
		if c, ok2 := m.clients[userID]; ok2 {
			msg := []byte(`{"event":"product_refresh"}`)
			_ = c.Conn.WriteMessage(websocket.TextMessage, msg)
			_ = productID // used in message payload in real impl
		}
	}
}
```

---

## Task 15: cmd — Cobra CLI

**Files:**

- Create: `cmd/cmds/root.go`
- Create: `cmd/cmds/rest.go`
- Create: `cmd/main.go`

- [ ] **Step 1: Write cmd/cmds/root.go**

```go
package cmds

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "[project_name] backend API server",
}

// Execute runs the root cobra command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(restCmd)
}
```

- [ ] **Step 2: Write cmd/cmds/rest.go**

```go
package cmds

import (
	"github.com/spf13/cobra"
	"gitlab.company.com/[project_name]/backend/configs"
	"gitlab.company.com/[project_name]/backend/pkg/logx"
	"gitlab.company.com/[project_name]/backend/protocol"
)

var restCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP API server",
	RunE:  runRest,
}

func runRest(_ *cobra.Command, _ []string) error {
	conf := configs.GetConfigs()
	logx.Init(conf.App.LogLevel)
	log := logx.GetLog()

	app, err := protocol.New(conf)
	if err != nil {
		log.WithError(err).Fatal("failed to initialize application")
		return err
	}

	app.StartCron()
	log.WithField("port", conf.App.Port).Info("starting server")
	return app.Start()
}
```

- [ ] **Step 3: Write cmd/main.go**

```go
package main

import "gitlab.company.com/[project_name]/backend/cmd/cmds"

func main() {
	cmds.Execute()
}
```

- [ ] **Step 4: Verify full build**

```bash
go build -o server ./cmd/main.go
```

Expected: `server` binary created, exit 0.

- [ ] **Step 5: Remove test binary**

```bash
rm server
```

---

## Task 16: DevOps Files

**Files:**

- Create: `Dockerfile`
- Create: `docker-compose.yml`
- Create: `.env.example`

- [ ] **Step 1: Write Dockerfile**

```dockerfile
# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY vendor/ vendor/
COPY . .
RUN go build -mod=vendor -o server ./cmd/main.go

# Runtime stage
FROM alpine:3.20
RUN apk add --no-cache tzdata ca-certificates
WORKDIR /app
COPY --from=builder /app/server .
COPY assets/ assets/
ENV ENV_FILE_PATH=/app/.env
EXPOSE 8080
CMD ["./server", "serve"]
```

- [ ] **Step 2: Write docker-compose.yml**

```yaml
version: "3.9"

services:
  api:
    build: .
    ports:
      - "8080:8080"
    env_file: .env
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy

  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: database
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 5s
      retries: 5

volumes:
  postgres_data:
```

- [ ] **Step 3: Write .env.example**

```env
# Application
APP_ENV=local
APP_PREFIX=SDD
APP_PORT=8080
APP_VERSION=1.0.0
LOG_LEVEL=info

# API Key (for external routes)
X_API_KEY=changeme

# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=database
POSTGRES_SSL_MODE=disable
POSTGRES_MAX_OPEN_CONNS=10
POSTGRES_MAX_IDLE_CONNS=5

# JWT
JWT_SECRET=changeme-use-a-long-random-secret-in-production

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# S3 (optional)
S3_BUCKET=
S3_REGION=ap-southeast-1
S3_ACCESS_KEY=
S3_SECRET_KEY=
S3_ENDPOINT=

# Email (optional)
EMAIL_HOST=
EMAIL_PORT=587
EMAIL_USER=
EMAIL_PASSWORD=
EMAIL_FROM=

# OAuth (optional)
OAUTH_GOOGLE_CLIENT_ID=
OAUTH_GOOGLE_CLIENT_SECRET=
OAUTH_GOOGLE_CALLBACK_URL=
```

- [ ] **Step 4: Create OpenAPI scaffold**

```yaml
# docs/openapi/openapi.yaml
openapi: 3.0.3
info:
  title: [project_name] Backend API
  version: 1.0.0
paths: {}
components: {}
```

```markdown
# docs/openapi/README.md

`docs/openapi/openapi.yaml` is source of truth.

- keep reusable schemas, parameters, and responses in `components/`
- keep endpoint path fragments in `paths/`
- update OpenAPI contract for every new endpoint
- keep contract compatible with Echo handlers
```

---

## Task 17: README

**Files:**

- Modify: `README.md`

- [ ] **Step 1: Rewrite README.md**

```markdown
# [project_name] — Backend API

Go backend service for the [project_name] platform. Built on clean architecture (hexagonal / ports & adapters) with Echo, GORM, and PostgreSQL.

## Architecture
```

HTTP → Echo Router → Middleware → Handler → Service (port) → Repository (port) → GORM → PostgreSQL

````

Dependency rule: everything depends inward. Domain has zero knowledge of HTTP or DB. See [docs/architechture/ARCHITECTURE.md](docs/architechture/ARCHITECTURE.md) for full reference.

## Quick Start

### Prerequisites

- Go 1.24+
- Docker + Docker Compose

### 1. Clone and copy env

```bash
cp .env.example .env
# edit .env with your local values
````

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

## Build

```bash
go build -o server ./cmd/main.go
./server serve
```

````

---

## Task 18: Final Verification

- [ ] **Step 1: Full build**

```bash
go build ./...
````

Expected: exit 0, no errors.

- [ ] **Step 2: Run tests**

```bash
go test ./... -v -count=1
```

Expected: `health_test.go` PASS. All others: no test files yet (OK).

- [ ] **Step 3: Vet**

```bash
go vet ./...
```

Expected: exit 0, no warnings.

- [ ] **Step 4: Verify directory structure matches architecture reference**

```bash
find . -type d -not -path '*/.git/*' -not -path '*/vendor/*' | sort
```

Compare against tree in `docs/architechture/ARCHITECTURE.md`.

- [ ] **Step 5: Smoke test binary (no DB required)**

```bash
go run ./cmd/main.go --help
```

Expected: cobra help output showing `serve` subcommand.

---

## Self-Review Checklist

### Spec Coverage

| Requirement                                     | Covered By                                                            |
| ----------------------------------------------- | --------------------------------------------------------------------- |
| Clean architecture / hexagonal                  | All layers follow port-interface pattern                              |
| Request lifecycle middleware chain              | Task 13 — protocol middlewares                                        |
| Dependency injection via `Dependencies` structs | All `NewXxx(deps)` constructors                                       |
| Standardized response models                    | Task 9 — handler/base + handler/common                                |
| Pagination with POST /query                     | handler/base/page.go BasePageRequest                                  |
| JWT auth middleware                             | pkg/jwt + protocol/http.go                                            |
| afterAuth Redis cache pattern                   | protocol/custom_middleware.go (stub)                                  |
| Permission middleware                           | protocol/permission_middleware.go (stub)                              |
| userAction logging                              | protocol/user_action_middleware.go                                    |
| Error handler (domain → HTTP)                   | handler/middleware/error_handler.go                                   |
| Error code convention PREFIX-CODE               | handler/common/common.go                                              |
| Transaction context helpers +`WithTx`           | internal/core/transaction/transaction.go                              |
| Repository tx awareness pattern                 | Documented in repository.go                                           |
| Context typed keys                              | internal/core/domain/context.go                                       |
| Bilingual Name/NameEn pattern                   | internal/core/enums/enums.go                                          |
| Cron with local guard                           | protocol/cron.go                                                      |
| WebSocket manager                               | websocket/websocket.go                                                |
| External API key middleware                     | handler/middleware/api_key_guard.go                                   |
| Swagger/OpenAPI-first API docs                  | docs/openapi/openapi.yaml + docs/openapi/README.md                    |
| AI workspace guidance                           | docs/AGENTS.md + docs/CONTEXT.md + docs/engineering/API_GUIDELINES.md |
| Testing + mock strategy guidance                | docs/engineering/TESTING.md + docs/AGENTS.md                          |
| Example feature (health)                        | Tasks 11-12                                                           |
| Docker + docker-compose                         | Task 16                                                               |
| .env.example                                    | Task 16                                                               |
| README                                          | Task 17                                                               |
| Anti-pattern guards                             | Documented in code comments                                           |
