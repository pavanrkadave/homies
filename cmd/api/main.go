package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/pavanrkadave/homies/config"
	"github.com/pavanrkadave/homies/internal/handler"
	"github.com/pavanrkadave/homies/internal/middleware"
	"github.com/pavanrkadave/homies/internal/repository/postgres"
	"github.com/pavanrkadave/homies/internal/usecase"
	"github.com/pavanrkadave/homies/pkg/database"
	"github.com/pavanrkadave/homies/pkg/response"
)

func main() {

	// Load Config
	cfg := config.Load()

	// Connect to the database
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal("failed to close database connection")
		}
	}(db)
	log.Println("✓ Connected to PostgreSQL successfully!")

	// Init Repositories
	userRepo := postgres.NewUserPostgresRepository(db)
	expenseRepo := postgres.NewExpensePostgresRepository(db)

	// Init UseCase
	userUC := usecase.NewUserUseCase(userRepo)
	expenseUC := usecase.NewExpenseUseCase(expenseRepo, userRepo)

	// Init Handlers
	userHandler := handler.NewUserHandler(userUC)
	expenseHandler := handler.NewExpenseHandler(expenseUC)
	healthHandler := handler.NewHealthHandler(db)

	mux := http.NewServeMux()

	// Healthcheck
	mux.HandleFunc("/health", healthHandler.Health)

	// API Routes
	mux.HandleFunc("/users", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			if request.URL.Query().Get("id") != "" {
				userHandler.GetUserByID(writer, request)
			} else {
				userHandler.GetAllUsers(writer, request)
			}
		case http.MethodPost:
			userHandler.CreateUser(writer, request)
		case http.MethodPut:
			userHandler.UpdateUser(writer, request)
		default:
			response.RespondWithError(writer, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	mux.HandleFunc("/expenses", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			if request.URL.Query().Get("id") != "" {
				expenseHandler.GetExpenseByID(writer, request)
			} else {
				expenseHandler.GetAllExpenses(writer, request)
			}
		case http.MethodPost:
			expenseHandler.CreateExpense(writer, request)
		case http.MethodPut:
			expenseHandler.UpdateExpense(writer, request)
		case http.MethodDelete:
			expenseHandler.DeleteExpense(writer, request)
		default:
			response.RespondWithError(writer, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	mux.HandleFunc("/expenses/equal-split", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodPost {
			expenseHandler.CreateExpenseWithEqualSplit(writer, request)
		} else {
			response.RespondWithError(writer, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	mux.HandleFunc("/expenses/user", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodGet {
			expenseHandler.GetExpenseByUser(writer, request)
		} else {
			response.RespondWithError(writer, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	mux.HandleFunc("/balances", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			expenseHandler.GetBalances(writer, request)
		default:
			response.RespondWithError(writer, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	middlewareHandler := middleware.Recovery(middleware.Logger(middleware.CORS(mux)))

	log.Printf("✓ Server starting on :%s with middleware enabled", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, middlewareHandler))
}
