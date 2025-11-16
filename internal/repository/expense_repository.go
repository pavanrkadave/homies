package repository

import (
	"context"

	"github.com/pavanrkadave/homies/internal/domain"
)

type ExpenseRepository interface {
	Create(ctx context.Context, expense *domain.Expense) error
	GetByID(ctx context.Context, id string) (*domain.Expense, error)
	GetAll(ctx context.Context) ([]*domain.Expense, error)
	GetByUserID(ctx context.Context, userID string) ([]*domain.Expense, error)
	Update(ctx context.Context, expense *domain.Expense) error
	Delete(ctx context.Context, id string) error
}
