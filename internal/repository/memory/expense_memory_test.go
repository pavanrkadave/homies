package memory

import (
	"context"
	"testing"
	"time"

	"github.com/pavanrkadave/homies/internal/domain"
)

func TestExpenseMemoryRepository_CreateAndGetByID(t *testing.T) {
	repo := NewExpenseMemoryRepository()
	ctx := context.Background()

	createExpense := &domain.Expense{
		ID:          "1",
		Description: "Subway",
		Amount:      13.50,
		Category:    "Food",
		PaidBy:      "1",
		Date:        time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Splits:      []domain.Split{},
	}

	err := repo.Create(ctx, createExpense)
	if err != nil {
		t.Fatalf("Unexpected error creating expense: %s", err)
	}

	expense, err := repo.GetByID(ctx, createExpense.ID)
	if err != nil {
		t.Fatalf("Unexpected error getting expense: %s", err)
	}
	if expense.ID != createExpense.ID {
		t.Fatalf("Expense ID does not match")
	}
}

func TestExpenseMemoryRepository_CreateAndGetByUserID(t *testing.T) {
	repo := NewExpenseMemoryRepository()
	ctx := context.Background()

	expenses := []domain.Expense{
		{ID: "1", Description: "Subway", Amount: 13.50, Category: "Food", PaidBy: "1", Date: time.Now(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Splits: []domain.Split{}},
		{ID: "2", Description: "KFC", Amount: 15.50, Category: "Food", PaidBy: "2", Date: time.Now(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Splits: []domain.Split{}},
		{ID: "3", Description: "PizzaHut", Amount: 14, Category: "Food", PaidBy: "2", Date: time.Now(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Splits: []domain.Split{
			{ExpenseID: "3", UserID: "1", Amount: 7},
			{ExpenseID: "3", UserID: "2", Amount: 7},
		}},
	}

	for _, expense := range expenses {
		err := repo.Create(ctx, &expense)
		if err != nil {
			t.Fatalf("should create expense but received error %+v", err)
		}
	}

	retrievedExpensesUser1, err := repo.GetByUserID(ctx, "1")
	if err != nil {
		t.Fatalf("Unexpected error getting expense: %s", err)
	}
	if len(retrievedExpensesUser1) != 2 {
		t.Fatalf("Expected 2 expenses for user 1, got %d", len(retrievedExpensesUser1))
	}

	retrievedExpensesUser2, err := repo.GetByUserID(ctx, "2")
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if len(retrievedExpensesUser2) != 2 {
		t.Fatalf("Expected 2 expenses for user 2, got %d", len(retrievedExpensesUser1))
	}
}

func TestExpenseMemoryRepository_CreateAndDelete(t *testing.T) {
	repo := NewExpenseMemoryRepository()
	ctx := context.Background()

	expenses := []domain.Expense{
		{ID: "1", Description: "Subway", Amount: 13.50, Category: "Food", PaidBy: "1", Date: time.Now(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Splits: []domain.Split{}},
		{ID: "2", Description: "KFC", Amount: 15.50, Category: "Food", PaidBy: "2", Date: time.Now(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Splits: []domain.Split{}},
		{ID: "3", Description: "PizzaHut", Amount: 14, Category: "Food", PaidBy: "2", Date: time.Now(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Splits: []domain.Split{
			{ExpenseID: "3", UserID: "1", Amount: 7},
			{ExpenseID: "3", UserID: "2", Amount: 7},
		}},
	}

	for _, expense := range expenses {
		err := repo.Create(ctx, &expense)
		if err != nil {
			t.Fatalf("should create expense but received error %+v", err)
		}
	}

	err := repo.Delete(ctx, "3")
	if err != nil {
		t.Fatalf("Unexpected error getting expense: %s", err)
	}

	retrievedExpenses, err := repo.GetAll(ctx)
	if err != nil {
		t.Fatalf("Unexpected error getting expenses: %s", err)
	}

	if len(retrievedExpenses) != 2 {
		t.Fatalf("Length of expenses should be 2 but was %d", len(retrievedExpenses))
	}

}
