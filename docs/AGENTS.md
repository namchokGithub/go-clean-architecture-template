# AI Agent Guide

## Purpose

Guide for Codex, Claude Code, Cursor, and other coding agents working in this repo.

## Architecture Rules

- Preserve Clean Architecture / Hexagonal direction
- Preserve flow: `Handler -> Service -> Repository`
- Use explicit `Dependencies` structs and `NewXxx(deps)` constructors
- Keep `context.Context` as first argument for service and repository methods
- Keep PostgreSQL-first and Echo-compatible direction

## Layering Rules

- Handlers: bind request, validate, call service, render response
- Services: own business logic, orchestration, transaction boundaries
- Repositories: own DB access only
- No business logic in handlers
- Repositories must not call other repositories
- Services may coordinate multiple repositories
- Prefer one port file per feature under `internal/core/port/`

## Coding Rules

- Keep changes production-oriented and Go idiomatic
- Avoid generic `utils/` package
- Prefer feature isolation over cross-cutting grab-bags
- Reuse shared technical packages only when transport-neutral
- Return domain types from repositories, not GORM models

## Transaction Rules

- Use `transaction.WithTx(...)` for transactional orchestration in service layer
- Keep transaction propagation in context
- Repositories must auto-detect tx via `transaction.FromContext(ctx)`
- Avoid manual commit/rollback flow in feature services unless strong reason

## API Documentation Rules

- Swagger/OpenAPI is source of truth
- Keep contract in `docs/openapi/openapi.yaml`
- Reuse schemas, parameters, and responses under `docs/openapi/components/`
- Update OpenAPI contract for every new endpoint
- Keep request/response shape compatible with Echo handlers

## Mock Strategy

The project uses `mockery` for interface-based mock generation.

Rules:

- Generate mocks from `internal/core/port`
- Store generated files in `tests/mocks/`
- Do NOT mock concrete implementations
- Do NOT manually edit generated mocks
- Prefer testing service behavior over implementation details

Example:

```bash
mockery --dir internal/core/port \
        --output tests/mocks \
        --all
```

## Documentation Update Rules

When adding or changing feature:

1. Update domain types if needed
2. Update `internal/core/port/[feature].go`
3. Update service and repository aggregates if needed
4. Register routes in `protocol/http.go`
5. Update `docs/openapi/openapi.yaml` and related components/paths
6. Update `docs/CONTEXT.md` only if project understanding changed materially

## Defaults for Agents

- Prefer small, targeted edits
- Do not rewrite architecture direction without explicit request
- Do not invent business requirements, endpoints, or integrations
- Leave TODO placeholders where repo context still missing
