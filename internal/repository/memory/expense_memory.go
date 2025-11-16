package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/pavanrkadave/homies/internal/domain"
)

type ExpenseMemoryRepository struct {
	expenses map[string]*domain.Expense
	mu       sync.RWMutex
}

func NewExpenseMemoryRepository() *ExpenseMemoryRepository {
	return &ExpenseMemoryRepository{
		expenses: make(map[string]*domain.Expense),
	}
}

func (repo *ExpenseMemoryRepository) Create(ctx context.Context, expense *domain.Expense) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.expenses[expense.ID] = expense
	return nil
}

func (repo *ExpenseMemoryRepository) GetByID(ctx context.Context, id string) (*domain.Expense, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	expense, ok := repo.expenses[id]
	if !ok {
		return nil, fmt.Errorf(`expense "%s" not found`, id)
	}
	return expense, nil
}

func (repo *ExpenseMemoryRepository) GetAll(ctx context.Context) ([]*domain.Expense, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	expenses := make([]*domain.Expense, 0, len(repo.expenses))
	for _, expense := range repo.expenses {
		expenses = append(expenses, expense)
	}
	return expenses, nil
}

func (repo *ExpenseMemoryRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.Expense, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	expenses := make([]*domain.Expense, 0, len(repo.expenses))
	for _, expense := range repo.expenses {
		if expense.PaidBy == userID {
			expenses = append(expenses, expense)
			continue
		}

		for _, split := range expense.Splits {
			if split.UserID == userID {
				expenses = append(expenses, expense)
				break
			}
		}
	}
	return expenses, nil
}

func (repo *ExpenseMemoryRepository) Update(ctx context.Context, expense *domain.Expense) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, ok := repo.expenses[expense.ID]; !ok {
		return fmt.Errorf("expense not found")
	}

	repo.expenses[expense.ID] = expense
	return nil
}

func (repo *ExpenseMemoryRepository) GetByDateRange(ctx context.Context, startDate, endDate string) ([]*domain.Expense, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	expenses := make([]*domain.Expense, 0)
	for _, expense := range repo.expenses {
		dateStr := expense.Date.Format("2006-01-02")
		if dateStr >= startDate && dateStr <= endDate {
			expenses = append(expenses, expense)
		}
	}
	return expenses, nil
}

func (repo *ExpenseMemoryRepository) GetByCategory(ctx context.Context, category string) ([]*domain.Expense, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	expenses := make([]*domain.Expense, 0)
	for _, expense := range repo.expenses {
		if expense.Category == category {
			expenses = append(expenses, expense)
		}
	}
	return expenses, nil
}

func (repo *ExpenseMemoryRepository) GetByFilters(ctx context.Context, category, startDate, endDate string) ([]*domain.Expense, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	expenses := make([]*domain.Expense, 0)
	for _, expense := range repo.expenses {
		match := true

		if category != "" && expense.Category != category {
			match = false
		}

		if startDate != "" || endDate != "" {
			dateStr := expense.Date.Format("2006-01-02")
			if startDate != "" && dateStr < startDate {
				match = false
			}
			if endDate != "" && dateStr > endDate {
				match = false
			}
		}

		if match {
			expenses = append(expenses, expense)
		}
	}
	return expenses, nil
}

func (repo *ExpenseMemoryRepository) Delete(ctx context.Context, id string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	delete(repo.expenses, id)
	return nil
}
