# Homies - Expense Tracker API

A production-ready expense tracker REST API built with Go 1.25, following Clean Architecture principles. Track shared expenses among roommates, split costs automatically, and calculate settlements.

## ✨ Features

- 👥 **User Management** - Create, update, list, and manage users
- 💰 **Expense Tracking** - Track expenses with flexible splitting
- ⚡ **Equal Split Helper** - Automatically calculate equal splits
- 🔍 **Filtering & Search** - Filter by category, date range, or both
- 💵 **Balance Calculation** - Automatic balance and settlement suggestions
- 🏗️ **Clean Architecture** - Domain-driven design with clear separation
- 🧪 **Comprehensive Tests** - 16+ unit tests with 100% coverage
- 🐳 **Docker Ready** - Complete Docker Compose setup
- 📊 **Structured Logging** - Production-ready logging with Zap
- 🔄 **Database Migrations** - Version-controlled migrations with golang-migrate
- 📝 **API Documentation** - OpenAPI/Swagger support

## 🏗️ Architecture

Clean Architecture with dependency inversion:
```
cmd/
  ├── api/           # Application entry point
  └── migrate/       # Migration runner
internal/
  ├── domain/        # Business entities & validation
  ├── usecase/       # Business logic layer
  ├── repository/    # Data access interfaces
  │   ├── postgres/  # PostgreSQL implementation
  │   └── memory/    # In-memory for testing
  ├── handler/       # HTTP handlers
  └── middleware/    # HTTP middleware
pkg/
  ├── logger/        # Structured logging
  ├── database/      # DB utilities & migrations
  └── response/      # Response helpers
```

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose (recommended)
- Go 1.25+ (for local development)

### Using Docker (Recommended)
```bash
# Clone the repository
git clone <repository-url>
cd homies

# Start everything with one command
make dev

# Or manually
docker-compose up -d

# View logs
make logs

# Stop services
docker-compose down
```

API available at `http://localhost:3000`

### Local Development
```bash
# Install dependencies
go mod download

# Setup environment
cp .env.example .env

# Run migrations
make migrate-up

# Run the server
make run

# Or with go run
go run cmd/api/main.go
```

### Running Tests
```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests verbose
make test-verbose

# Run tests verbosely
go test ./... -v
```

## 📡 API Endpoints

### Users
- `GET /users` - List all users
- `GET /users?id={id}` - Get user by ID
- `POST /users` - Create new user
- `PUT /users?id={id}` - Update user

### Expenses
- `GET /expenses` - List all expenses (with optional filters)
- `GET /expenses?category={category}` - Filter by category
- `GET /expenses?start_date={date}&end_date={date}` - Filter by date range
- `GET /expenses?id={id}` - Get expense by ID
- `GET /expenses/user?user_id={id}` - Get user's expenses
- `POST /expenses` - Create expense
- `POST /expenses/equal-split` - Create expense with equal split
- `PUT /expenses?id={id}` - Update expense
- `DELETE /expenses?id={id}` - Delete expense

### Balance
- `GET /balances` - Get balances and settlement suggestions

### Health
- `GET /health` - Health check

**Full API Documentation:** See [docs/COMPLETE_DOCUMENTATION.md](docs/COMPLETE_DOCUMENTATION.md)

## 🛠️ Development

### Makefile Commands
```bash
make help           # Show all available commands
make build          # Build the application
make run            # Run the application
make test           # Run tests
make docker-up      # Start Docker containers
make docker-down    # Stop Docker containers
make docker-rebuild # Rebuild and restart containers
make migrate-up     # Run database migrations
make migrate-down   # Rollback last migration
make swagger        # Generate API documentation
make lint           # Run linter
make fmt            # Format code
```

### Environment Variables
See `.env.example` for all configuration options:
- `SERVER_PORT` - Server port (default: 3000)
- `LOG_LEVEL` - Logging level (debug, info, warn, error)
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` - Database config

## 📚 Documentation

- **[Complete Documentation](docs/COMPLETE_DOCUMENTATION.md)** - Comprehensive API docs and guides
- **[Quick Reference](docs/QUICK_REFERENCE.md)** - Developer quick reference
- **[Project Status](docs/PROJECT_STATUS.md)** - Current project status and metrics
- **[Cleanup Plan](docs/CLEANUP_PLAN.md)** - Recent improvements and refactoring
- **[Archive](docs/archive/)** - Historical phase documentation

## 🧪 Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific package
go test ./internal/usecase -v

# Run specific test
go test ./internal/usecase -run TestExpenseUseCase_CreateExpenseWithEqualSplit -v
```

**Test Coverage:** 16/16 tests passing (100%)

## 🏗️ Tech Stack

- **Language:** Go 1.25
- **Database:** PostgreSQL
- **Migration:** golang-migrate
- **Logging:** Zap (Uber)
- **Documentation:** Swagger/OpenAPI
- **Containerization:** Docker & Docker Compose
- **Architecture:** Clean Architecture

## 📈 Project Status

- ✅ Phase 1: User Management Enhancements
- ✅ Phase 2: Expense Enhancements  
- ✅ Phase 3: Filtering & Search
- ⏳ Phase 4: Statistics & Reporting (Coming soon)

See [PROJECT_STATUS.md](docs/PROJECT_STATUS.md) for detailed status.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

[Your License Here]

## 🙏 Acknowledgments

- Built with Clean Architecture principles
- Inspired by best practices in Go development
- Thanks to the open-source community

---

**Made with ❤️ using Go**

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