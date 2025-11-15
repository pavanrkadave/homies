package usecase

import (
	"context"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/pavanrkadave/homies/internal/domain"
	"github.com/pavanrkadave/homies/internal/repository"
)

type ExpenseUseCase interface {
	CreateExpense(ctx context.Context, description, category, paidBy string, amount float64, splits []domain.Split) (*domain.Expense, error)
	GetExpense(ctx context.Context, id string) (*domain.Expense, error)
	GetAllExpenses(ctx context.Context) ([]*domain.Expense, error)
	GetExpensesByUser(ctx context.Context, userID string) ([]*domain.Expense, error)
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
