package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/pavanrkadave/homies/internal/domain"
)

type mockExpenseRepository struct {
	expenses map[string]*domain.Expense
}

func (m *mockExpenseRepository) Create(ctx context.Context, expense *domain.Expense) error {
	m.expenses[expense.ID] = expense
	return nil
}

func (m *mockExpenseRepository) GetByID(ctx context.Context, id string) (*domain.Expense, error) {
	expense, ok := m.expenses[id]
	if !ok {
		return nil, errors.New("expense not found")
	}
	return expense, nil
}

func (m *mockExpenseRepository) GetAll(ctx context.Context) ([]*domain.Expense, error) {
	var expenses []*domain.Expense
	for _, expense := range m.expenses {
		expenses = append(expenses, expense)
	}
	return expenses, nil
}

func (m *mockExpenseRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.Expense, error) {
	var expenses []*domain.Expense
	for _, expense := range m.expenses {
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

func (m *mockExpenseRepository) Update(ctx context.Context, expense *domain.Expense) error {
	if _, ok := m.expenses[expense.ID]; !ok {
		return errors.New("expense not found")
	}
	m.expenses[expense.ID] = expense
	return nil
}

func (m *mockExpenseRepository) Delete(ctx context.Context, id string) error {
	delete(m.expenses, id)
	return nil
}

func newMockExpenseRepository() *mockExpenseRepository {
	return &mockExpenseRepository{
		expenses: make(map[string]*domain.Expense),
	}
}

func TestExpenseUseCase_UpdateExpense(t *testing.T) {
	expenseRepo := newMockExpenseRepository()
	userRepo := newMockUserRepository()
	expenseUC := NewExpenseUseCase(expenseRepo, userRepo)
	ctx := context.Background()

	// Create test users
	user1 := &domain.User{ID: "user1", Name: "User1", Email: "user1@test.com"}
	user2 := &domain.User{ID: "user2", Name: "User2", Email: "user2@test.com"}
	_ = userRepo.Create(ctx, user1)
	_ = userRepo.Create(ctx, user2)

	// Create an expense
	splits := []domain.Split{
		{UserID: user1.ID, Amount: 50.0},
		{UserID: user2.ID, Amount: 50.0},
	}
	expense, err := expenseUC.CreateExpense(ctx, "Dinner", "food", user1.ID, 100.0, splits)
	if err != nil {
		t.Fatalf("Failed to create expense: %v", err)
	}

	// Update the expense
	newSplits := []domain.Split{
		{UserID: user1.ID, Amount: 60.0},
		{UserID: user2.ID, Amount: 40.0},
	}
	updatedExpense, err := expenseUC.UpdateExpense(ctx, expense.ID, "Updated Dinner", "restaurant", 100.0, newSplits)
	if err != nil {
		t.Fatalf("Failed to update expense: %v", err)
	}

	// Verify updates
	if updatedExpense.Description != "Updated Dinner" {
		t.Errorf("Expected description 'Updated Dinner', got: %v", updatedExpense.Description)
	}
	if updatedExpense.Category != "restaurant" {
		t.Errorf("Expected category 'restaurant', got: %v", updatedExpense.Category)
	}
	if len(updatedExpense.Splits) != 2 {
		t.Errorf("Expected 2 splits, got: %v", len(updatedExpense.Splits))
	}
	if updatedExpense.Splits[0].Amount != 60.0 {
		t.Errorf("Expected split amount 60.0, got: %v", updatedExpense.Splits[0].Amount)
	}
}

func TestExpenseUseCase_UpdateExpense_NotFound(t *testing.T) {
	expenseRepo := newMockExpenseRepository()
	userRepo := newMockUserRepository()
	expenseUC := NewExpenseUseCase(expenseRepo, userRepo)
	ctx := context.Background()

	splits := []domain.Split{
		{UserID: "user1", Amount: 50.0},
	}
	_, err := expenseUC.UpdateExpense(ctx, "nonexistent", "Test", "test", 50.0, splits)
	if err == nil {
		t.Fatal("Expected error when updating non-existent expense, got nil")
	}
}

func TestExpenseUseCase_UpdateExpense_ValidationError(t *testing.T) {
	expenseRepo := newMockExpenseRepository()
	userRepo := newMockUserRepository()
	expenseUC := NewExpenseUseCase(expenseRepo, userRepo)
	ctx := context.Background()

	// Create test user
	user1 := &domain.User{ID: "user1", Name: "User1", Email: "user1@test.com"}
	_ = userRepo.Create(ctx, user1)

	// Create an expense
	splits := []domain.Split{
		{UserID: user1.ID, Amount: 100.0},
	}
	expense, err := expenseUC.CreateExpense(ctx, "Test", "test", user1.ID, 100.0, splits)
	if err != nil {
		t.Fatalf("Failed to create expense: %v", err)
	}

	// Try to update with invalid splits (sum doesn't match amount)
	invalidSplits := []domain.Split{
		{UserID: user1.ID, Amount: 50.0},
	}
	_, err = expenseUC.UpdateExpense(ctx, expense.ID, "Test", "test", 100.0, invalidSplits)
	if err == nil {
		t.Fatal("Expected validation error for invalid splits, got nil")
	}
}
