# API Guidelines

## Contract

- OpenAPI-first
- Source of truth: `docs/openapi/openapi.yaml`
- Every new endpoint must update contract in same change
- Reuse schemas, parameters, and responses where possible

## Handler Rules

- Keep Echo handlers thin
- Bind and validate input in handler layer
- Delegate business logic to service layer
- Return consistent response envelopes

## Request Design

- Use `context.Context` through all downstream calls
- Paginated queries may use `POST /query` when filter body is needed
- Keep field names stable and explicit
- Prefer feature-based route grouping

## Response Design

- Reuse common success and error response shapes
- Include machine-readable error code
- Keep error payload predictable for clients and agents

## Transactions

- Start transactional orchestration in service layer
- Use `transaction.WithTx(...)`
- Do not open or commit transactions in handler layer

## Persistence Boundary

- Repository layer owns SQL and GORM access
- Repositories return domain types
- Repositories do not call other repositories

## Change Checklist

1. Update OpenAPI contract
2. Update handler, service, repository flow
3. Reuse shared schemas if possible
4. Keep Echo behavior aligned with documented contract
