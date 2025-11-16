package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

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

func (m *mockExpenseRepository) GetByDateRange(ctx context.Context, startDate, endDate string) ([]*domain.Expense, error) {
	var expenses []*domain.Expense
	for _, expense := range m.expenses {
		dateStr := expense.Date.Format("2006-01-02")
		if dateStr >= startDate && dateStr <= endDate {
			expenses = append(expenses, expense)
		}
	}
	return expenses, nil
}

func (m *mockExpenseRepository) GetByCategory(ctx context.Context, category string) ([]*domain.Expense, error) {
	var expenses []*domain.Expense
	for _, expense := range m.expenses {
		if expense.Category == category {
			expenses = append(expenses, expense)
		}
	}
	return expenses, nil
}

func (m *mockExpenseRepository) GetByFilters(ctx context.Context, category, startDate, endDate string) ([]*domain.Expense, error) {
	var expenses []*domain.Expense
	for _, expense := range m.expenses {
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

func TestExpenseUseCase_CreateExpenseWithEqualSplit(t *testing.T) {
	expenseRepo := newMockExpenseRepository()
	userRepo := newMockUserRepository()
	expenseUC := NewExpenseUseCase(expenseRepo, userRepo)
	ctx := context.Background()

	// Create test users
	user1 := &domain.User{ID: "user1", Name: "User1", Email: "user1@test.com"}
	user2 := &domain.User{ID: "user2", Name: "User2", Email: "user2@test.com"}
	user3 := &domain.User{ID: "user3", Name: "User3", Email: "user3@test.com"}
	_ = userRepo.Create(ctx, user1)
	_ = userRepo.Create(ctx, user2)
	_ = userRepo.Create(ctx, user3)

	// Create expense with equal split
	userIDs := []string{user1.ID, user2.ID, user3.ID}
	expense, err := expenseUC.CreateExpenseWithEqualSplit(ctx, "Team Dinner", "food", user1.ID, 100.0, userIDs)
	if err != nil {
		t.Fatalf("Failed to create expense with equal split: %v", err)
	}

	// Verify expense was created
	if expense.Description != "Team Dinner" {
		t.Errorf("Expected description 'Team Dinner', got: %v", expense.Description)
	}
	if expense.Amount != 100.0 {
		t.Errorf("Expected amount 100.0, got: %v", expense.Amount)
	}

	// Verify splits
	if len(expense.Splits) != 3 {
		t.Fatalf("Expected 3 splits, got: %v", len(expense.Splits))
	}

	// Calculate total of splits
	var total float64
	for _, split := range expense.Splits {
		total += split.Amount
	}

	// Verify total matches expense amount
	if total != 100.0 {
		t.Errorf("Expected total splits to be 100.0, got: %v", total)
	}
}

func TestExpenseUseCase_CreateExpenseWithEqualSplit_NoUsers(t *testing.T) {
	expenseRepo := newMockExpenseRepository()
	userRepo := newMockUserRepository()
	expenseUC := NewExpenseUseCase(expenseRepo, userRepo)
	ctx := context.Background()

	// Try to create expense with no users
	_, err := expenseUC.CreateExpenseWithEqualSplit(ctx, "Test", "test", "user1", 100.0, []string{})
	if err == nil {
		t.Fatal("Expected error when creating equal split with no users, got nil")
	}
}

func TestExpenseUseCase_CreateExpenseWithEqualSplit_UnevenAmount(t *testing.T) {
	expenseRepo := newMockExpenseRepository()
	userRepo := newMockUserRepository()
	expenseUC := NewExpenseUseCase(expenseRepo, userRepo)
	ctx := context.Background()

	// Create test users
	user1 := &domain.User{ID: "user1", Name: "User1", Email: "user1@test.com"}
	user2 := &domain.User{ID: "user2", Name: "User2", Email: "user2@test.com"}
	user3 := &domain.User{ID: "user3", Name: "User3", Email: "user3@test.com"}
	_ = userRepo.Create(ctx, user1)
	_ = userRepo.Create(ctx, user2)
	_ = userRepo.Create(ctx, user3)

	// Create expense with amount that doesn't divide evenly (100 / 3 = 33.33...)
	userIDs := []string{user1.ID, user2.ID, user3.ID}
	expense, err := expenseUC.CreateExpenseWithEqualSplit(ctx, "Uneven Split", "test", user1.ID, 100.0, userIDs)
	if err != nil {
		t.Fatalf("Failed to create expense with uneven split: %v", err)
	}

	// Verify total still equals expense amount (rounding handled correctly)
	var total float64
	for _, split := range expense.Splits {
		total += split.Amount
	}

	if total != 100.0 {
		t.Errorf("Expected total splits to be 100.0 even with rounding, got: %v", total)
	}
}

func TestExpenseUseCase_GetExpensesByCategory(t *testing.T) {
	expenseRepo := newMockExpenseRepository()
	userRepo := newMockUserRepository()
	expenseUC := NewExpenseUseCase(expenseRepo, userRepo)
	ctx := context.Background()

	// Create test user
	user1 := &domain.User{ID: "user1", Name: "User1", Email: "user1@test.com"}
	_ = userRepo.Create(ctx, user1)

	// Create expenses with different categories
	_, _ = expenseUC.CreateExpense(ctx, "Groceries", "food", user1.ID, 100.0, []domain.Split{{UserID: user1.ID, Amount: 100.0}})
	_, _ = expenseUC.CreateExpense(ctx, "Restaurant", "food", user1.ID, 50.0, []domain.Split{{UserID: user1.ID, Amount: 50.0}})
	_, _ = expenseUC.CreateExpense(ctx, "Movie", "entertainment", user1.ID, 20.0, []domain.Split{{UserID: user1.ID, Amount: 20.0}})

	// Get expenses by category
	expenses, err := expenseUC.GetExpensesByCategory(ctx, "food")
	if err != nil {
		t.Fatalf("Failed to get expenses by category: %v", err)
	}

	// Verify we got 2 food expenses
	if len(expenses) != 2 {
		t.Errorf("Expected 2 food expenses, got: %v", len(expenses))
	}
}

func TestExpenseUseCase_GetExpensesByCategory_EmptyCategory(t *testing.T) {
	expenseRepo := newMockExpenseRepository()
	userRepo := newMockUserRepository()
	expenseUC := NewExpenseUseCase(expenseRepo, userRepo)
	ctx := context.Background()

	// Try to get expenses with empty category
	_, err := expenseUC.GetExpensesByCategory(ctx, "")
	if err == nil {
		t.Fatal("Expected error for empty category, got nil")
	}
}

func TestExpenseUseCase_GetExpensesByFilters(t *testing.T) {
	expenseRepo := newMockExpenseRepository()
	userRepo := newMockUserRepository()
	expenseUC := NewExpenseUseCase(expenseRepo, userRepo)
	ctx := context.Background()

	// Create test user
	user1 := &domain.User{ID: "user1", Name: "User1", Email: "user1@test.com"}
	_ = userRepo.Create(ctx, user1)

	// Create test expenses
	_, _ = expenseUC.CreateExpense(ctx, "Groceries", "food", user1.ID, 100.0, []domain.Split{{UserID: user1.ID, Amount: 100.0}})
	_, _ = expenseUC.CreateExpense(ctx, "Movie", "entertainment", user1.ID, 20.0, []domain.Split{{UserID: user1.ID, Amount: 20.0}})

	// Test with category filter only
	expenses, err := expenseUC.GetExpensesByFilters(ctx, "food", "", "")
	if err != nil {
		t.Fatalf("Failed to get expenses by filters: %v", err)
	}

	if len(expenses) != 1 {
		t.Errorf("Expected 1 food expense, got: %v", len(expenses))
	}
}

func TestExpenseUseCase_GetExpensesByFilters_NoFilters(t *testing.T) {
	expenseRepo := newMockExpenseRepository()
	userRepo := newMockUserRepository()
	expenseUC := NewExpenseUseCase(expenseRepo, userRepo)
	ctx := context.Background()

	// Create test user
	user1 := &domain.User{ID: "user1", Name: "User1", Email: "user1@test.com"}
	_ = userRepo.Create(ctx, user1)

	// Create test expenses
	_, _ = expenseUC.CreateExpense(ctx, "Expense 1", "food", user1.ID, 100.0, []domain.Split{{UserID: user1.ID, Amount: 100.0}})
	_, _ = expenseUC.CreateExpense(ctx, "Expense 2", "entertainment", user1.ID, 50.0, []domain.Split{{UserID: user1.ID, Amount: 50.0}})

	// Get all expenses (no filters)
	expenses, err := expenseUC.GetExpensesByFilters(ctx, "", "", "")
	if err != nil {
		t.Fatalf("Failed to get expenses: %v", err)
	}

	// Should return all expenses
	if len(expenses) != 2 {
		t.Errorf("Expected 2 expenses, got: %v", len(expenses))
	}
}

func TestExpenseUseCase_GetUserStats(t *testing.T) {
	expenseRepo := newMockExpenseRepository()
	userRepo := newMockUserRepository()
	expenseUC := NewExpenseUseCase(expenseRepo, userRepo)
	ctx := context.Background()

	// Create test users
	user1 := &domain.User{ID: "user1", Name: "User1", Email: "user1@test.com"}
	user2 := &domain.User{ID: "user2", Name: "User2", Email: "user2@test.com"}
	_ = userRepo.Create(ctx, user1)
	_ = userRepo.Create(ctx, user2)

	// Create expenses
	// User1 paid 150 total (100 food + 50 entertainment)
	_, _ = expenseUC.CreateExpense(ctx, "Groceries", "food", user1.ID, 100.0, []domain.Split{
		{UserID: user1.ID, Amount: 50.0},
		{UserID: user2.ID, Amount: 50.0},
	})
	_, _ = expenseUC.CreateExpense(ctx, "Movie", "entertainment", user1.ID, 50.0, []domain.Split{
		{UserID: user1.ID, Amount: 25.0},
		{UserID: user2.ID, Amount: 25.0},
	})

	// Get stats for user1
	stats, err := expenseUC.GetUserStats(ctx, user1.ID)
	if err != nil {
		t.Fatalf("Failed to get user stats: %v", err)
	}

	// Verify stats
	if stats.UserID != user1.ID {
		t.Errorf("Expected user ID %s, got: %s", user1.ID, stats.UserID)
	}
	if stats.TotalPaid != 150.0 {
		t.Errorf("Expected total paid 150.0, got: %v", stats.TotalPaid)
	}
	if stats.TotalOwed != 75.0 {
		t.Errorf("Expected total owed 75.0, got: %v", stats.TotalOwed)
	}
	if stats.NetBalance != 75.0 {
		t.Errorf("Expected net balance 75.0, got: %v", stats.NetBalance)
	}
	if stats.ExpenseCount != 2 {
		t.Errorf("Expected expense count 2, got: %v", stats.ExpenseCount)
	}

	// Check category breakdown
	if stats.ByCategory["food"] != 100.0 {
		t.Errorf("Expected food category 100.0, got: %v", stats.ByCategory["food"])
	}
	if stats.ByCategory["entertainment"] != 50.0 {
		t.Errorf("Expected entertainment category 50.0, got: %v", stats.ByCategory["entertainment"])
	}
}

func TestExpenseUseCase_GetUserStats_UserNotFound(t *testing.T) {
	expenseRepo := newMockExpenseRepository()
	userRepo := newMockUserRepository()
	expenseUC := NewExpenseUseCase(expenseRepo, userRepo)
	ctx := context.Background()

	// Try to get stats for non-existent user
	_, err := expenseUC.GetUserStats(ctx, "nonexistent")
	if err == nil {
		t.Fatal("Expected error for non-existent user, got nil")
	}
}

func TestExpenseUseCase_GetMonthlySummary(t *testing.T) {
	expenseRepo := newMockExpenseRepository()
	userRepo := newMockUserRepository()
	expenseUC := NewExpenseUseCase(expenseRepo, userRepo)
	ctx := context.Background()

	// Create test user
	user1 := &domain.User{ID: "user1", Name: "User1", Email: "user1@test.com"}
	_ = userRepo.Create(ctx, user1)

	// Create expenses with specific dates
	expense1, _ := expenseUC.CreateExpense(ctx, "Groceries", "food", user1.ID, 100.0, []domain.Split{{UserID: user1.ID, Amount: 100.0}})
	expense2, _ := expenseUC.CreateExpense(ctx, "Restaurant", "food", user1.ID, 50.0, []domain.Split{{UserID: user1.ID, Amount: 50.0}})
	expense3, _ := expenseUC.CreateExpense(ctx, "Movie", "entertainment", user1.ID, 30.0, []domain.Split{{UserID: user1.ID, Amount: 30.0}})

	// Set dates to November 2025
	expense1.Date = time.Date(2025, 11, 5, 0, 0, 0, 0, time.UTC)
	expense2.Date = time.Date(2025, 11, 15, 0, 0, 0, 0, time.UTC)
	expense3.Date = time.Date(2025, 11, 20, 0, 0, 0, 0, time.UTC)

	// Get monthly summary for November 2025
	summary, err := expenseUC.GetMonthlySummary(ctx, 2025, 11)
	if err != nil {
		t.Fatalf("Failed to get monthly summary: %v", err)
	}

	// Verify summary
	if summary.Year != 2025 {
		t.Errorf("Expected year 2025, got: %v", summary.Year)
	}
	if summary.Month != 11 {
		t.Errorf("Expected month 11, got: %v", summary.Month)
	}
	if summary.TotalExpenses != 180.0 {
		t.Errorf("Expected total expenses 180.0, got: %v", summary.TotalExpenses)
	}
	if summary.ExpenseCount != 3 {
		t.Errorf("Expected expense count 3, got: %v", summary.ExpenseCount)
	}

	// Check category breakdown
	if summary.ByCategory["food"] != 150.0 {
		t.Errorf("Expected food category 150.0, got: %v", summary.ByCategory["food"])
	}
	if summary.ByCategory["entertainment"] != 30.0 {
		t.Errorf("Expected entertainment category 30.0, got: %v", summary.ByCategory["entertainment"])
	}

	// Check top category
	if summary.TopCategory != "food" {
		t.Errorf("Expected top category 'food', got: %s", summary.TopCategory)
	}

	// Check average per day (180 / 30 days in November)
	if summary.AveragePerDay != 6.0 {
		t.Errorf("Expected average per day 6.0, got: %v", summary.AveragePerDay)
	}
}

func TestExpenseUseCase_GetMonthlySummary_InvalidMonth(t *testing.T) {
	expenseRepo := newMockExpenseRepository()
	userRepo := newMockUserRepository()
	expenseUC := NewExpenseUseCase(expenseRepo, userRepo)
	ctx := context.Background()

	// Try with invalid month
	_, err := expenseUC.GetMonthlySummary(ctx, 2025, 13)
	if err == nil {
		t.Fatal("Expected error for invalid month, got nil")
	}

	// Try with month 0
	_, err = expenseUC.GetMonthlySummary(ctx, 2025, 0)
	if err == nil {
		t.Fatal("Expected error for month 0, got nil")
	}
}
