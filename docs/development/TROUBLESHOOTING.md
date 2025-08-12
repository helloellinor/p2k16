# Development Setup Troubleshooting

This guide helps resolve common issues encountered during P2K16 development setup and migration.

## Common Issues

### Database Connection Issues

#### Problem: "connection refused" or "database does not exist"
```
Error: dial tcp [::1]:5432: connect: connection refused
```

**Solutions:**
1. **Check PostgreSQL is running**
   ```bash
   # macOS (Homebrew)
   brew services start postgresql
   
   # Linux (systemd)
   sudo systemctl start postgresql
   
   # Docker
   docker run --name p2k16-postgres \
     -e POSTGRES_USER=p2k16 \
     -e POSTGRES_PASSWORD=p2k16 \
     -e POSTGRES_DB=p2k16 \
     -p 5432:5432 \
     -d postgres:13
   ```

2. **Verify database exists**
   ```bash
   psql -h localhost -U p2k16 -l
   # Should show p2k16 database
   ```

3. **Create database if missing**
   ```bash
   createdb -h localhost -U p2k16 p2k16
   ```

#### Problem: "relation does not exist" errors
```
Error: relation "accounts" does not exist
```

**Solution: Run database migrations**
```bash
# Check if migrations directory exists
ls migrations/

# Run migrations (if using Flyway)
flyway -url=jdbc:postgresql://localhost:5432/p2k16 \
       -user=p2k16 \
       -password=p2k16 \
       migrate

# Or run SQL directly
psql -h localhost -U p2k16 -d p2k16 -f database-setup.sql
```

#### Problem: "column does not exist" (e.g., missing created_by)
```
Error: column "created_by" of relation "accounts" does not exist
```

**Solution: Update database schema**
```sql
-- Connect to database
psql -h localhost -U p2k16 -d p2k16

-- Add missing columns
ALTER TABLE accounts ADD COLUMN IF NOT EXISTS created_by INT;
ALTER TABLE accounts ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT NOW();

-- Verify schema
\d accounts
```

### Go Application Issues

#### Problem: Go server won't start
```
Error: listen tcp :8080: bind: address already in use
```

**Solutions:**
1. **Check what's using the port**
   ```bash
   lsof -i :8080
   # or
   netstat -tulpn | grep :8080
   ```

2. **Kill existing process**
   ```bash
   kill -9 <PID>
   ```

3. **Use different port**
   ```bash
   make run PORT=8081
   ```

#### Problem: "module not found" errors
```
Error: cannot find module github.com/helloellinor/p2k16/internal/models
```

**Solutions:**
```bash
# Download dependencies
go mod download

# Tidy up modules
go mod tidy

# Clear module cache if corrupted
go clean -modcache
go mod download
```

#### Problem: Build failures
```
Error: undefined: gin.Context
```

**Solutions:**
```bash
# Clean build cache
go clean -cache

# Rebuild from scratch
make clean
make build

# Check Go version (requires 1.21+)
go version
```

### Python Application Issues

#### Problem: Python virtual environment issues
```
Error: No module named 'flask'
```

**Solutions:**
1. **Activate virtual environment**
   ```bash
   # fish shell
   source env/bin/activate.fish
   
   # bash/zsh
   source env/bin/activate
   ```

2. **Install dependencies**
   ```bash
   pip install -r requirements.txt
   ```

3. **Python version compatibility (use Python 3.11)**
   ```bash
   # Check Python version
   python --version
   
   # Install Python 3.11 on macOS
   brew install python@3.11
   
   # Create new virtual environment
   python3.11 -m venv env
   source env/bin/activate.fish
   pip install --upgrade pip setuptools wheel
   pip install -r requirements.txt
   ```

#### Problem: "pkgutil.ImpImporter" error with Python 3.12+
```
AttributeError: module 'pkgutil' has no attribute 'ImpImporter'
```

**Solution: Downgrade to Python 3.11**
```bash
# Remove current environment
rm -rf env

# Install Python 3.11
brew install python@3.11  # macOS
# or use pyenv to manage versions

# Create new environment with Python 3.11
python3.11 -m venv env
source env/bin/activate.fish
pip install -r requirements.txt
```

#### Problem: Flask application won't start
```
Error: Address already in use
```

**Solutions:**
```bash
# Check what's using port 5000
lsof -i :5000

# Kill process
kill -9 <PID>

# Or run on different port
flask run --port 5001
```

### Environment Configuration Issues

#### Problem: Environment variables not loaded
**Solution: Check .env file**
```bash
# Copy example environment file
cp .env.example .env

# Edit with your settings
nano .env

# Verify environment loading
env | grep DB_
```

#### Problem: Fish shell configuration issues
**Solution: Set up fish environment**
```bash
# Install fish if needed
brew install fish  # macOS
sudo apt install fish  # Ubuntu

# Source settings
source .settings.fish

# Or manually set environment variables
set -x DB_HOST localhost
set -x DB_PORT 5432
set -x DB_USER p2k16
```

### Migration-Specific Issues

#### Problem: Session compatibility between Python and Go
**Symptoms:** Login works in one system but not the other

**Solutions:**
1. **Check session stores**
   ```python
   # Python session debugging
   print(f"Session data: {session}")
   print(f"Session ID: {session.sid}")
   ```

   ```go
   // Go session debugging  
   session := sessions.Default(c)
   log.Printf("Session data: %+v", session.Values)
   ```

2. **Verify cookie settings**
   ```python
   # Python session config
   app.config['SESSION_COOKIE_DOMAIN'] = 'localhost'
   app.config['SESSION_COOKIE_SECURE'] = False  # for development
   ```

   ```go
   // Go session config
   store.Options.Domain = "localhost"
   store.Options.Secure = false  // for development
   ```

#### Problem: API compatibility issues
**Symptoms:** Same request works in Python but fails in Go

**Solutions:**
1. **Compare request/response formats**
   ```bash
   # Test Python endpoint
   curl -X POST http://localhost:5000/api/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"test"}' \
     -v

   # Test Go endpoint  
   curl -X POST http://localhost:8081/api/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"test"}' \
     -v
   ```

2. **Check error logs**
   ```bash
   # Go logs
   make run 2>&1 | tee go.log

   # Python logs
   tail -f /var/log/p2k16/app.log
   ```

### HTMX Issues

#### Problem: HTMX requests not working
**Symptoms:** Forms submit but nothing happens

**Solutions:**
1. **Check HTMX library is loaded**
   ```html
   <script src="https://unpkg.com/htmx.org@1.8.4"></script>
   ```

2. **Verify HTMX attributes**
   ```html
   <form hx-post="/api/login" hx-target="#result">
     <!-- form fields -->
   </form>
   <div id="result"></div>
   ```

3. **Check browser console for errors**
   ```javascript
   // Browser console debugging
   htmx.logAll();
   ```

4. **Check Content Security Policy**
   ```html
   <meta http-equiv="Content-Security-Policy" 
         content="script-src 'self' 'unsafe-inline' unpkg.com;">
   ```

### Docker Issues

#### Problem: Docker containers won't start
**Solutions:**
```bash
# Check Docker is running
docker info

# Check container logs
docker logs p2k16-postgres

# Remove and recreate containers
docker rm -f p2k16-postgres
docker run --name p2k16-postgres \
  -e POSTGRES_USER=p2k16 \
  -e POSTGRES_PASSWORD=p2k16 \
  -e POSTGRES_DB=p2k16 \
  -p 5432:5432 \
  -d postgres:13
```

## Diagnostic Commands

### Database Diagnostics
```bash
# Check database connection
psql -h localhost -U p2k16 -d p2k16 -c "SELECT 1;"

# Check table structure
psql -h localhost -U p2k16 -d p2k16 -c "\d accounts"

# Check data
psql -h localhost -U p2k16 -d p2k16 -c "SELECT COUNT(*) FROM accounts;"

# Check active connections
psql -h localhost -U p2k16 -d p2k16 -c "SELECT COUNT(*) FROM pg_stat_activity;"
```

### Go Application Diagnostics
```bash
# Check Go installation
go version

# Check dependencies
go mod verify

# Check build
go build -v ./cmd/server

# Run tests
go test -v ./...

# Check for race conditions
go test -race ./...
```

### Python Application Diagnostics
```bash
# Check Python installation
python --version

# Check virtual environment
which python
which pip

# Check dependencies
pip list

# Check Flask app
python -c "from web.src.p2k16 import app; print('Flask app loads OK')"
```

### Network Diagnostics
```bash
# Check ports
netstat -tulpn | grep -E ":(5000|8080|8081|5432)"

# Test connectivity
curl -I http://localhost:5000/
curl -I http://localhost:8080/

# Check DNS resolution (if using domain names)
nslookup localhost
```

## Recovery Procedures

### Reset Development Environment
```bash
# 1. Stop all services
pkill -f python
pkill -f p2k16

# 2. Reset database
dropdb -h localhost -U p2k16 p2k16
createdb -h localhost -U p2k16 p2k16
psql -h localhost -U p2k16 -d p2k16 -f database-setup.sql

# 3. Clean Go build
make clean
go clean -cache -modcache

# 4. Rebuild Python environment
rm -rf env
python3.11 -m venv env
source env/bin/activate.fish
pip install -r requirements.txt

# 5. Start fresh
make dev-migration
```

### Emergency Rollback
If migration testing goes wrong:
```bash
# 1. Stop Go system
pkill -f p2k16-server

# 2. Ensure Python system is running
source .settings.fish
p2k16-run-web

# 3. Reset database to known good state
psql -h localhost -U p2k16 -d p2k16 < backup.sql

# 4. Clear any problematic sessions
redis-cli FLUSHDB  # if using Redis for sessions
```

## Getting Help

### Log Collection
Before asking for help, collect relevant logs:
```bash
# Go application logs
make run 2>&1 | tee logs/go-$(date +%Y%m%d-%H%M%S).log

# Python application logs
python -m web.main 2>&1 | tee logs/python-$(date +%Y%m%d-%H%M%S).log

# Database logs
tail -100 /var/log/postgresql/postgresql-13-main.log

# System logs
journalctl -u postgresql -n 50
```

### Environment Information
```bash
# System information
uname -a
cat /etc/os-release  # Linux
sw_vers             # macOS

# Development tools
go version
python --version
psql --version
docker --version

# Network configuration
ifconfig
netstat -rn
```

### Creating Minimal Reproduction
1. Start with clean environment
2. Document exact steps taken
3. Provide error messages in full
4. Include relevant configuration files
5. Share logs (sanitized of sensitive data)

---

**Last Updated**: [Current Date]  
**Troubleshooting Maintainer**: [To be assigned]