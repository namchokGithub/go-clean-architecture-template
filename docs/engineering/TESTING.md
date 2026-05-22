## Testing & Mock Strategy

The project uses `mockery` for generating mocks from port interfaces.

### Rules

- Generate mocks from interfaces only
- Do NOT mock concrete implementations
- Store generated mocks under:

```txt
tests/mocks/
```

- Do NOT manually edit generated mock files
- Service layer tests should mock repositories and external dependencies through ports

### Recommended Command

```bash
mockery --dir internal/core/port \
        --output tests/mocks \
        --all
```

### Testing Direction

- Unit tests should live close to features (`*_test.go`)
- Integration tests should live under:

```txt
tests/integration/
```

- Prefer behavior-focused tests over excessive mocking
- Repository logic should be validated with PostgreSQL integration tests when possible
