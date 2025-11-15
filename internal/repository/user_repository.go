package repository

import (
	"context"

	"github.com/pavanrkadave/homies/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetAll(ctx context.Context) ([]*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}
