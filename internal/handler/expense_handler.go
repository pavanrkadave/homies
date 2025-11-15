package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/pavanrkadave/homies/internal/domain"
	"github.com/pavanrkadave/homies/internal/usecase"
)

type ExpenseHandler struct {
	expenseUc usecase.ExpenseUseCase
}

func NewExpenseHandler(expenseUc usecase.ExpenseUseCase) *ExpenseHandler {
	return &ExpenseHandler{expenseUc: expenseUc}
}

type ExpenseRequest struct {
	Description string         `json:"description"`
	Amount      float64        `json:"amount"`
	Category    string         `json:"category"`
	PaidBy      string         `json:"paid_by"`
	Splits      []SplitRequest `json:"splits"`
}
type SplitRequest struct {
	UserId string  `json:"user_id"`
	Amount float64 `json:"amount"`
}

type ExpenseResponse struct {
	ID          string          `json:"id"`
	Description string          `json:"description"`
	Amount      float64         `json:"amount"`
	Category    string          `json:"category"`
	PaidBy      string          `json:"paid_by"`
	Date        time.Time       `json:"date"`
	CreatedAt   time.Time       `json:"created_at"`
	Splits      []SplitResponse `json:"splits"`
}
type SplitResponse struct {
	UserId string  `json:"user_id"`
	Amount float64 `json:"amount"`
}

func (h *ExpenseHandler) CreateExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req ExpenseRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	splits := make([]domain.Split, len(req.Splits))
	for i, split := range req.Splits {
		splits[i] = domain.Split{
			UserID: split.UserId,
			Amount: split.Amount,
		}
	}
	expense, err := h.expenseUc.CreateExpense(r.Context(), req.Description, req.Category, req.PaidBy, req.Amount, splits)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	responseSplits := make([]SplitResponse, len(expense.Splits))
	for i, split := range expense.Splits {
		responseSplits[i] = SplitResponse{
			UserId: split.UserID,
			Amount: split.Amount,
		}
	}

	expenseResponse := ExpenseResponse{
		ID:          expense.ID,
		Description: expense.Description,
		Amount:      expense.Amount,
		Category:    expense.Category,
		PaidBy:      expense.PaidBy,
		Date:        expense.Date,
		CreatedAt:   expense.CreatedAt,
		Splits:      responseSplits,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(expenseResponse)
}

func (h *ExpenseHandler) GetAllExpenses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	expenses, err := h.expenseUc.GetAllExpenses(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	expenseResponse := make([]ExpenseResponse, len(expenses))
	for i, expense := range expenses {

		splits := make([]SplitResponse, len(expense.Splits))
		for i, split := range expense.Splits {
			splits[i] = SplitResponse{
				UserId: split.UserID,
				Amount: split.Amount,
			}
		}

		expenseResponse[i] = ExpenseResponse{
			ID:          expense.ID,
			Description: expense.Description,
			Amount:      expense.Amount,
			Category:    expense.Category,
			PaidBy:      expense.PaidBy,
			Date:        expense.Date,
			CreatedAt:   expense.CreatedAt,
			Splits:      splits,
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(expenseResponse)
}

func (h *ExpenseHandler) GetBalances(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	balances, err := h.expenseUc.CalculateBalances(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(balances)
}
