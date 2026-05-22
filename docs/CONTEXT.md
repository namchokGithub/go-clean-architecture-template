# Project Context

## Overview

- Project type: Golang backend API
- Runtime direction: Echo + PostgreSQL
- Architecture: Clean Architecture / Hexagonal
- Dependency style: explicit constructor injection
- Request style: context-first across service and repository layers
- API contract: Swagger/OpenAPI under `docs/openapi/`

## Current Phase

- Early initialization / scaffold phase
- Core engineering foundation being established
- Business flows not fully defined yet

## Business / Domain

- Product/domain summary: TODO
- Primary actors: TODO
- Key business workflows: TODO
- Critical invariants: TODO

## External Integrations

- PostgreSQL: primary datastore
- Redis: caching / middleware support
- S3: optional file storage
- Other integrations: TODO

## Future Modules

- Auth / identity: TODO
- User / role / permission: TODO
- Core business features: TODO
- Background jobs / cron: TODO
- Realtime / websocket flows: TODO

## Architecture Summary

- Request flow: `Handler -> Service -> Repository`
- Handlers stay thin
- Services own orchestration and transactions
- Repositories own persistence
- Transactions use `transaction.WithTx(...)`
- Repositories read tx from context when present

## Notes for Agents

- Do not assume missing business details
- Prefer TODO over invented requirements
- Keep docs and OpenAPI contract aligned with implemented behavior
