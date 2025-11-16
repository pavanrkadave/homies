package handler

import (
	"encoding/json"
	"fmt"
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

// CreateExpense godoc
// @Summary      Create a new expense
// @Description  Create a new expense with custom splits
// @Tags         expenses
// @Accept       json
// @Produce      json
// @Param        expense  body      ExpenseRequest  true  "Expense data"
// @Success      201      {object}  ExpenseResponse
// @Failure      400      {object}  map[string]string
// @Router       /expenses [post]
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

type EqualSplitRequest struct {
	Description string   `json:"description"`
	Amount      float64  `json:"amount"`
	Category    string   `json:"category"`
	PaidBy      string   `json:"paid_by"`
	UserIDs     []string `json:"user_ids"`
}

// CreateExpenseWithEqualSplit godoc
// @Summary      Create expense with equal split
// @Description  Create a new expense with equal splits among specified users
// @Tags         expenses
// @Accept       json
// @Produce      json
// @Param        expense  body      EqualSplitRequest  true  "Equal split expense data"
// @Success      201      {object}  ExpenseResponse
// @Failure      400      {object}  map[string]string
// @Router       /expenses/equal-split [post]
func (h *ExpenseHandler) CreateExpenseWithEqualSplit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req EqualSplitRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	expense, err := h.expenseUc.CreateExpenseWithEqualSplit(r.Context(), req.Description, req.Category, req.PaidBy, req.Amount, req.UserIDs)
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, ToExpenseResponse(expense))
}

// GetAllExpenses godoc
// @Summary      Get all expenses
// @Description  Retrieve all expenses with optional filters (category, date range)
// @Tags         expenses
// @Produce      json
// @Param        category    query     string  false  "Filter by category"
// @Param        start_date  query     string  false  "Start date (YYYY-MM-DD)"
// @Param        end_date    query     string  false  "End date (YYYY-MM-DD)"
// @Success      200         {array}   ExpenseResponse
// @Failure      400         {object}  map[string]string
// @Router       /expenses [get]
func (h *ExpenseHandler) GetAllExpenses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Check for filter query parameters
	category := r.URL.Query().Get("category")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	var expenses []*domain.Expense
	var err error

	// Use filters if any are provided
	if category != "" || startDate != "" || endDate != "" {
		expenses, err = h.expenseUc.GetExpensesByFilters(r.Context(), category, startDate, endDate)
	} else {
		expenses, err = h.expenseUc.GetAllExpenses(r.Context())
	}

	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusOK, ToExpenseResponses(expenses))
}

// GetBalances godoc
// @Summary      Get all balances
// @Description  Calculate and retrieve balances between all users
// @Tags         balances
// @Produce      json
// @Success      200  {array}   domain.Balance
// @Failure      500  {object}  map[string]string
// @Router       /balances [get]
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

// GetExpenseByUser godoc
// @Summary      Get expenses by user
// @Description  Retrieve all expenses for a specific user
// @Tags         expenses
// @Produce      json
// @Param        user_id  query     string  true  "User ID"
// @Success      200      {array}   ExpenseResponse
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Router       /expenses/user [get]
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

// UpdateExpense godoc
// @Summary      Update an expense
// @Description  Update an existing expense by ID
// @Tags         expenses
// @Accept       json
// @Produce      json
// @Param        id       query     string          true  "Expense ID"
// @Param        expense  body      ExpenseRequest  true  "Updated expense data"
// @Success      200      {object}  ExpenseResponse
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Router       /expenses [put]
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

// DeleteExpense godoc
// @Summary      Delete an expense
// @Description  Delete an expense by ID
// @Tags         expenses
// @Param        id   query     string  true  "Expense ID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /expenses [delete]
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

// GetUserStats godoc
// @Summary      Get user statistics
// @Description  Get spending statistics for a specific user
// @Tags         statistics
// @Produce      json
// @Param        user_id  query     string  true  "User ID"
// @Success      200      {object}  domain.UserStats
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Router       /users/stats [get]
func (h *ExpenseHandler) GetUserStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		response.RespondWithError(w, http.StatusBadRequest, "user_id parameter is required")
		return
	}

	stats, err := h.expenseUc.GetUserStats(r.Context(), userID)
	if err != nil {
		response.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusOK, stats)
}

// GetMonthlySummary godoc
// @Summary      Get monthly summary
// @Description  Get expense summary for a specific month
// @Tags         statistics
// @Produce      json
// @Param        year   query     int  true  "Year"
// @Param        month  query     int  true  "Month (1-12)"
// @Success      200    {object}  domain.MonthlySummary
// @Failure      400    {object}  map[string]string
// @Router       /expenses/monthly [get]
func (h *ExpenseHandler) GetMonthlySummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	yearStr := r.URL.Query().Get("year")
	monthStr := r.URL.Query().Get("month")

	if yearStr == "" || monthStr == "" {
		response.RespondWithError(w, http.StatusBadRequest, "year and month parameters are required")
		return
	}

	year := 0
	month := 0
	if _, err := fmt.Sscanf(yearStr, "%d", &year); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "invalid year format")
		return
	}

	if _, err := fmt.Sscanf(monthStr, "%d", &month); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "invalid month format")
		return
	}

	summary, err := h.expenseUc.GetMonthlySummary(r.Context(), year, month)
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusOK, summary)
}
