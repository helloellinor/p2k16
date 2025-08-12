# Deployment Strategy

This document outlines the production deployment strategy for migrating from Python/Flask to Go + HTMX.

## Overview

The deployment strategy ensures zero-downtime migration through parallel deployment, feature flagging, and gradual traffic shifting.

## Deployment Architecture

### Current State (Python Only)
```
                   ┌─────────────────┐
    Internet ──────│  Load Balancer  │
                   │    (Nginx)      │
                   └─────────┬───────┘
                             │
                    ┌────────▼────────┐
                    │ Python Flask    │
                    │ Application     │
                    │ (Port 5000)     │
                    └────────┬────────┘
                             │
                     ┌───────▼────────┐
                     │  PostgreSQL    │
                     │   Database     │
                     └────────────────┘
```

### Migration State (Parallel Deployment)
```
                   ┌─────────────────┐
    Internet ──────│  Load Balancer  │
                   │    (Nginx)      │
                   │  Feature-based  │
                   │    Routing      │
                   └─────────┬───────┘
                             │
                   ┌─────────┼─────────┐
                   │         │         │
          ┌────────▼────────┐ ┌───────▼──────┐
          │ Python Flask    │ │ Go + HTMX    │
          │ Legacy System   │ │ New System   │
          │ (Port 5000)     │ │ (Port 8080)  │
          └────────┬────────┘ └───────┬──────┘
                   │                  │
                   └─────────┬────────┘
                             │
                     ┌───────▼────────┐
                     │  PostgreSQL    │
                     │   Database     │
                     │   (Shared)     │
                     └────────────────┘
```

### Target State (Go Only)
```
                   ┌─────────────────┐
    Internet ──────│  Load Balancer  │
                   │    (Nginx)      │
                   └─────────┬───────┘
                             │
                    ┌────────▼────────┐
                    │ Go + HTMX       │
                    │ Application     │
                    │ (Port 8080)     │
                    └────────┬────────┘
                             │
                     ┌───────▼────────┐
                     │  PostgreSQL    │
                     │   Database     │
                     └────────────────┘
```

## Deployment Phases

### Phase 1: Preparation
**Duration**: 1 week

- [ ] Set up production Go environment
- [ ] Configure monitoring and logging
- [ ] Set up health checks for both systems
- [ ] Create deployment pipelines
- [ ] Test rollback procedures

### Phase 2: Parallel Deployment
**Duration**: 16 weeks (during migration development)

- [ ] Deploy Go system alongside Python
- [ ] Configure feature-based routing
- [ ] Implement session sharing between systems
- [ ] Set up monitoring dashboards
- [ ] Test compatibility between systems

### Phase 3: Gradual Migration
**Duration**: 3 weeks

- [ ] Migrate authentication endpoints (Week 1)
- [ ] Migrate user management endpoints (Week 2)
- [ ] Migrate remaining endpoints (Week 3)
- [ ] Monitor performance and errors
- [ ] Implement rollback triggers

### Phase 4: Full Migration
**Duration**: 1 week

- [ ] Route 100% traffic to Go system
- [ ] Deprecate Python system
- [ ] Clean up deployment configuration
- [ ] Update documentation
- [ ] Archive Python codebase

## Traffic Routing Strategy

### Feature-Based Routing (Nginx Configuration)

```nginx
upstream python_backend {
    server 127.0.0.1:5000;
}

upstream go_backend {
    server 127.0.0.1:8080;
}

server {
    listen 80;
    server_name p2k16.example.com;

    # Route migrated endpoints to Go
    location /api/auth/ {
        proxy_pass http://go_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # Route new HTMX endpoints to Go
    location /htmx/ {
        proxy_pass http://go_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # Route everything else to Python (for now)
    location / {
        proxy_pass http://python_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # Health checks
    location /health/python {
        proxy_pass http://python_backend/health;
    }

    location /health/go {
        proxy_pass http://go_backend/health;
    }
}
```

### A/B Testing Configuration

```nginx
# Split traffic based on user ID for gradual rollout
map $cookie_user_id $backend_pool {
    ~*[0-4]$ go_backend;      # 50% to Go
    default  python_backend;   # 50% to Python
}

location /api/users/ {
    proxy_pass http://$backend_pool;
}
```

## Session Management During Migration

### Session Compatibility Layer

```go
// internal/session/compatibility.go
type SessionCompatibility struct {
    pythonStore sessions.Store
    goStore     sessions.Store
}

func (s *SessionCompatibility) MigrateSession(c *gin.Context) {
    // Check for Python session
    pythonSession := s.pythonStore.Get(c.Request, "session")
    if pythonSession != nil && !pythonSession.IsNew {
        // Migrate to Go session format
        goSession := s.goStore.Get(c.Request, "session")
        goSession.Values["user_id"] = pythonSession.Values["user_id"]
        goSession.Values["migrated"] = true
        goSession.Save(c.Request, c.Writer)
    }
}
```

### Session Sharing Strategy

```python
# Python session compatibility (web/session_compat.py)
import redis

class SessionCompat:
    def __init__(self):
        self.redis_client = redis.Redis(host='localhost', port=6379, db=1)
    
    def share_session(self, session_id, user_data):
        # Store session in Redis for Go access
        self.redis_client.setex(
            f"session:{session_id}", 
            3600,  # 1 hour
            json.dumps(user_data)
        )
```

## Database Migration Strategy

### Schema Compatibility

```sql
-- Ensure schema compatibility during migration
-- Add any missing columns for Go system
ALTER TABLE accounts ADD COLUMN IF NOT EXISTS last_login_go TIMESTAMP;
ALTER TABLE sessions ADD COLUMN IF NOT EXISTS created_by_system VARCHAR(10) DEFAULT 'python';

-- Create indices for Go system performance
CREATE INDEX IF NOT EXISTS idx_accounts_email_enabled ON accounts(email, enabled);
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
```

### Data Migration Scripts

```bash
#!/bin/bash
# scripts/migrate-data.sh

echo "Starting data migration..."

# Backup current database
pg_dump p2k16 > "backup_$(date +%Y%m%d_%H%M%S).sql"

# Run schema updates
psql -h localhost -U p2k16 -d p2k16 -f migrations/go_compatibility.sql

# Verify data integrity
python scripts/verify_migration.py

echo "Data migration complete"
```

## Monitoring & Alerting

### Key Metrics to Monitor

```yaml
# Prometheus monitoring config
- name: response_time
  query: histogram_quantile(0.95, http_request_duration_seconds)
  threshold: 0.5  # 500ms
  
- name: error_rate
  query: rate(http_requests_total{status=~"4..|5.."}[5m])
  threshold: 0.01  # 1%
  
- name: database_connections
  query: postgres_connections_active
  threshold: 50
  
- name: memory_usage
  query: process_resident_memory_bytes
  threshold: 536870912  # 512MB
```

### Alerting Rules

```yaml
# Alert when Go system has higher error rate than Python
- alert: GoSystemErrors
  expr: rate(http_requests_total{service="go",status=~"5.."}[5m]) > 
        rate(http_requests_total{service="python",status=~"5.."}[5m]) * 2
  for: 2m
  annotations:
    summary: "Go system error rate exceeds Python by 2x"
    
# Alert when migration causes database issues
- alert: DatabaseConnectionSpike
  expr: postgres_connections_active > 100
  for: 1m
  annotations:
    summary: "Database connections spiking during migration"
```

## Rollback Procedures

### Automatic Rollback Triggers

```bash
#!/bin/bash
# scripts/auto-rollback.sh

# Check error rates every 30 seconds
while true; do
    GO_ERROR_RATE=$(curl -s "http://prometheus:9090/api/v1/query?query=rate(http_requests_total{service=\"go\",status=~\"5..\"}[2m])" | jq '.data.result[0].value[1]' | tr -d '"')
    
    if (( $(echo "$GO_ERROR_RATE > 0.05" | bc -l) )); then
        echo "ERROR: Go system error rate too high ($GO_ERROR_RATE), rolling back..."
        ./rollback-to-python.sh
        break
    fi
    
    sleep 30
done
```

### Manual Rollback Process

```bash
#!/bin/bash
# scripts/rollback-to-python.sh

echo "Rolling back to Python system..."

# 1. Update Nginx configuration
cp nginx/python-only.conf /etc/nginx/sites-available/p2k16
nginx -s reload

# 2. Stop Go system
systemctl stop p2k16-go

# 3. Ensure Python system is running
systemctl start p2k16-python
systemctl status p2k16-python

# 4. Migrate sessions back to Python format if needed
python scripts/migrate-sessions-to-python.py

echo "Rollback complete. Python system active."
```

## Performance Validation

### Load Testing

```yaml
# artillery load test config
config:
  target: 'http://localhost'
  phases:
    - duration: 60
      arrivalRate: 10
    - duration: 120
      arrivalRate: 50
    - duration: 60
      arrivalRate: 100

scenarios:
  - name: "Login flow"
    requests:
      - post:
          url: "/api/auth/login"
          json:
            email: "test@example.com"
            password: "test"
  
  - name: "Dashboard access"
    requests:
      - get:
          url: "/dashboard"
          headers:
            Cookie: "session_id=test_session"
```

### Performance Benchmarks

```bash
# Run performance comparison
./scripts/benchmark-comparison.sh

# Expected results:
# Python Flask: ~150ms avg response time
# Go + HTMX:    ~50ms avg response time (3x improvement)
```

## Security Considerations

### Session Security

```go
// Ensure secure session configuration
store := sessions.NewCookieStore([]byte(sessionSecret))
store.Options = &sessions.Options{
    Path:     "/",
    MaxAge:   3600,
    HttpOnly: true,
    Secure:   true,  // HTTPS only in production
    SameSite: http.SameSiteStrictMode,
}
```

### HTTPS Configuration

```nginx
server {
    listen 443 ssl http2;
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    # Security headers
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains";
}
```

## Deployment Checklist

### Pre-Migration
- [ ] Database backup completed
- [ ] Go system deployed and tested
- [ ] Monitoring configured
- [ ] Rollback procedures tested
- [ ] Performance benchmarks established

### During Migration
- [ ] Traffic routing configured
- [ ] Session compatibility verified
- [ ] Error rates monitored
- [ ] Performance metrics tracked
- [ ] Rollback triggers active

### Post-Migration
- [ ] All traffic routed to Go system
- [ ] Python system deprecated
- [ ] Monitoring updated
- [ ] Documentation updated
- [ ] Team trained on new system

---

**Deployment Lead**: [To be assigned]  
**Last Updated**: [Current Date]  
**Next Review**: Weekly during migration