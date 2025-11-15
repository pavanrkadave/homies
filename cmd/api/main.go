package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/pavanrkadave/homies/config"
	"github.com/pavanrkadave/homies/internal/handler"
	"github.com/pavanrkadave/homies/internal/repository/postgres"
	"github.com/pavanrkadave/homies/internal/usecase"
	"github.com/pavanrkadave/homies/pkg/database"
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

	// API Routes
	http.HandleFunc("/users", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			userHandler.GetAllUsers(writer, request)
		case http.MethodPost:
			userHandler.CreateUser(writer, request)
		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/expenses", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			if request.URL.Query().Get("id") != "" {
				expenseHandler.GetExpenseByID(writer, request)
			} else {
				expenseHandler.GetAllExpenses(writer, request)
			}
		case http.MethodPost:
			expenseHandler.CreateExpense(writer, request)
		case http.MethodDelete:
			expenseHandler.DeleteExpense(writer, request)
		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/expenses/user", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodGet {
			expenseHandler.GetExpenseByUser(writer, request)
		} else {
			writer.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/balances", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			expenseHandler.GetBalances(writer, request)
		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	log.Printf("✓ Server starting on :%s", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, nil))
}
