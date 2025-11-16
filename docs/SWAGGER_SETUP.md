# Swagger/OpenAPI Documentation Setup

## Overview

The Homies Expense Tracker API now includes complete Swagger/OpenAPI documentation with an interactive UI for testing and exploring all endpoints.

## Access Swagger UI

Once the application is running, access the Swagger UI at:
```
http://localhost:3000/swagger/index.html
```

## Available Endpoints

### Total: 8 Paths with Multiple HTTP Methods

#### 1. **Health Check**
- `GET /health` - Check API and database health status

#### 2. **Users**
- `POST /users` - Create a new user
- `GET /users` - Get all users
- `GET /users?id={id}` - Get user by ID
- `PUT /users?id={id}` - Update user

#### 3. **Expenses**
- `POST /expenses` - Create expense with custom splits
- `GET /expenses` - Get all expenses (with optional filters)
  - Query params: `category`, `start_date`, `end_date`
- `GET /expenses?id={id}` - Get expense by ID
- `PUT /expenses?id={id}` - Update expense
- `DELETE /expenses?id={id}` - Delete expense
- `POST /expenses/equal-split` - Create expense with equal splits
- `GET /expenses/user?user_id={id}` - Get expenses by user
- `GET /expenses/monthly?year={year}&month={month}` - Get monthly summary

#### 4. **Balances**
- `GET /balances` - Calculate and retrieve all balances

#### 5. **Statistics**
- `GET /users/stats?user_id={id}` - Get user spending statistics

## Generating Swagger Documentation

### Manual Generation
```bash
make swagger
```

This command:
1. Checks if `swag` is installed (installs if missing)
2. Scans the codebase for swagger annotations
3. Generates swagger files in `docs/swagger/`
4. Outputs: `docs.go`, `swagger.json`, `swagger.yaml`

### Development with Live Swagger UI
```bash
make swagger-serve
```

This starts the application with Swagger UI available immediately.

## Swagger Annotations

All handlers include complete swagger annotations:

### Example Annotation Structure
```go
// CreateUser godoc
// @Summary      Create a new user
// @Description  Create a new user with name and email
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      CreateUserRequest  true  "User data"
// @Success      201   {object}  UserResponse
// @Failure      400   {object}  map[string]string
// @Router       /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    // Implementation...
}
```

## API Information

- **Title**: Homies Expense Tracker API
- **Version**: 1.0
- **Description**: A production-ready expense tracker REST API for roommates to track shared expenses, split costs, and calculate settlements. Built with Go following Clean Architecture principles.
- **Host**: localhost:3000
- **Base Path**: /
- **Schemes**: http, https

## Tags Organization

Endpoints are organized into logical groups:

1. **health** - Health check endpoints
2. **users** - User management endpoints
3. **expenses** - Expense CRUD and operations
4. **balances** - Balance calculation endpoints
5. **statistics** - Analytics and reporting endpoints

## Dependencies

The following packages were added for Swagger support:

```go
import (
    _ "github.com/pavanrkadave/homies/docs/swagger"
    httpSwagger "github.com/swaggo/http-swagger"
)
```

**Installed packages:**
- `github.com/swaggo/swag/cmd/swag` - CLI tool for generating docs
- `github.com/swaggo/http-swagger` - HTTP handler for Swagger UI
- `github.com/swaggo/files` - Static file serving
- `github.com/go-openapi/*` - OpenAPI specification libraries

## Makefile Commands

```makefile
# Generate Swagger documentation
make swagger

# Start app with Swagger UI
make swagger-serve

# View help for all commands
make help
```

## Testing with Swagger UI

1. Start the application:
   ```bash
   make run
   # or
   make docker-up
   ```

2. Open browser to http://localhost:3000/swagger/index.html

3. Explore endpoints by:
   - Clicking "Try it out" on any endpoint
   - Filling in required parameters
   - Clicking "Execute"
   - Viewing the response

## Regenerating Documentation

If you add new endpoints or modify existing ones:

1. Add swagger annotations to the handler function
2. Run `make swagger` to regenerate documentation
3. Restart the application
4. Refresh the Swagger UI page

## Common Swagger Annotations

### Request Parameters
```go
// @Param id query string true "User ID"           // Query parameter
// @Param user body CreateUserRequest true "Data"  // Request body
```

### Response Types
```go
// @Success 200 {object} UserResponse              // Single object
// @Success 200 {array} UserResponse               // Array
// @Failure 400 {object} map[string]string         // Error response
```

### Tags and Metadata
```go
// @Tags users                    // Group in Swagger UI
// @Summary Short description     // Brief summary
// @Description Longer description // Detailed info
// @Accept json                   // Content-Type
// @Produce json                  // Response type
```

## Known Issues

### Warning: Duplicate Routes
You may see: `warning: route GET /users is declared multiple times`

This is expected because:
- `GET /users` returns all users
- `GET /users?id={id}` returns a specific user

Both use the same HTTP method and path, differentiated only by query parameters. Swagger picks up both but issues a warning. This does not affect functionality.

## Production Deployment

For production, you may want to:

1. **Disable Swagger UI** (security):
   ```go
   if cfg.Environment != "production" {
       mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)
   }
   ```

2. **Serve static swagger.json** for API consumers:
   ```
   https://api.example.com/swagger.json
   ```

3. **Use external Swagger UI** pointing to your API's swagger.json

## Additional Resources

- [Swaggo Documentation](https://github.com/swaggo/swag)
- [Swagger/OpenAPI Specification](https://swagger.io/specification/)
- [HTTP Swagger Handler](https://github.com/swaggo/http-swagger)

## Troubleshooting

### Swagger UI not loading
- Ensure application is running
- Check that `/swagger/` route is registered
- Verify docs package is imported: `_ "github.com/pavanrkadave/homies/docs/swagger"`

### Documentation not updating
- Run `make swagger` to regenerate
- Restart the application
- Clear browser cache

### Build errors after adding swagger
- Run `go mod tidy`
- Ensure all swagger annotations are syntactically correct
- Check that request/response types exist

## Summary

✅ Complete API documentation with 8 endpoint paths
✅ Interactive Swagger UI for testing
✅ Automatic schema generation from Go structs
✅ Organized by tags (health, users, expenses, balances, statistics)
✅ Production-ready with detailed descriptions
✅ Easy to maintain - just add annotations to new endpoints

