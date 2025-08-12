# P2K16 Go + HTMX Migration

This directory contains the new Go-based implementation of P2K16 using HTMX for frontend interactivity.

## Running the Go Server

1. **Install Go 1.21+** if not already installed

2. **Set up environment variables** (copy `.env.example` to `.env` and modify as needed):
   ```bash
   cp .env.example .env
   ```

3. **Build and run the application**:
   ```bash
   go build -o p2k16-server ./cmd/server
   ./p2k16-server
   ```

   Or run directly:
   ```bash
   go run ./cmd/server
   ```

4. **Access the application** at http://localhost:8080

## Project Structure

```
├── cmd/server/          # Main application entry point
├── internal/
│   ├── database/        # Database connection logic
│   ├── handlers/        # HTTP request handlers
│   ├── middleware/      # HTTP middleware (logging, CORS, etc.)
│   └── models/          # Data models and repositories
├── web/templates/       # HTML templates (future)
├── static/             # Static assets (CSS, JS, images)
└── config/             # Configuration files
```

## Features Implemented

- [x] Basic Go project structure
- [x] PostgreSQL database connection (compatible with existing schema)
- [x] HTMX-based frontend (replaces AngularJS)
- [x] Basic authentication framework
- [x] Home page with HTMX interactions
- [x] Login page with form handling
- [ ] Session management
- [ ] Badge system
- [ ] Door access control
- [ ] Tool management
- [ ] Membership management

## Migration Strategy

This Go application is designed to coexist with the existing Flask application during migration:

1. **Database compatibility**: Uses the same PostgreSQL schema as the Flask app
2. **Incremental migration**: Features can be migrated one by one
3. **API compatibility**: Maintains similar API endpoints where possible

## Technology Stack

- **Backend**: Go with Gin web framework
- **Frontend**: HTMX for dynamic interactions, Bootstrap for styling
- **Database**: PostgreSQL (existing schema)
- **Authentication**: bcrypt password hashing (compatible with Flask-Bcrypt)