package usecase

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/pavanrkadave/homies/internal/domain"
	"github.com/pavanrkadave/homies/internal/repository"
)

type ExpenseUseCase interface {
	CreateExpense(ctx context.Context, description, category, paidBy string, amount float64, splits []domain.Split) (*domain.Expense, error)
	CreateExpenseWithEqualSplit(ctx context.Context, description, category, paidBy string, amount float64, userIDs []string) (*domain.Expense, error)
	GetExpense(ctx context.Context, id string) (*domain.Expense, error)
	GetAllExpenses(ctx context.Context) ([]*domain.Expense, error)
	GetExpensesByUser(ctx context.Context, userID string) ([]*domain.Expense, error)
	GetExpensesByDateRange(ctx context.Context, startDate, endDate string) ([]*domain.Expense, error)
	GetExpensesByCategory(ctx context.Context, category string) ([]*domain.Expense, error)
	GetExpensesByFilters(ctx context.Context, category, startDate, endDate string) ([]*domain.Expense, error)
	UpdateExpense(ctx context.Context, id, description, category string, amount float64, splits []domain.Split) (*domain.Expense, error)
	DeleteExpense(ctx context.Context, id string) error
	CalculateBalances(ctx context.Context) (*domain.BalanceSummary, error)
}

type expenseUseCase struct {
	expenseRepo repository.ExpenseRepository
	userRepo    repository.UserRepository
}

func NewExpenseUseCase(expenseRepo repository.ExpenseRepository, userRepo repository.UserRepository) ExpenseUseCase {
	return &expenseUseCase{
		expenseRepo: expenseRepo,
		userRepo:    userRepo,
	}
}

func (e *expenseUseCase) CreateExpense(ctx context.Context, description, category, paidBy string, amount float64, splits []domain.Split) (*domain.Expense, error) {
	expenseId := uuid.New().String()

	user, err := e.userRepo.GetByID(ctx, paidBy)
	if err != nil {
		return nil, err
	}

	for _, split := range splits {
		_, err := e.userRepo.GetByID(ctx, split.UserID)
		if err != nil {
			return nil, err
		}
	}

	expense := &domain.Expense{
		ID:          expenseId,
		Description: description,
		Category:    category,
		PaidBy:      user.ID,
		Amount:      amount,
		Splits:      splits,
		Date:        time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = expense.Validate()
	if err != nil {
		return nil, err
	}

	err = e.expenseRepo.Create(ctx, expense)
	if err != nil {
		return nil, err
	}

	return expense, nil
}

func (e *expenseUseCase) CreateExpenseWithEqualSplit(ctx context.Context, description, category, paidBy string, amount float64, userIDs []string) (*domain.Expense, error) {
	if len(userIDs) == 0 {
		return nil, errors.New("at least one user is required for equal split")
	}

	// Validate all users exist
	for _, userID := range userIDs {
		_, err := e.userRepo.GetByID(ctx, userID)
		if err != nil {
			return nil, err
		}
	}

	// Calculate equal split amount
	splitAmount := amount / float64(len(userIDs))

	// Handle rounding: calculate remainder and add to last user
	var splits []domain.Split
	var totalAllocated float64

	for i, userID := range userIDs {
		if i == len(userIDs)-1 {
			// Last user gets the remainder to ensure exact total
			splits = append(splits, domain.Split{
				UserID: userID,
				Amount: amount - totalAllocated,
			})
		} else {
			// Round to 2 decimal places
			roundedAmount := math.Round(splitAmount*100) / 100
			splits = append(splits, domain.Split{
				UserID: userID,
				Amount: roundedAmount,
			})
			totalAllocated += roundedAmount
		}
	}

	return e.CreateExpense(ctx, description, category, paidBy, amount, splits)
}

func (e *expenseUseCase) GetExpense(ctx context.Context, id string) (*domain.Expense, error) {
	return e.expenseRepo.GetByID(ctx, id)
}

func (e *expenseUseCase) GetAllExpenses(ctx context.Context) ([]*domain.Expense, error) {
	return e.expenseRepo.GetAll(ctx)
}

func (e *expenseUseCase) GetExpensesByUser(ctx context.Context, userID string) ([]*domain.Expense, error) {
	_, err := e.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return e.expenseRepo.GetByUserID(ctx, userID)
}

func (e *expenseUseCase) GetExpensesByDateRange(ctx context.Context, startDate, endDate string) ([]*domain.Expense, error) {
	if startDate == "" || endDate == "" {
		return nil, errors.New("both start_date and end_date are required")
	}
	return e.expenseRepo.GetByDateRange(ctx, startDate, endDate)
}

func (e *expenseUseCase) GetExpensesByCategory(ctx context.Context, category string) ([]*domain.Expense, error) {
	if category == "" {
		return nil, errors.New("category is required")
	}
	return e.expenseRepo.GetByCategory(ctx, category)
}

func (e *expenseUseCase) GetExpensesByFilters(ctx context.Context, category, startDate, endDate string) ([]*domain.Expense, error) {
	// If no filters provided, return all expenses
	if category == "" && startDate == "" && endDate == "" {
		return e.expenseRepo.GetAll(ctx)
	}

	// Validate date range if provided
	if (startDate != "" && endDate == "") || (startDate == "" && endDate != "") {
		return nil, errors.New("both start_date and end_date must be provided together")
	}

	return e.expenseRepo.GetByFilters(ctx, category, startDate, endDate)
}

func (e *expenseUseCase) UpdateExpense(ctx context.Context, id, description, category string, amount float64, splits []domain.Split) (*domain.Expense, error) {
	// Get existing expense
	expense, err := e.expenseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate users in splits exist
	for _, split := range splits {
		_, err := e.userRepo.GetByID(ctx, split.UserID)
		if err != nil {
			return nil, err
		}
	}

	// Update expense fields
	err = expense.Update(description, category, amount, splits)
	if err != nil {
		return nil, err
	}

	// Save to repository
	err = e.expenseRepo.Update(ctx, expense)
	if err != nil {
		return nil, err
	}

	return expense, nil
}

func (e *expenseUseCase) DeleteExpense(ctx context.Context, id string) error {
	return e.expenseRepo.Delete(ctx, id)
}

func (e *expenseUseCase) CalculateBalances(ctx context.Context) (*domain.BalanceSummary, error) {

	expenses, err := e.expenseRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	balanceMap := make(map[string]float64)

	for _, expense := range expenses {
		balanceMap[expense.PaidBy] += expense.Amount
		for _, split := range expense.Splits {
			balanceMap[split.UserID] -= split.Amount
		}
	}

	var balances []domain.Balance
	for userID, amount := range balanceMap {
		balances = append(balances, domain.Balance{
			UserID: userID,
			Amount: amount,
		})
	}

	var settlements []domain.Settlement
	settlements = calculateSettlements(balances)

	return &domain.BalanceSummary{
		Balances:    balances,
		Settlements: settlements,
	}, nil

}

func calculateSettlements(balances []domain.Balance) []domain.Settlement {
	var settlements []domain.Settlement

	var debtors []domain.Balance
	var creditors []domain.Balance

	for _, balance := range balances {
		if balance.Amount < 0 {
			debtors = append(debtors, balance)
		}
		if balance.Amount > 0 {
			creditors = append(creditors, balance)
		}
	}

	for len(debtors) > 0 && len(creditors) > 0 {
		debtor := debtors[0]
		creditor := creditors[0]

		settlementAmount := math.Min(-debtor.Amount, creditor.Amount)

		settlements = append(settlements, domain.Settlement{
			From:   debtor.UserID,
			To:     creditor.UserID,
			Amount: settlementAmount,
		})

		debtors[0].Amount += settlementAmount
		creditors[0].Amount -= settlementAmount

		if math.Abs(debtors[0].Amount) < 0.01 {
			debtors = debtors[1:]
		}

		if math.Abs(creditors[0].Amount) < 0.01 {
			creditors = creditors[1:]
		}
	}
	return settlements
}
