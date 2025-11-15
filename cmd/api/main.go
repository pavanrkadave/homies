package main

import (
	"log"
	"net/http"

	"github.com/pavanrkadave/homies/internal/handler"
	"github.com/pavanrkadave/homies/internal/repository/memory"
	"github.com/pavanrkadave/homies/internal/usecase"
)

func main() {
	// Init all Dependencies
	userRepo := memory.NewUserMemoryRepository()
	userUseCase := usecase.NewUserUseCase(userRepo)
	userHandler := handler.NewUserHandler(userUseCase)

	expenseRepo := memory.NewExpenseMemoryRepository()
	expenseUseCase := usecase.NewExpenseUseCase(expenseRepo, userRepo)
	expenseHandler := handler.NewExpenseHandler(expenseUseCase)

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
			expenseHandler.GetAllExpenses(writer, request)
		case http.MethodPost:
			expenseHandler.CreateExpense(writer, request)
		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	log.Println("Server starting on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
