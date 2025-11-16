.PHONY: help build run test clean docker-up docker-down docker-rebuild migrate-up migrate-down swagger logs

help:
	@echo "Homies Expense Tracker - Available Commands:"
	@echo ""
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application locally"
	@echo "  make test           - Run all tests"
	@echo "  make test-verbose   - Run tests with verbose output"
	@echo "  make clean          - Clean build artifacts"
	@echo ""
	@echo "  make docker-up      - Start Docker containers"
	@echo "  make docker-down    - Stop Docker containers"
	@echo "  make docker-rebuild - Rebuild and restart containers"
	@echo "  make logs           - View application logs"
	@echo ""
	@echo "  make migrate-up     - Run database migrations"
	@echo "  make migrate-down   - Rollback last migration"
	@echo ""
	@echo "  make swagger        - Generate Swagger documentation"
	@echo "  make swagger-serve  - Serve Swagger UI locally"
	@echo ""
	@echo "  make lint           - Run linter"
	@echo "  make fmt            - Format code"

build:
	@echo "Building application..."
	@go build -o bin/homies cmd/api/main.go
	@echo "✓ Build complete: bin/homies"

run:
	@echo "Starting application..."
	@go run cmd/api/main.go

test:
	@echo "Running tests..."
	@go test ./...

test-verbose:
	@echo "Running tests (verbose)..."
	@go test ./... -v

test-coverage:
	@echo "Running tests with coverage..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated: coverage.html"

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "✓ Clean complete"

docker-up:
	@echo "Starting Docker containers..."
	@docker-compose up -d
	@echo "✓ Containers started"
	@docker-compose ps

docker-down:
	@echo "Stopping Docker containers..."
	@docker-compose down
	@echo "✓ Containers stopped"

docker-rebuild:
	@echo "Rebuilding and restarting containers..."
	@docker-compose down
	@docker-compose up -d --build
	@echo "✓ Containers rebuilt and started"
	@docker-compose ps

logs:
	@docker-compose logs -f app

logs-db:
	@docker-compose logs -f postgres

migrate-up:
	@echo "Running migrations..."
	@go run cmd/migrate/main.go up

migrate-down:
	@echo "Rolling back last migration..."
	@go run cmd/migrate/main.go down

migrate-create:
	@read -p "Enter migration name: " name; \
	timestamp=$$(date +%s); \
	touch migrations/$${timestamp}_$${name}.up.sql; \
	touch migrations/$${timestamp}_$${name}.down.sql; \
	echo "✓ Created migration files:"; \
	echo "  migrations/$${timestamp}_$${name}.up.sql"; \
	echo "  migrations/$${timestamp}_$${name}.down.sql"

swagger:
	@echo "Generating Swagger documentation..."
	@which swag > /dev/null || (echo "Installing swag..." && go install github.com/swaggo/swag/cmd/swag@latest)
	@swag init -g cmd/api/main.go -o docs/swagger
	@echo "✓ Swagger docs generated: docs/swagger/"

swagger-serve:
	@echo "Swagger UI available at: http://localhost:3000/swagger/index.html"
	@echo "Make sure the application is running (make run or make docker-up)"

lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	@golangci-lint run
	@echo "✓ Lint complete"

fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✓ Format complete"

mod-tidy:
	@echo "Tidying go.mod..."
	@go mod tidy
	@echo "✓ go.mod tidied"

install-tools:
	@echo "Installing development tools..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✓ Tools installed"

db-shell:
	@echo "Connecting to PostgreSQL..."
	@docker-compose exec postgres psql -U postgres -d homies

# Development workflow
dev: docker-up migrate-up swagger
	@echo "✓ Development environment ready!"
	@echo "  - Database: running"
	@echo "  - Migrations: applied"
	@echo "  - Swagger: generated"
	@echo ""
	@echo "Run 'make run' to start the application"

# Production build
prod-build:
	@echo "Building for production..."
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/homies cmd/api/main.go
	@echo "✓ Production build complete: bin/homies"

