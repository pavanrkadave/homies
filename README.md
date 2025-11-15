# Homies - Expense Tracker API

A production-ready expense tracker REST API built with Go, following Clean Architecture principles.

## Features

- 👥 User management (Create, List)
- 💰 Expense tracking with flexible splitting
- ✅ Automatic validation (splits must sum to total)
- 🏗️ Clean Architecture (Domain, Use Case, Repository, Handler)
- 🧪 Comprehensive unit tests
- 🔒 Type-safe with Go's strong typing

## Architecture
```
cmd/api/              # Application entry point
internal/
  ├── domain/         # Business entities & validation
  ├── usecase/        # Business logic
  ├── repository/     # Data access layer
  │   └── memory/     # In-memory implementation
  └── handler/        # HTTP handlers
```

## Getting Started

### Prerequisites
- Go 1.25+

### Installation
```bash
# Clone the repository
git clone https://github.com/pavanrkadave/homies.git
cd homies

# Install dependencies
go mod download

# Run the server
go run cmd/api/main.go
```

Server starts on `http://localhost:3000`

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -cover

# Run tests verbosely
go test ./... -v
```

## API Endpoints

### Users

**Create User**
```bash
POST /users
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com"
}
```

**Get All Users**
```bash
GET /users
```

### Expenses

**Create Expense**
```bash
POST /expenses
Content-Type: application/json

{
  "description": "Groceries",
  "amount": 100.00,
  "category": "groceries",
  "paid_by": "user-id",
  "splits": [
    {"user_id": "user-id-1", "amount": 50.00},
    {"user_id": "user-id-2", "amount": 50.00}
  ]
}
```

**Get All Expenses**
```bash
GET /expenses
```

## Project Status

✅ Core Features Implemented
- User CRUD operations
- Expense creation with validation
- Split tracking

🚧 Coming Soon
- Database persistence (PostgreSQL)
- Balance calculations
- Settlement tracking
- Authentication & Authorization

## Tech Stack

- **Language**: Go 1.25
- **Architecture**: Clean Architecture
- **Testing**: Go's built-in testing package
- **Storage**: In-memory (PostgreSQL coming soon)

## License

MIT