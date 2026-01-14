# Testing Guide

## Overview

This project includes comprehensive testing at multiple levels:
- **Unit Tests:** Test individual functions and components
- **Integration Tests:** Test database operations and tenant isolation
- **E2E Tests:** Test complete user flows
- **Security Tests:** Test security mechanisms
- **Load Tests:** Test performance under load
- **Frontend Tests:** Test React components and stores

## Backend Tests (Go)

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test file
go test ./tests/integration/tenant_isolation_test.go

# Run with verbose output
go test -v ./...

# Run only integration tests
go test ./tests/integration/...

# Run only security tests
go test ./tests/security/...
```

### Test Structure

```
tests/
├── testutil/               # Test utilities and helpers
│   └── database.go         # Database test helpers
├── integration/            # Integration tests
│   ├── tenant_isolation_test.go
│   └── tier_enforcement_test.go
├── e2e/                    # End-to-end tests
│   └── subscription_flow_test.go
├── security/               # Security tests
│   └── security_test.go
└── load/                   # Load tests
    └── k6-load-test.js
```

### Test Database Setup

Integration tests use temporary MySQL databases:

1. Tests create a unique database per test run
2. Schema is migrated automatically
3. Test data is seeded
4. Database is cleaned up after tests

**Requirements:**
- MySQL server running on localhost:3306
- User: root
- Password: test_password

### Writing Tests

Example test structure:

```go
func TestMyFeature(t *testing.T) {
    // Setup
    testDB := testutil.SetupTestDB(t)
    defer testDB.Teardown(t)
    testDB.RunMigrations(t)

    // Test
    t.Run("MySubTest", func(t *testing.T) {
        // Arrange
        ctx := context.Background()

        // Act
        result, err := myFunction(ctx)

        // Assert
        if err != nil {
            t.Fatalf("Expected no error, got %v", err)
        }
        if result != expected {
            t.Errorf("Expected %v, got %v", expected, result)
        }
    })
}
```

### Tenant Isolation Tests

Critical tests to prevent cross-tenant data leaks:

```bash
# Run tenant isolation tests
go test -v ./tests/integration/tenant_isolation_test.go

# Check coverage
go test -cover ./tests/integration/tenant_isolation_test.go
```

**What's tested:**
- ✅ Gin repository isolation
- ✅ User repository isolation
- ✅ Audit log isolation
- ✅ Count queries are tenant-scoped
- ✅ Cross-tenant data access prevention

### Tier Enforcement Tests

Tests for subscription tier limits:

```bash
go test -v ./tests/integration/tier_enforcement_test.go
```

**What's tested:**
- ✅ Gin count limits per tier
- ✅ Photo limits per gin
- ✅ Feature access per tier
- ✅ Multi-user restriction (Enterprise only)
- ✅ Storage limits

### Security Tests

```bash
go test -v ./tests/security/security_test.go
```

**What's tested:**
- ✅ Password hashing
- ✅ SQL injection prevention
- ✅ API key format
- ✅ XSS prevention
- ✅ Authorization bypass prevention
- ✅ Audit logging

## Load Tests (k6)

### Running Load Tests

```bash
# Install k6
# macOS: brew install k6
# Linux: sudo apt install k6
# Windows: choco install k6

# Run load test
cd tests/load
k6 run k6-load-test.js

# Run with custom duration
k6 run --duration 30m k6-load-test.js

# Run against production (careful!)
k6 run --env BASE_URL=https://api.ginapp.com k6-load-test.js

# Run with HTML report
k6 run --out html=report.html k6-load-test.js
```

### Load Test Scenarios

The load test simulates:
- 100 concurrent users → 200 concurrent users
- Creating gins
- Listing gins
- Searching
- Getting stats
- Updating gins

### Performance Thresholds

- p95 response time < 500ms
- Error rate < 5%
- All endpoints respond within 1 second

## Frontend Tests (Vitest)

### Running Frontend Tests

```bash
cd frontend

# Run tests
npm test

# Run with coverage
npm run test:coverage

# Run in watch mode
npm run test:watch

# Run UI mode
npm run test:ui
```

### Test Structure

```
frontend/src/
├── stores/
│   └── __tests__/
│       └── authStore.test.ts
├── components/
│   └── __tests__/
│       └── Layout.test.tsx
└── tests/
    └── setup.ts
```

### Writing Frontend Tests

Example component test:

```typescript
import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { Login } from '../Login';

describe('Login Component', () => {
  it('renders login form', () => {
    render(<Login />);
    expect(screen.getByLabelText('Email')).toBeInTheDocument();
    expect(screen.getByLabelText('Password')).toBeInTheDocument();
  });
});
```

Example store test:

```typescript
import { describe, it, expect } from 'vitest';
import { useAuthStore } from '../authStore';

describe('authStore', () => {
  it('should login successfully', async () => {
    const { login } = useAuthStore.getState();
    await login('test@example.com', 'password');

    const state = useAuthStore.getState();
    expect(state.isAuthenticated).toBe(true);
  });
});
```

## Continuous Integration

Tests run automatically on:
- Pull requests
- Commits to main branch
- Nightly builds

### GitHub Actions Workflow

```yaml
name: Tests

on: [push, pull_request]

jobs:
  backend-tests:
    runs-on: ubuntu-latest
    services:
      mysql:
        image: mysql:8.0
        env:
          MYSQL_ROOT_PASSWORD: test_password
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - run: go test -v -cover ./...

  frontend-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
      - run: npm ci
      - run: npm test
```

## Test Coverage Goals

| Component | Target | Current |
|-----------|--------|---------|
| Backend Repositories | 90% | TBD |
| Backend Services | 85% | TBD |
| Frontend Stores | 80% | TBD |
| Frontend Components | 70% | TBD |
| **Overall** | **80%** | **TBD** |

## Best Practices

### General
- ✅ Write tests before fixing bugs
- ✅ Test behavior, not implementation
- ✅ Use descriptive test names
- ✅ Keep tests independent
- ✅ Clean up after tests

### Backend
- ✅ Use table-driven tests
- ✅ Mock external dependencies
- ✅ Test error cases
- ✅ Verify tenant isolation
- ✅ Test transaction rollbacks

### Frontend
- ✅ Test user interactions
- ✅ Mock API calls
- ✅ Test loading/error states
- ✅ Test accessibility
- ✅ Snapshot critical UIs

## Debugging Tests

### Backend

```bash
# Run with race detector
go test -race ./...

# Run with CPU profiling
go test -cpuprofile=cpu.prof ./...

# Run specific test
go test -run TestTenantIsolation_GinRepository ./tests/integration/...

# Run with verbose output and test logs
go test -v -args -test.v ./...
```

### Frontend

```bash
# Debug specific test
npm test -- authStore.test.ts

# Run with Node inspector
node --inspect-brk ./node_modules/vitest/vitest.mjs run
```

## Common Issues

### Backend

**Issue:** Tests fail with "database not found"
**Solution:** Ensure MySQL is running and accessible

**Issue:** Tenant isolation test fails
**Solution:** Check that all queries include `WHERE tenant_id = ?`

### Frontend

**Issue:** Tests timeout
**Solution:** Increase timeout in vitest.config.ts

**Issue:** Module not found
**Solution:** Check path aliases in vitest.config.ts

## Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Vitest Documentation](https://vitest.dev/)
- [k6 Documentation](https://k6.io/docs/)
- [React Testing Library](https://testing-library.com/react)
