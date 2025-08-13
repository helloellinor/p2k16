# P2K16 Local Development Setup

This page documents the exact, tested steps to get P2K16 running locally on macOS (fish shell) and Linux. It covers Python, Docker/Postgres, Flyway, and both the Python web app and the new Go server.

---

## TL;DR (macOS + fish)

```fish
# Prereqs
brew install colima docker libpq openjdk python@3.11
colima start

# Shell env
echo 'set -gx PATH /opt/homebrew/opt/libpq/bin $PATH' >> ~/.config/fish/config.fish
echo 'set -gx JAVA_HOME (brew --prefix openjdk)/libexec/openjdk.jdk/Contents/Home' >> ~/.config/fish/config.fish
echo 'set -gx PATH $JAVA_HOME/bin $PATH' >> ~/.config/fish/config.fish
source ~/.config/fish/config.fish

# Python env (3.11)
python3.11 -m venv env
source env/bin/activate.fish
pip install --upgrade pip setuptools wheel cython packaging
pip install -r requirements.txt

# Project PATH/env
source .settings.fish

# Start Postgres
cd docker/p2k16; docker-compose up -d; cd -

# Seed roles/db (password: postgres)
psql -h localhost -p 2016 -U postgres -f database-setup.sql

# Apply DB migrations (alternative: p2k16-run-web runs this too)
bin/flyway -url=jdbc:postgresql://localhost:2016/p2k16 -user=postgres -password=postgres -locations=filesystem:$(pwd)/migrations migrate

# Start Python web app (Gunicorn on :5000)
p2k16-run-web

# Start Go server (Gin on :8081)
make run PORT=8081
```

---

## 1) Prerequisites

- Docker backend
  - macOS: Docker Desktop or Colima (recommended). For Colima: `brew install colima docker && colima start`.
- PostgreSQL client tools
  - macOS: `brew install libpq` and add to PATH (below).
- Java runtime for Flyway
  - macOS: `brew install openjdk`; set JAVA_HOME (below).
- Python 3.11
  - macOS: `brew install python@3.11`.
- Go toolchain (for Go server)
  - macOS: `brew install go` or from golang.org.

### macOS PATH and env

```fish
echo 'set -gx PATH /opt/homebrew/opt/libpq/bin $PATH' >> ~/.config/fish/config.fish
echo 'set -gx JAVA_HOME (brew --prefix openjdk)/libexec/openjdk.jdk/Contents/Home' >> ~/.config/fish/config.fish
echo 'set -gx PATH $JAVA_HOME/bin $PATH' >> ~/.config/fish/config.fish
source ~/.config/fish/config.fish
```

---

## 2) Python environment (3.11)

Python 3.11 is recommended; 3.12+ breaks some legacy deps.

```fish
python3.11 -m venv env
source env/bin/activate.fish
pip install --upgrade pip setuptools wheel cython packaging
pip install -r requirements.txt
```

If editable install is needed for the Python web package:

```fish
pip install -e web
```

Common build errors (matplotlib, setuptools/packaging, cython) are resolved by the upgrades above. See Troubleshooting below for more.

---

## 3) Project environment

The repo includes helper scripts in `bin/`. Source the local settings to put them on PATH and set PG defaults:

```fish
source .settings.fish
```

---

## 4) Database: Docker Postgres

Start the database:

```fish
cd docker/p2k16
docker-compose up -d
cd -
```

The container exposes Postgres on localhost:2016 (mapped from 5432).

Initialize roles and the database:

```fish
psql -h localhost -p 2016 -U postgres -f database-setup.sql
```

Run Flyway migrations (idempotent):

```fish
bin/flyway -url=jdbc:postgresql://localhost:2016/p2k16 \
  -user=postgres -password=postgres \
  -locations=filesystem:(pwd)/migrations migrate
```

Note: The `p2k16-run-web` script also performs migrations automatically using explicit flags.

---

## 5) Run the apps

### Python web app (Flask/Gunicorn)

```fish
p2k16-run-web
```

This will: install Python deps, apply migrations, install bower assets, and start Gunicorn on http://localhost:5000.

### Go server (Gin)

The Go server can be run via Makefile. It now reads DB settings from environment and defaults DB_PORT to 2016.

```fish
make run PORT=8081
```

Environment variables used by the server (with defaults):

- DB_HOST=localhost
- DB_PORT=2016
- DB_USER=p2k16-web
- DB_PASSWORD=p2k16-web
- DB_NAME=p2k16
- PORT=8080 (override with `PORT=8081` if 8080 is busy)

---

## 6) Tests

Use the bundled runner (also applies migrations):

```fish
bin/p2k16-run-test
```

Or run Go tests:

```fish
go test ./...
```

---

## Troubleshooting

- pq: column "created_by" does not exist (or similar errors when logging in with Go server)
  - Your database schema is missing a column expected by the Go backend. To fix:
    ```sql
    ALTER TABLE account ADD COLUMN created_by integer;
    ALTER TABLE account ADD COLUMN updated_by integer;
    ```
  - After running these, restart the Go server and try logging in again.

- role "p2k16-web" does not exist
  - You are likely connecting to port 5432 instead of 2016. Set `DB_PORT=2016` or use the Makefile `run` target. The Go server now defaults to 2016.

- ERROR: flyway.url must be set / Unable to connect to the database
  - Use explicit flags to Flyway as shown above, or ensure your config points to `jdbc:postgresql://localhost:2016/p2k16` with user/password `postgres`.

- Python build errors (cython_sources, pkgutil.ImpImporter, canonicalize_version)
  - Upgrade build tooling in your venv: `pip install --upgrade pip setuptools wheel cython packaging` and use Python 3.11.

- Port 8080 already in use for Go server
  - Run with `make run PORT=8081`.

- Permission denied during migration or cached plan errors
  - Stop containers, clear local data dir under `docker/p2k16/pdb/`, start fresh, and re-run migrations.

---

## Reference

- DB connection for local dev: `localhost:2016`, db `p2k16`, roles created by `database-setup.sql`.
- Flyway CLI is vendored via `bin/flyway` and auto-downloads v7.8.1.
- Python app config: `infrastructure/config-local.cfg`.
- Makefile `run` exports DB_* and PORT envs to the Go server process.
