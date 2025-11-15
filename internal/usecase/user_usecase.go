package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pavanrkadave/homies/internal/domain"
	"github.com/pavanrkadave/homies/internal/repository"
)

type UserUseCase interface {
	CreateUser(ctx context.Context, name, email string) (*domain.User, error)
	GetUser(ctx context.Context, id string) (*domain.User, error)
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
	UpdateUser(ctx context.Context, id, name, email string) (*domain.User, error)
}

type userUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

func (u *userUseCase) CreateUser(ctx context.Context, name, email string) (*domain.User, error) {
	userId := uuid.New().String()
	user := &domain.User{
		ID:        userId,
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := user.Validate()
	if err != nil {
		return nil, err
	}

	err = u.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUseCase) GetUser(ctx context.Context, id string) (*domain.User, error) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUseCase) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	users, err := u.userRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userUseCase) UpdateUser(ctx context.Context, id, name, email string) (*domain.User, error) {
	// Get existing user
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if email is being changed and if new email already exists
	if email != user.Email {
		existingUser, err := u.userRepo.GetByEmail(ctx, email)
		if err == nil && existingUser != nil {
			return nil, domain.ErrEmailAlreadyExists
		}
	}

	// Update user fields
	user.Name = name
	user.Email = email
	user.UpdatedAt = time.Now()

	// Validate updated user
	if err := user.Validate(); err != nil {
		return nil, err
	}

	// Save to repository
	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
