# Go Development Setup Guide

This guide covers setting up the Go development environment for the P2K16 migration project.

## Prerequisites

- **Go 1.21+**: Install from [golang.org](https://golang.org/dl/)
- **PostgreSQL 13+**: For database development
- **Make**: For build automation
- **Git**: Version control

## Quick Start

```bash
# Clone and setup (if not already done)
git clone https://github.com/helloellinor/p2k16.git
cd p2k16

# Install Go dependencies
go mod download

# Set up development environment
make dev-setup

# Run in demo mode (no database required)
make demo

# Or run with database
make run
```

## Project Structure

```
p2k16/
├── cmd/                    # Application entry points
│   ├── demo/              # Demo mode without database
│   ├── server/            # Main web server
│   └── test/              # Test utilities
├── internal/              # Private application code
│   ├── database/          # Database connection and setup
│   ├── handlers/          # HTTP handlers (controllers)
│   ├── middleware/        # HTTP middleware
│   └── models/            # Data models and repositories
├── static/                # Static assets (CSS, JS, images)
├── templates/             # HTML templates
├── docs/                  # Documentation
│   ├── go/               # Go-specific documentation
│   ├── migration/        # Migration documentation
│   └── development/      # Development guides
└── Makefile              # Build automation
```

## Development Environment

### Environment Variables

Create a `.env` file in the project root:

```env
# Database connection
DB_HOST=localhost
DB_PORT=5432
DB_USER=p2k16
DB_PASSWORD=p2k16
DB_NAME=p2k16

# Server configuration
PORT=8080
GIN_MODE=debug

# Session configuration
SESSION_SECRET=your-session-secret-key

# Development flags
DEMO_MODE=false
LOG_LEVEL=debug
```

### Database Setup

```bash
# Start PostgreSQL (Docker)
docker run --name p2k16-postgres \
  -e POSTGRES_USER=p2k16 \
  -e POSTGRES_PASSWORD=p2k16 \
  -e POSTGRES_DB=p2k16 \
  -p 5432:5432 \
  -d postgres:13

# Run database migrations
make db-migrate

# Verify database setup
make db-check
```

### Development Commands

```bash
# Development workflow
make dev-setup          # Initial setup
make run                # Run server with auto-reload
make test               # Run tests
make lint               # Run linters
make build              # Build binary

# Database management
make db-migrate         # Run migrations
make db-reset           # Reset database
make db-seed            # Seed test data

# Testing
make test-unit          # Unit tests only
make test-integration   # Integration tests only
make test-coverage      # Coverage report

# Demo mode (no database)
make demo               # Run in demo mode
make demo-data          # Generate demo data
```

## Code Organization

### Handler Pattern

```go
// internal/handlers/auth.go
package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/helloellinor/p2k16/internal/models"
)

type AuthHandler struct {
    accountRepo models.AccountRepository
}

func NewAuthHandler(accountRepo models.AccountRepository) *AuthHandler {
    return &AuthHandler{accountRepo: accountRepo}
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Implementation...
    c.JSON(http.StatusOK, gin.H{"status": "success"})
}
```

### Repository Pattern

```go
// internal/models/account.go
package models

import (
    "database/sql"
    "time"
)

type Account struct {
    ID        int       `json:"id" db:"id"`
    Email     string    `json:"email" db:"email"`
    Name      string    `json:"name" db:"name"`
    Enabled   bool      `json:"enabled" db:"enabled"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type AccountRepository interface {
    Create(account *Account) error
    GetByID(id int) (*Account, error)
    GetByEmail(email string) (*Account, error)
    Update(account *Account) error
    Delete(id int) error
}

type PostgresAccountRepository struct {
    db *sql.DB
}

func NewPostgresAccountRepository(db *sql.DB) AccountRepository {
    return &PostgresAccountRepository{db: db}
}

func (r *PostgresAccountRepository) Create(account *Account) error {
    query := `
        INSERT INTO accounts (email, name, enabled, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id`
    
    err := r.db.QueryRow(query, account.Email, account.Name, account.Enabled, time.Now()).Scan(&account.ID)
    return err
}
```

### Middleware Pattern

```go
// internal/middleware/auth.go
package middleware

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"
)

func RequireAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        session := sessions.Default(c)
        userID := session.Get("user_id")
        
        if userID == nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
            c.Abort()
            return
        }
        
        c.Set("user_id", userID)
        c.Next()
    }
}
```

## Testing

### Unit Testing

```go
// internal/handlers/auth_test.go
package handlers

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type MockAccountRepository struct {
    mock.Mock
}

func (m *MockAccountRepository) GetByEmail(email string) (*models.Account, error) {
    args := m.Called(email)
    return args.Get(0).(*models.Account), args.Error(1)
}

func TestAuthHandler_Login(t *testing.T) {
    mockRepo := new(MockAccountRepository)
    handler := NewAuthHandler(mockRepo)
    
    // Test implementation...
    assert.NotNil(t, handler)
}
```

### Integration Testing

```go
// internal/handlers/auth_integration_test.go
package handlers

import (
    "testing"
    "net/http/httptest"
    "github.com/gin-gonic/gin"
)

func TestAuthHandler_Login_Integration(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    // Setup test server
    router := gin.New()
    handler := NewAuthHandler(NewPostgresAccountRepository(db))
    router.POST("/auth/login", handler.Login)
    
    // Test request
    req := httptest.NewRequest("POST", "/auth/login", strings.NewReader(`{"email":"test@example.com","password":"test"}`))
    resp := httptest.NewRecorder()
    
    router.ServeHTTP(resp, req)
    
    assert.Equal(t, 200, resp.Code)
}
```

## HTMX Integration

### Template Structure

```html
<!-- templates/base.html -->
<!DOCTYPE html>
<html>
<head>
    <title>P2K16</title>
    <script src="https://unpkg.com/htmx.org@1.8.4"></script>
    <link rel="stylesheet" href="/static/css/main.css">
</head>
<body>
    <div id="app">
        {{template "content" .}}
    </div>
    
    <div id="notifications" hx-swap-oob="true"></div>
</body>
</html>
```

### HTMX Components

```html
<!-- templates/components/form.html -->
<form hx-post="{{.Action}}" 
      hx-target="{{.Target}}" 
      hx-indicator="#loading"
      class="htmx-form">
    
    {{range .Fields}}
    <div class="form-group">
        <label for="{{.Name}}">{{.Label}}</label>
        <input type="{{.Type}}" 
               name="{{.Name}}" 
               id="{{.Name}}" 
               value="{{.Value}}"
               {{if .Required}}required{{end}}>
    </div>
    {{end}}
    
    <button type="submit" class="btn btn-primary">
        <span class="htmx-indicator" id="loading">Loading...</span>
        {{.SubmitText}}
    </button>
</form>
```

### HTMX Handlers

```go
// internal/handlers/htmx.go
package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func (h *AuthHandler) LoginForm(c *gin.Context) {
    data := gin.H{
        "Action": "/auth/login",
        "Target": "#main-content",
        "Fields": []FormField{
            {Name: "email", Type: "email", Label: "Email", Required: true},
            {Name: "password", Type: "password", Label: "Password", Required: true},
        },
        "SubmitText": "Login",
    }
    
    c.HTML(http.StatusOK, "components/form.html", data)
}
```

## Debugging

### VS Code Configuration

Create `.vscode/launch.json`:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Go Server",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/server",
            "env": {
                "GIN_MODE": "debug",
                "DB_HOST": "localhost"
            },
            "args": []
        }
    ]
}
```

### Debug Logging

```go
// internal/logger/logger.go
package logger

import (
    "github.com/sirupsen/logrus"
    "os"
)

var Log *logrus.Logger

func init() {
    Log = logrus.New()
    Log.SetOutput(os.Stdout)
    
    if os.Getenv("LOG_LEVEL") == "debug" {
        Log.SetLevel(logrus.DebugLevel)
    }
}

// Usage in handlers
import "github.com/helloellinor/p2k16/internal/logger"

func (h *AuthHandler) Login(c *gin.Context) {
    logger.Log.Debug("Login attempt for user", c.PostForm("email"))
    // ...
}
```

## Performance Monitoring

### Metrics Collection

```go
// internal/middleware/metrics.go
package middleware

import (
    "time"
    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
        Name: "http_duration_seconds",
        Help: "Duration of HTTP requests.",
    }, []string{"path"})
)

func MetricsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        duration := time.Since(start)
        httpDuration.WithLabelValues(c.Request.URL.Path).Observe(duration.Seconds())
    }
}
```

## Common Issues

### Database Connection
```bash
# Check database connection
make db-check

# Reset database if corrupted
make db-reset && make db-migrate
```

### Go Module Issues
```bash
# Clean module cache
go clean -modcache

# Update dependencies
go mod tidy && go mod download
```

### Build Issues
```bash
# Clean build cache
go clean -cache

# Rebuild from scratch
make clean && make build
```

## Migration Development

### Parallel Development Setup

```bash
# Terminal 1: Python system
source .settings.fish
p2k16-run-web  # Runs on :5000

# Terminal 2: Go system  
make run PORT=8081  # Runs on :8081

# Terminal 3: Database monitoring
watch -n 2 'psql -h localhost -U p2k16 -d p2k16 -c "SELECT COUNT(*) FROM accounts;"'
```

### Migration Testing

```bash
# Test both systems
curl http://localhost:5000/api/auth/login  # Python
curl http://localhost:8081/api/auth/login  # Go

# Compare responses
make compare-apis
```

## Resources

### Documentation
- [Gin Framework Guide](https://gin-gonic.com/docs/)
- [HTMX Documentation](https://htmx.org/)
- [Go Database/SQL Tutorial](https://go.dev/doc/tutorial/database-access)

### Tools
- [Air](https://github.com/cosmtrek/air) - Live reload for Go
- [golangci-lint](https://golangci-lint.run/) - Go linter
- [Delve](https://github.com/go-delve/delve) - Go debugger

---

**Last Updated**: [Current Date]  
**Go Lead**: [To be assigned]