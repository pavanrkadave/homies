package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/pavanrkadave/homies/internal/domain"
	"github.com/pavanrkadave/homies/internal/usecase"
	"github.com/pavanrkadave/homies/pkg/response"
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
		response.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req ExpenseRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
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
		response.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, ToExpenseResponse(expense))
}

func (h *ExpenseHandler) GetAllExpenses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	expenses, err := h.expenseUc.GetAllExpenses(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusOK, ToExpenseResponses(expenses))
}

func (h *ExpenseHandler) GetBalances(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	balances, err := h.expenseUc.CalculateBalances(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusOK, balances)
}

func (h *ExpenseHandler) GetExpenseByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		response.RespondWithError(w, http.StatusBadRequest, "id parameter is required")
		return
	}

	expense, err := h.expenseUc.GetExpense(r.Context(), id)
	if err != nil {
		response.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusOK, ToExpenseResponse(expense))
}

func (h *ExpenseHandler) GetExpenseByUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		response.RespondWithError(w, http.StatusBadRequest, "user_id parameter is required")
		return
	}

	expenses, err := h.expenseUc.GetExpensesByUser(r.Context(), userID)
	if err != nil {
		response.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusOK, ToExpenseResponses(expenses))
}

func (h *ExpenseHandler) UpdateExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		response.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		response.RespondWithError(w, http.StatusBadRequest, "id parameter is required")
		return
	}

	var req ExpenseRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	splits := make([]domain.Split, len(req.Splits))
	for i, split := range req.Splits {
		splits[i] = domain.Split{
			UserID: split.UserId,
			Amount: split.Amount,
		}
	}

	expense, err := h.expenseUc.UpdateExpense(r.Context(), id, req.Description, req.Category, req.Amount, splits)
	if err != nil {
		if err.Error() == "expense not found" {
			response.RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		response.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusOK, ToExpenseResponse(expense))
}

func (h *ExpenseHandler) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		response.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		response.RespondWithError(w, http.StatusBadRequest, "id parameter is required")
		return
	}

	err := h.expenseUc.DeleteExpense(r.Context(), id)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
