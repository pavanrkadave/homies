package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/pavanrkadave/homies/internal/domain"
)

type UserMemoryRepository struct {
	users map[string]*domain.User
	mu    sync.RWMutex
}

func NewUserMemoryRepository() *UserMemoryRepository {
	return &UserMemoryRepository{
		users: make(map[string]*domain.User),
	}
}

func (repo *UserMemoryRepository) Create(ctx context.Context, user *domain.User) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.users[user.ID] = user
	return nil
}

func (repo *UserMemoryRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	user, exists := repo.users[id]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (repo *UserMemoryRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	for _, user := range repo.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (repo *UserMemoryRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	var users []*domain.User
	for _, user := range repo.users {
		users = append(users, user)
	}
	return users, nil
}

func (repo *UserMemoryRepository) Update(ctx context.Context, user *domain.User) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.users[user.ID]; !exists {
		return fmt.Errorf("user not found")
	}

	repo.users[user.ID] = user
	return nil
}
