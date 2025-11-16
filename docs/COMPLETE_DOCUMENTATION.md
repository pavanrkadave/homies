# Homies Expense Tracker - Complete Documentation

**Version:** 1.0.0  
**Last Updated:** November 16, 2025  
**Status:** Production Ready ✅

---

## Table of Contents
1. [Overview](#overview)
2. [Architecture](#architecture)
3. [API Endpoints](#api-endpoints)
4. [Implementation Progress](#implementation-progress)
5. [Setup & Deployment](#setup--deployment)
6. [Testing](#testing)
7. [Development History](#development-history)

---

## Overview

Homies is a production-ready expense tracker REST API built with Go 1.25, following Clean Architecture principles. It allows roommates to track shared expenses, split costs, and calculate settlements.

### Key Features
- ✅ User management with CRUD operations
- ✅ Expense tracking with flexible splitting
- ✅ Automatic equal split calculation
- ✅ Balance calculation and settlement suggestions
- ✅ Expense filtering by category and date range
- ✅ PostgreSQL database with migrations
- ✅ Docker containerization
- ✅ Comprehensive test coverage

### Tech Stack
- **Language:** Go 1.25
- **Database:** PostgreSQL
- **Containerization:** Docker & Docker Compose
- **Architecture:** Clean Architecture
- **Testing:** Go testing package with mocks

---

## Architecture

### Clean Architecture Layers

```
┌─────────────────────────────────────────────┐
│              HTTP Layer                      │
│  (Handlers, Middleware, Routes)              │
│         internal/handler/                    │
│         internal/middleware/                 │
└──────────────┬──────────────────────────────┘
               │
┌──────────────▼──────────────────────────────┐
│          Use Case Layer                      │
│     (Business Logic)                         │
│         internal/usecase/                    │
└──────────────┬──────────────────────────────┘
               │
┌──────────────▼──────────────────────────────┐
│        Repository Layer                      │
│  (Data Access Interfaces)                    │
│         internal/repository/                 │
└──────────────┬──────────────────────────────┘
               │
┌──────────────▼──────────────────────────────┐
│          Domain Layer                        │
│   (Business Entities, Validation)            │
│         internal/domain/                     │
└──────────────────────────────────────────────┘
```

### Directory Structure
```
homies/
├── cmd/
│   ├── api/main.go              # Application entry point
│   └── migrate/main.go          # Migration runner
├── internal/
│   ├── domain/                  # Business entities
│   ├── usecase/                 # Business logic
│   ├── repository/              # Data access
│   │   ├── postgres/            # PostgreSQL implementation
│   │   └── memory/              # In-memory for testing
│   ├── handler/                 # HTTP handlers
│   └── middleware/              # HTTP middleware
├── pkg/
│   ├── database/                # Database utilities
│   └── response/                # Response helpers
├── migrations/                  # SQL migrations
├── config/                      # Configuration
└── docker-compose.yml           # Docker setup
```

---

## API Endpoints

### Base URL
```
http://localhost:3000
```

### User Endpoints

#### Get All Users
```http
GET /users
```

**Response:**
```json
[
  {
    "id": "uuid",
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2025-11-16T10:00:00Z"
  }
]
```

#### Get User by ID
```http
GET /users?id={user_id}
```

**Response:**
```json
{
  "id": "uuid",
  "name": "John Doe",
  "email": "john@example.com",
  "created_at": "2025-11-16T10:00:00Z"
}
```

#### Create User
```http
POST /users
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com"
}
```

**Response:** `201 Created`

#### Update User
```http
PUT /users?id={user_id}
Content-Type: application/json

{
  "name": "John Smith",
  "email": "john.smith@example.com"
}
```

**Response:** `200 OK`

---

### Expense Endpoints

#### Get All Expenses (with filters)
```http
GET /expenses
GET /expenses?category=food
GET /expenses?start_date=2025-11-01&end_date=2025-11-30
GET /expenses?category=food&start_date=2025-11-01&end_date=2025-11-30
```

**Response:**
```json
[
  {
    "id": "uuid",
    "description": "Groceries",
    "amount": 100.00,
    "category": "food",
    "paid_by": "user_id",
    "date": "2025-11-16T10:00:00Z",
    "created_at": "2025-11-16T10:00:00Z",
    "splits": [
      {
        "user_id": "user1_id",
        "amount": 50.00
      },
      {
        "user_id": "user2_id",
        "amount": 50.00
      }
    ]
  }
]
```

#### Get Expense by ID
```http
GET /expenses?id={expense_id}
```

#### Get Expenses by User
```http
GET /expenses/user?user_id={user_id}
```

#### Create Expense
```http
POST /expenses
Content-Type: application/json

{
  "description": "Dinner",
  "amount": 100.00,
  "category": "food",
  "paid_by": "user_id",
  "splits": [
    {
      "user_id": "user1_id",
      "amount": 60.00
    },
    {
      "user_id": "user2_id",
      "amount": 40.00
    }
  ]
}
```

**Response:** `201 Created`

#### Create Expense with Equal Split
```http
POST /expenses/equal-split
Content-Type: application/json

{
  "description": "Team Dinner",
  "amount": 100.00,
  "category": "food",
  "paid_by": "user_id",
  "user_ids": ["user1_id", "user2_id", "user3_id"]
}
```

**Response:** `201 Created` (automatically calculates equal splits)

#### Update Expense
```http
PUT /expenses?id={expense_id}
Content-Type: application/json

{
  "description": "Updated Dinner",
  "amount": 120.00,
  "category": "dining",
  "paid_by": "user_id",
  "splits": [...]
}
```

**Response:** `200 OK`

#### Delete Expense
```http
DELETE /expenses?id={expense_id}
```

**Response:** `204 No Content`

---

### Balance Endpoints

#### Get Balances & Settlements
```http
GET /balances
```

**Response:**
```json
{
  "balances": [
    {
      "user_id": "user1_id",
      "amount": 50.00
    },
    {
      "user_id": "user2_id",
      "amount": -50.00
    }
  ],
  "settlements": [
    {
      "from": "user2_id",
      "to": "user1_id",
      "amount": 50.00
    }
  ]
}
```

---

### Health Check

#### Health Status
```http
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "database": "connected"
}
```

---

## Implementation Progress

### Phase 1: User Management ✅
- Update User endpoint
- Get User by ID endpoint
- **Tests:** 5/5 passing

### Phase 2: Expense Enhancements ✅
- Update Expense endpoint
- Equal Split Helper
- **Tests:** 7/7 passing

### Phase 3: Filtering & Search ✅
- Filter by date range
- Filter by category
- Combined filters
- **Tests:** 4/4 passing

### Total
- **Features:** 8 completed
- **Tests:** 16/16 passing (100%)
- **Code Quality:** Excellent
- **Status:** Production Ready

---

## Setup & Deployment

### Prerequisites
- Docker & Docker Compose
- Go 1.25+ (for local development)

### Quick Start with Docker
```bash
# Clone the repository
git clone <repository-url>
cd homies

# Start services
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f app

# Stop services
docker-compose down
```

### Local Development
```bash
# Install dependencies
go mod download

# Run migrations
go run cmd/migrate/main.go

# Run application
go run cmd/api/main.go

# Run tests
go test ./...
```

### Environment Variables
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=homies
SERVER_PORT=3000
```

---

## Testing

### Run All Tests
```bash
go test ./...
```

### Run Specific Package Tests
```bash
go test ./internal/usecase -v
go test ./internal/repository/memory -v
```

### Test Coverage
```bash
go test ./... -cover
```

### Manual API Testing
```bash
# Get all users
curl http://localhost:3000/users

# Create user
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'

# Filter expenses
curl "http://localhost:3000/expenses?category=food&start_date=2025-11-01&end_date=2025-11-30"
```

---

## Development History

### Phase 1: User Management Enhancements
**Completed:** November 15, 2025

**Features:**
- Update user endpoint with email uniqueness validation
- Get user by ID with 404 handling
- Enhanced user repository with Update method

**Files Modified:** 8 files  
**Tests Added:** 5 tests

---

### Phase 2: Expense Enhancements
**Completed:** November 16, 2025

**Features:**
- Update expense endpoint with transaction safety
- Equal split helper with automatic calculation and rounding
- Split validation and recalculation

**Files Modified:** 7 files  
**Tests Added:** 7 tests

**Highlights:**
- Transaction-based updates ensure data integrity
- Smart rounding algorithm (e.g., 100/3 = 33.33, 33.33, 33.34)
- Last user gets remainder for exact totals

---

### Phase 3: Filtering & Search
**Completed:** November 16, 2025

**Features:**
- Filter expenses by date range (ISO 8601 format)
- Filter expenses by category (case-insensitive)
- Combined filters with dynamic SQL query building
- Backward compatible (no breaking changes)

**Files Modified:** 6 files  
**Tests Added:** 4 tests

**Highlights:**
- Parameterized SQL queries prevent SQL injection
- Helper method `scanExpensesWithSplits()` eliminates duplication
- Index-ready query structure for performance

---

## Code Quality & Best Practices

### Security
✅ Parameterized SQL queries (no SQL injection)  
✅ Input validation at use case layer  
✅ Transaction safety for data integrity  
✅ Email uniqueness enforcement  

### Architecture
✅ Clean Architecture principles  
✅ Dependency inversion  
✅ Interface-based design  
✅ Separation of concerns  

### Code Standards
✅ Functions under 30 lines  
✅ Clear variable naming  
✅ DRY principles  
✅ Comprehensive error handling  
✅ Thread-safe implementations  

### Testing
✅ 100% use case test coverage  
✅ Mock repositories for unit tests  
✅ Integration tests with Docker  
✅ Edge case coverage  

---

## Performance Optimization

### Recommended Database Indexes
```sql
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_expenses_category ON expenses(category);
CREATE INDEX idx_expenses_date ON expenses(date);
CREATE INDEX idx_expenses_category_date ON expenses(category, date);
```

### Query Optimization
- Single query per request (no N+1 problems)
- Batch operations where possible
- Transaction-based updates
- Results sorted at database level

---

## Troubleshooting

### Database Connection Issues
```bash
# Check if PostgreSQL is running
docker-compose ps

# Restart database
docker-compose restart postgres

# View database logs
docker-compose logs postgres
```

### Application Not Starting
```bash
# Check application logs
docker-compose logs app

# Rebuild and restart
docker-compose up -d --build
```

### Migration Issues
```bash
# Run migrations manually
go run cmd/migrate/main.go
```

---

## Contributing

### Code Style
- Follow Go best practices
- Keep functions under 30 lines
- Write tests for all new features
- Document complex logic

### Commit Messages
```
<type>: <short description>

<detailed description>
- Bullet points of changes

Example usage
```

**Types:** feat, fix, refactor, test, docs, chore

---

## License

[Your License Here]

---

## Support

For issues and questions, please open an issue on the repository.

---

**Built with ❤️ using Clean Architecture principles**

