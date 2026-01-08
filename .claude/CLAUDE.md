# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`init-database` is a Go CLI tool that initializes Senzing databases. It creates the database schema from SQL files and sets up the default Senzing configuration. Part of the `senzing-tools` suite.

## Common Commands

```bash
# Build
make build                      # Build binary to target/{os}-{arch}/init-database

# Test (requires Senzing C library at /opt/senzing/er/lib)
make setup                      # One-time setup: creates /tmp/sqlite/G2C.db
make test                       # Run all tests with SQLite

# Test with specific databases
make test-postgresql            # Test with PostgreSQL (requires running instance)
make test-mysql                 # Test with MySQL
make test-mssql                 # Test with MS SQL

# Run single test
go test -tags "libsqlite3 linux" -v -run TestFunctionName ./package/...

# Lint
make lint                       # Runs golangci-lint, govulncheck, cspell
make fix                        # Auto-fix lint issues (gofumpt, wsl, etc.)

# Coverage
make coverage                   # Generate and open coverage report in browser

# Clean
make clean                      # Remove build artifacts, test databases, docker containers
```

## Architecture

The tool follows a layered initialization pattern:

```console
cmd/root.go          → CLI entry (Cobra), parses args, creates BasicInitializer
    ↓
initializer/         → Orchestrates initialization sequence
    ↓
    ├── senzingschema/   → Creates database schema from SQL files
    ├── senzingconfig/   → Installs Senzing configuration (datasources, etc.)
    └── senzingload/     → Optionally loads TruthSet data
```

Each package defines an interface (`Initializer`, `SenzingSchema`, `SenzingConfig`, `SenzingLoad`) with a `Basic*` implementation. All implementations support the observer pattern for event notification.

**Platform-specific code**: Files suffixed with `_linux.go`, `_darwin.go`, `_windows.go` contain OS-specific implementations (e.g., `initializer_linux.go` handles filename cleaning).

## Key Patterns

- **Error wrapping**: Always use `wraperror.Errorf(err, "context")` for errors
- **Logging**: Each package has a `ComponentID` constant (6501-6504) and `IDMessages` map for structured logging
- **Configuration**: Environment variables use `SENZING_TOOLS_` prefix (e.g., `SENZING_TOOLS_DATABASE_URL`)
- **Singletons**: `Basic*` structs use mutex-protected lazy initialization for dependent services

## Environment Variables

- `SENZING_TOOLS_DATABASE_URL` - Database connection string (e.g., `sqlite3://na:na@nowhere/tmp/sqlite/G2C.db`)
- `SENZING_TOOLS_ENGINE_CONFIGURATION_JSON` - Senzing engine settings JSON
- `SENZING_TOOLS_SQL_FILE` - Path to SQL schema file (auto-detected from database type)
- `LD_LIBRARY_PATH` - Must include `/opt/senzing/er/lib` for tests

## Linting Configuration

Strict golangci-lint config at `.github/linters/.golangci.yaml`. Key settings:

- Max line length: 120 (golines)
- Max function complexity: 25 (cyclop)
- Max cognitive complexity: 40 (gocognit)
- Max function length: 150 lines / 80 statements
