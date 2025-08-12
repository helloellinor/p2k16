# API Compatibility Guide

This document outlines the API compatibility strategy between the Python/Flask backend and the new Go + HTMX system during migration.

## Overview

To ensure zero-downtime migration, the Go backend must maintain API compatibility with existing Python endpoints while introducing new HTMX-enhanced interfaces.

## Compatibility Strategy

### 1. Endpoint Mapping

| Python Route | HTTP Method | Go Route | Status | Notes |
|-------------|-------------|----------|--------|-------|
| `/api/auth/login` | POST | `/api/auth/login` | ✅ Compatible | Session handling differs |
| `/api/auth/logout` | POST | `/api/auth/logout` | ✅ Compatible | |
| `/api/accounts/` | GET | `/api/accounts/` | ❌ Not implemented | User listing |
| `/api/accounts/<id>` | GET | `/api/accounts/<id>` | ❌ Not implemented | User details |
| `/api/badges/` | GET | `/api/badges/` | ❌ Not implemented | Badge listing |
| `/api/tools/` | GET | `/api/tools/` | 🚧 Partial | Basic tool listing only |
| `/api/memberships/` | GET | `/api/memberships/` | ❌ Not implemented | Membership status |

### 2. Response Format Compatibility

#### Python Response Format
```json
{
  "status": "success",
  "data": {...},
  "message": "Optional message"
}
```

#### Go Response Format (Compatible)
```json
{
  "status": "success", 
  "data": {...},
  "message": "Optional message"
}
```

### 3. Authentication Compatibility

#### Session Management
- **Python**: Flask-Session with server-side storage
- **Go**: HTTP-only cookies with server-side session store
- **Compatibility**: Session migration middleware converts between formats

#### Authorization Headers
Both systems support the same authorization patterns:
- Cookie-based sessions (primary)
- Bearer token authentication (for API clients)

## Testing Strategy

### Automated Compatibility Tests

Create test suites that verify API compatibility:

```bash
# Run compatibility tests
npm run test-api-compat

# Test specific endpoints
npm run test-auth-compat
npm run test-user-compat
```

### Test Implementation

```javascript
// Example compatibility test
describe('Auth API Compatibility', () => {
  test('login endpoint returns same format', async () => {
    const pythonResponse = await fetch('http://localhost:5000/api/auth/login', {
      method: 'POST',
      body: JSON.stringify({email: 'test@example.com', password: 'test'})
    });
    
    const goResponse = await fetch('http://localhost:8080/api/auth/login', {
      method: 'POST', 
      body: JSON.stringify({email: 'test@example.com', password: 'test'})
    });
    
    expect(goResponse.data).toMatchSchema(pythonResponse.data);
  });
});
```

## Migration Phases

### Phase 1: Authentication (✅ Complete)
- [x] Login/logout API compatibility
- [x] Session format compatibility
- [x] Error response compatibility

### Phase 2: User Management (🚧 In Progress)
- [ ] User listing API (`/api/accounts/`)
- [ ] User detail API (`/api/accounts/<id>`)
- [ ] User profile update API
- [ ] Password reset API

### Phase 3: Badge System (❌ Not Started)
- [ ] Badge listing API (`/api/badges/`)
- [ ] Badge detail API (`/api/badges/<id>`)
- [ ] Badge award API (`/api/badges/<id>/award`)

### Phase 4: Tools & Access (❌ Not Started)
- [ ] Tool listing API (`/api/tools/`)
- [ ] Tool checkout API (`/api/tools/<id>/checkout`)
- [ ] Door access API (`/api/doors/<id>/access`)

### Phase 5: Membership & Payments (❌ Not Started)
- [ ] Membership status API (`/api/memberships/`)
- [ ] Payment processing API (`/api/payments/`)
- [ ] Invoice API (`/api/invoices/`)

## Breaking Changes (To Avoid)

### Avoid These Changes
- Changing response format structure
- Modifying required request parameters
- Changing authentication mechanisms
- Altering error code meanings

### Acceptable Changes
- Adding optional response fields
- Adding optional request parameters
- Improving performance
- Enhancing error messages

## Rollback Strategy

### Immediate Rollback
If compatibility issues arise:

1. **Route Traffic Back to Python**
   ```nginx
   # Nginx config change
   location /api/ {
     proxy_pass http://python-backend:5000;
   }
   ```

2. **Database Consistency Check**
   ```sql
   -- Verify no data corruption
   SELECT COUNT(*) FROM accounts WHERE created_at > NOW() - INTERVAL '1 hour';
   ```

3. **Session Migration**
   ```bash
   # Migrate Go sessions back to Python format
   python scripts/migrate_sessions.py --direction=go-to-python
   ```

## Monitoring & Alerts

### Key Metrics to Monitor
- **API Response Time**: Go vs Python comparison
- **Error Rates**: 4xx/5xx errors by endpoint
- **Session Compatibility**: Cross-system session usage
- **Data Consistency**: Database state validation

### Alert Thresholds
- Response time > 500ms (Go should be <100ms)
- Error rate > 1% for migrated endpoints
- Session failures > 0.1% of total sessions

## Tools & Scripts

### Compatibility Testing Tools
```bash
# Install testing dependencies
npm install --save-dev api-compatibility-suite

# Run full compatibility test suite
npm run test:compatibility

# Generate compatibility report
npm run report:compatibility
```

### Migration Helper Scripts
```bash
# Validate API compatibility before deployment
./scripts/validate-compatibility.sh

# Monitor compatibility during migration
./scripts/monitor-compatibility.sh --duration=1h
```

## Next Steps

1. **Complete Phase 2 APIs**: Implement user management endpoints with full compatibility
2. **Enhance Test Coverage**: Add more comprehensive compatibility tests
3. **Performance Baselines**: Establish performance comparison benchmarks
4. **Monitoring Setup**: Deploy compatibility monitoring in production

---

*Last updated: [Current Date]*  
*Next review: Weekly during active migration phases*