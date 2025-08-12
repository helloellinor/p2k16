# Migration Testing Strategy

This document outlines the comprehensive testing approach for the P2K16 migration from Python/Flask to Go + HTMX.

## Testing Overview

The migration testing strategy ensures:
- **Zero data loss** during migration
- **Feature parity** between Python and Go systems
- **Performance improvements** are measurable
- **Rollback capability** if issues arise

## Testing Pyramid

```
    ┌─────────────────┐
    │   E2E Tests     │  ← User workflows, browser automation
    │                 │
    ├─────────────────┤
    │ Integration     │  ← API compatibility, database consistency
    │ Tests           │
    ├─────────────────┤
    │   Unit Tests    │  ← Go handlers, repositories, models
    │                 │
    └─────────────────┘
```

## Test Categories

### 1. Unit Tests

**Purpose**: Test individual Go components in isolation

**Coverage Areas**:
- Database models and validation
- Repository layer functionality
- Handler business logic
- Middleware behavior
- Utility functions

**Example**:
```go
func TestAccountRepository_Create(t *testing.T) {
    repo := setupTestDB(t)
    defer cleanupTestDB(t, repo)
    
    account := &models.Account{
        Email: "test@example.com",
        Name:  "Test User",
    }
    
    err := repo.Create(account)
    assert.NoError(t, err)
    assert.NotZero(t, account.ID)
}
```

**Run Commands**:
```bash
# Run all unit tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/handlers
```

### 2. Integration Tests

**Purpose**: Test interactions between components and external systems

**Coverage Areas**:
- Database integration
- HTTP endpoint behavior
- Session management
- Authentication flows
- HTMX interactions

**Example**:
```go
func TestAuthHandler_Login_Integration(t *testing.T) {
    app := setupTestApp(t)
    defer teardownTestApp(t, app)
    
    // Create test user
    user := createTestUser(t, app.db)
    
    // Test login request
    req := httptest.NewRequest("POST", "/api/auth/login", 
        strings.NewReader(`{"email":"test@example.com","password":"test"}`))
    resp := httptest.NewRecorder()
    
    app.router.ServeHTTP(resp, req)
    
    assert.Equal(t, 200, resp.Code)
    assert.Contains(t, resp.Header().Get("Set-Cookie"), "session_id")
}
```

### 3. API Compatibility Tests

**Purpose**: Ensure Go endpoints match Python API behavior

**Coverage Areas**:
- Request/response format compatibility
- Error handling consistency
- Authentication behavior
- Data validation rules

**Example**:
```javascript
// Jest/Node.js compatibility test
describe('API Compatibility', () => {
  test('login endpoint compatibility', async () => {
    const pythonResp = await fetch('http://localhost:5000/api/auth/login', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({email: 'test@example.com', password: 'test'})
    });
    
    const goResp = await fetch('http://localhost:8080/api/auth/login', {
      method: 'POST', 
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({email: 'test@example.com', password: 'test'})
    });
    
    expect(goResp.status).toBe(pythonResp.status);
    expect(await goResp.json()).toMatchObject(await pythonResp.json());
  });
});
```

### 4. End-to-End Tests

**Purpose**: Test complete user workflows across both systems

**Coverage Areas**:
- User registration and login flows
- Badge management workflows  
- Tool checkout processes
- Membership management
- Payment processing

**Example**:
```javascript
// Playwright E2E test
test('user can login and view dashboard', async ({ page }) => {
  await page.goto('http://localhost:8080');
  
  await page.fill('#email', 'test@example.com');
  await page.fill('#password', 'testpassword');
  await page.click('#login-button');
  
  await expect(page.locator('#dashboard')).toBeVisible();
  await expect(page.locator('#user-name')).toContainText('Test User');
});
```

### 5. Performance Tests

**Purpose**: Verify Go system meets performance targets

**Coverage Areas**:
- Response time benchmarks
- Memory usage monitoring
- Concurrent user handling
- Database query performance

**Example**:
```go
func BenchmarkLoginHandler(b *testing.B) {
    app := setupBenchmarkApp(b)
    defer teardownBenchmarkApp(b, app)
    
    req := httptest.NewRequest("POST", "/api/auth/login", 
        strings.NewReader(`{"email":"test@example.com","password":"test"}`))
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        resp := httptest.NewRecorder()
        app.router.ServeHTTP(resp, req)
    }
}
```

## Test Environment Setup

### Local Development Testing

```bash
# Set up test database
make test-db-setup

# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run compatibility tests (requires both systems)
make test-compatibility
```

### CI/CD Pipeline Testing

```yaml
# .github/workflows/test.yml
name: Test Suite
on: [push, pull_request]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: '1.21'
      - run: make test
      
  integration-tests:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:13
        env:
          POSTGRES_PASSWORD: test
        options: --health-cmd pg_isready --health-interval 10s
    steps:
      - uses: actions/checkout@v3
      - run: make test-integration
      
  compatibility-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - run: make test-compatibility
```

## Test Data Management

### Test Database Setup

```sql
-- test_data.sql
INSERT INTO accounts (email, name, enabled) VALUES 
('test1@example.com', 'Test User 1', true),
('test2@example.com', 'Test User 2', true),
('admin@example.com', 'Admin User', true);

INSERT INTO badges (name, description) VALUES
('Basic Access', 'Basic facility access'),
('Tool Certified', 'Certified for basic tools');
```

### Test Data Factories

```go
// testdata/factories.go
func CreateTestAccount(db *sql.DB, email string) *models.Account {
    account := &models.Account{
        Email:   email,
        Name:    "Test User",
        Enabled: true,
    }
    
    err := db.Create(account)
    if err != nil {
        panic(err)
    }
    
    return account
}

func CreateTestBadge(db *sql.DB, name string) *models.Badge {
    badge := &models.Badge{
        Name:        name,
        Description: "Test badge",
    }
    
    err := db.Create(badge)
    if err != nil {
        panic(err)
    }
    
    return badge
}
```

## Phase-Specific Testing

### Phase 1: Foundation Testing
- [x] Authentication flow tests
- [x] Session management tests  
- [x] Basic HTMX interaction tests
- [ ] Session validation middleware tests
- [ ] Error handling tests

### Phase 2: User & Badge Testing
- [ ] User CRUD operation tests
- [ ] Badge management tests
- [ ] User profile interface tests
- [ ] Badge awarding workflow tests

### Phase 3: Membership Testing
- [ ] Membership status tests
- [ ] Payment processing tests
- [ ] Circle management tests
- [ ] Permission system tests

### Phase 4: Tools & Doors Testing
- [ ] Tool checkout/checkin tests
- [ ] Door access control tests
- [ ] Real-time status update tests
- [ ] Access logging tests

### Phase 5: Events & Reports Testing
- [ ] Event logging tests
- [ ] Report generation tests
- [ ] Analytics dashboard tests
- [ ] Data export tests

## Test Automation

### Continuous Testing

```bash
# Watch mode for development
make test-watch

# Pre-commit testing
make test-pre-commit

# Deploy validation testing
make test-deploy-validate
```

### Test Reporting

```bash
# Generate test coverage report
make coverage-report

# Generate compatibility report  
make compatibility-report

# Generate performance benchmark report
make benchmark-report
```

## Quality Gates

### Code Coverage Requirements
- **Unit Tests**: Minimum 80% coverage
- **Integration Tests**: Critical paths covered
- **API Compatibility**: 100% endpoint compatibility

### Performance Requirements
- **Response Time**: <100ms for 95% of requests
- **Memory Usage**: <50MB under normal load
- **Concurrent Users**: Support 100+ simultaneous users

### Compatibility Requirements
- **API Parity**: 100% compatibility for migrated endpoints
- **Data Integrity**: Zero data loss during migration
- **Session Compatibility**: Seamless session migration

## Troubleshooting

### Common Test Issues

#### Database Connection Issues
```bash
# Check test database connection
make test-db-check

# Reset test database
make test-db-reset
```

#### Test Data Cleanup
```bash
# Clean test data between runs
make test-cleanup

# Verify test isolation
make test-isolation-check
```

#### Compatibility Test Failures
```bash
# Debug compatibility issues
make debug-compatibility

# Compare API responses
make compare-api-responses
```

## Tools & Dependencies

### Go Testing Tools
```go
// go.mod testing dependencies
require (
    github.com/stretchr/testify v1.8.0
    github.com/DATA-DOG/go-sqlmock v1.5.0
    github.com/testcontainers/testcontainers-go v0.15.0
)
```

### JavaScript Testing Tools
```json
{
  "devDependencies": {
    "jest": "^29.0.0",
    "@playwright/test": "^1.25.0",
    "supertest": "^6.2.0"
  }
}
```

## Next Steps

1. **Complete Phase 1 Tests**: Finish foundation testing coverage
2. **Set Up CI/CD**: Implement automated testing pipeline
3. **Performance Baselines**: Establish benchmark comparisons
4. **Compatibility Monitoring**: Deploy automated compatibility checking

---

*Last updated: [Current Date]*  
*Testing lead: [To be assigned]*