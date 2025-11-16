# Homies Expense Tracker - Cleanup & Improvements Plan

**Date:** November 16, 2025  
**Version:** 1.1.0

---

## Summary of Improvements

This document outlines the cleanup and improvements made to the Homies project before implementing Phase 4.

---

## 1. Documentation Consolidation ✅

### Problem
- Multiple scattered documentation files (PHASE1_COMPLETE.md, PHASE2_COMPLETE.md, PHASE2_SUMMARY.md, etc.)
- Redundant information across files
- Difficult to find specific information

### Solution
Created centralized documentation structure:

```
docs/
└── COMPLETE_DOCUMENTATION.md  # Single source of truth
```

**COMPLETE_DOCUMENTATION.md** includes:
- Overview and features
- Architecture explanation
- Complete API documentation
- Setup and deployment guide
- Testing instructions
- Development history
- Troubleshooting guide

### Files to Archive/Remove
Move to `docs/archive/`:
- PHASE1_COMPLETE.md
- PHASE2_COMPLETE.md
- PHASE2_SUMMARY.md
- PHASE3_COMPLETE.md
- PHASE3_SUMMARY.md
- IMPLEMENTATION_SUMMARY.md

Keep in root:
- README.md (quick start guide)
- PROJECT_STATUS.md (current status dashboard)
- QUICK_REFERENCE.md (developer quick reference)

---

## 2. Structured Logging ✅

### Library Choice: Zap by Uber

**Why Zap?**
- ⚡ Blazing fast (zero-allocation)
- 🏗️ Structured logging with fields
- 📊 Multiple log levels (debug, info, warn, error, fatal)
- 🎯 Production-ready
- 🔧 Highly configurable

### Implementation

**New Package:** `pkg/logger/logger.go`

**Features:**
- Level-based logging (debug, info, warn, error, fatal)
- Environment-based configuration (LOG_LEVEL env var)
- Production mode with JSON encoding
- Development mode with colored console output
- ISO 8601 timestamps
- Caller information

**Usage:**
```go
import "github.com/pavanrkadave/homies/pkg/logger"
import "go.uber.org/zap"

// Initialize logger
logger.InitLogger("info")  // or "debug", "warn", "error"

// Log messages
logger.Info("Server starting", zap.Int("port", 3000))
logger.Debug("Processing request", zap.String("method", "GET"))
logger.Warn("Slow query detected", zap.Duration("duration", time.Second))
logger.Error("Failed to connect", zap.Error(err))
```

### Environment Variables
```env
LOG_LEVEL=info  # debug, info, warn, error
```

---

## 3. Database Migrations - golang-migrate ✅

### Library Choice: golang-migrate

**Why golang-migrate?**
- 📦 Industry standard
- 🔄 Up and down migrations
- 🎯 Version tracking
- 🛡️ Dirty state detection
- 🔧 CLI and library support
- 📊 Multiple database support

### Features

**New Package:** `pkg/database/migrate_new.go`

**Capabilities:**
- ✅ Run migrations (Up)
- ✅ Rollback migrations (Down)
- ✅ Migrate to specific version
- ✅ Version tracking
- ✅ Dirty state detection
- ✅ Structured logging integration

**Migration Functions:**
```go
// Run all pending migrations
RunMigrations(db, "migrations/")

// Rollback last migration
RollbackMigration(db, "migrations/")

// Migrate to specific version
MigrateToVersion(db, "migrations/", 2)
```

### Migration File Naming Convention

**Format:** `{version}_{description}.{up|down}.sql`

```
migrations/
├── 001_create_users_table.up.sql
├── 001_create_users_table.down.sql
├── 002_create_expenses_table.up.sql
├── 002_create_expenses_table.down.sql
├── 003_create_splits_table.up.sql
└── 003_create_splits_table.down.sql
```

### Benefits
- ✅ Track migration history
- ✅ Easy rollback capability
- ✅ Prevent dirty states
- ✅ Team collaboration friendly
- ✅ CI/CD integration ready

---

## 4. OpenAPI/Swagger Documentation ✅

### Library Choice: swaggo/swag

**Why Swag?**
- 📝 Generate OpenAPI 3.0 spec from Go annotations
- 🌐 Interactive Swagger UI
- 🔄 Auto-updates from code
- 🎯 Type-safe
- 📊 Widely adopted

### Implementation

**Installation:**
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

**Generate Docs:**
```bash
swag init -g cmd/api/main.go -o docs/swagger
```

**Swagger UI Endpoint:**
```
http://localhost:3000/swagger/index.html
```

### Annotation Example

```go
// @Summary Create a new user
// @Description Create a new user with name and email
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User details"
// @Success 201 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Router /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    // ...
}
```

### Main API Documentation
```go
// @title Homies Expense Tracker API
// @version 1.0
// @description REST API for tracking shared expenses among roommates
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@homies.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3000
// @BasePath /
// @schemes http https
```

---

## 5. Code Organization Improvements

### Project Structure (Updated)

```
homies/
├── cmd/
│   ├── api/main.go              # Application entry (with Swagger)
│   └── migrate/main.go          # Migration CLI
├── docs/
│   ├── COMPLETE_DOCUMENTATION.md  # Main docs
│   ├── swagger/                   # Auto-generated
│   └── archive/                   # Old phase docs
├── internal/
│   ├── domain/
│   ├── usecase/
│   ├── repository/
│   ├── handler/                   # With Swagger annotations
│   └── middleware/
│       └── logging.go             # Request logging middleware
├── pkg/
│   ├── logger/                    # Structured logging
│   ├── database/
│   │   ├── postgres.go
│   │   ├── migrate.go (old)
│   │   └── migrate_new.go (new)
│   └── response/
├── migrations/                    # Renamed files
│   ├── 001_create_users_table.up.sql
│   ├── 001_create_users_table.down.sql
│   ├── 002_create_expenses_table.up.sql
│   ├── 002_create_expenses_table.down.sql
│   ├── 003_create_splits_table.up.sql
│   └── 003_create_splits_table.down.sql
├── .env.example                   # Environment template
├── Makefile                       # Common commands
├── README.md                      # Quick start
├── PROJECT_STATUS.md             # Current state
└── docker-compose.yml
```

---

## 6. Configuration Improvements

### New .env.example
```env
# Server
SERVER_PORT=3000
LOG_LEVEL=info

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=homies
DB_SSL_MODE=disable

# Application
ENV=development  # development, production
```

### Config Package Enhancement
```go
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Logger   LoggerConfig
}

type LoggerConfig struct {
    Level string
    Mode  string  // "development" or "production"
}
```

---

## 7. Makefile for Common Tasks

```makefile
.PHONY: help build run test clean docker-up docker-down migrate-up migrate-down swagger

help:
	@echo "Available commands:"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the application"
	@echo "  make test         - Run tests"
	@echo "  make docker-up    - Start Docker containers"
	@echo "  make docker-down  - Stop Docker containers"
	@echo "  make migrate-up   - Run migrations"
	@echo "  make migrate-down - Rollback last migration"
	@echo "  make swagger      - Generate Swagger docs"

build:
	go build -o bin/homies cmd/api/main.go

run:
	go run cmd/api/main.go

test:
	go test ./... -v

clean:
	rm -rf bin/

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

migrate-up:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down

swagger:
	swag init -g cmd/api/main.go -o docs/swagger
```

---

## 8. Middleware Enhancements

### Request Logging Middleware
```go
// internal/middleware/logging.go
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        logger.Info("Incoming request",
            zap.String("method", r.Method),
            zap.String("path", r.URL.Path),
            zap.String("remote_addr", r.RemoteAddr),
        )
        
        next.ServeHTTP(w, r)
        
        logger.Info("Request completed",
            zap.String("method", r.Method),
            zap.String("path", r.URL.Path),
            zap.Duration("duration", time.Since(start)),
        )
    })
}
```

---

## 9. Testing Improvements

### Add Test Logging
```go
func TestMain(m *testing.M) {
    // Initialize logger for tests
    logger.InitLogger("error")  // Only show errors in tests
    code := m.Run()
    logger.Sync()
    os.Exit(code)
}
```

---

## 10. Docker Improvements

### Updated Dockerfile
```dockerfile
# Use specific Go version
FROM golang:1.25-alpine AS builder

# Install dependencies including swag
RUN apk add --no-cache git
RUN go install github.com/swaggo/swag/cmd/swag@latest

# ... rest of build steps ...

# Generate Swagger docs
RUN swag init -g cmd/api/main.go -o docs/swagger
```

---

## Implementation Checklist

### Phase 1: Documentation ✅
- [x] Create docs/ directory
- [x] Consolidate all docs into COMPLETE_DOCUMENTATION.md
- [ ] Move old phase docs to docs/archive/
- [ ] Update README.md to reference new structure

### Phase 2: Logging ✅
- [x] Install Zap
- [x] Create pkg/logger package
- [ ] Update main.go to initialize logger
- [ ] Replace all log.Println with structured logging
- [ ] Add logging middleware
- [ ] Update config to include LOG_LEVEL

### Phase 3: Migrations ✅
- [x] Install golang-migrate
- [x] Create migrate_new.go with new functions
- [ ] Rename migration files to new format
- [ ] Update cmd/migrate/main.go
- [ ] Test up/down migrations
- [ ] Document migration commands

### Phase 4: Swagger ✅
- [x] Install swaggo/swag
- [ ] Add main API annotations
- [ ] Add annotations to all handlers
- [ ] Generate Swagger docs
- [ ] Add Swagger UI route
- [ ] Test Swagger endpoint

### Phase 5: Code Cleanup
- [ ] Create Makefile
- [ ] Create .env.example
- [ ] Update config package
- [ ] Add request logging middleware
- [ ] Update Dockerfile
- [ ] Update docker-compose.yml
- [ ] Run go mod tidy

### Phase 6: Testing
- [ ] Run all tests
- [ ] Update test documentation
- [ ] Add integration test examples

---

## Next Steps After Cleanup

1. **Test Everything**
   - All unit tests pass
   - Docker containers start correctly
   - Migrations work (up and down)
   - Swagger UI loads
   - All endpoints documented

2. **Update Documentation**
   - README.md with new structure
   - Update QUICK_REFERENCE.md
   - Add migration guide
   - Add logging guide

3. **Ready for Phase 4**
   - Statistics & Reporting
   - User spending summary
   - Monthly summaries

---

## Benefits of These Improvements

### For Developers
- ✅ Centralized documentation
- ✅ Better debugging with structured logs
- ✅ Easy migration management
- ✅ Interactive API documentation
- ✅ Clear project structure

### For Operations
- ✅ Production-ready logging
- ✅ Database version control
- ✅ Easy rollback capability
- ✅ Health monitoring ready
- ✅ CI/CD friendly

### For Users
- ✅ Self-documenting API
- ✅ Better error messages
- ✅ Faster debugging
- ✅ More reliable system

---

**Status:** Ready for Implementation  
**Estimated Time:** 2-3 hours  
**Priority:** High (before Phase 4)

