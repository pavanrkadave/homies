package usecase

import (
	"context"
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
