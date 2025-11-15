package handler

import (
	"time"

	"github.com/pavanrkadave/homies/internal/domain"
)

// ToUserResponse converts a domain.User to UserResponse
func ToUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}
}

// ToUserResponses converts multiple users to response DTOs
func ToUserResponses(users []*domain.User) []UserResponse {
	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = ToUserResponse(user)
	}
	return responses
}

// ToExpenseResponse converts a domain.Expense to ExpenseResponse
func ToExpenseResponse(expense *domain.Expense) ExpenseResponse {
	splits := make([]SplitResponse, len(expense.Splits))
	for i, split := range expense.Splits {
		splits[i] = SplitResponse{
			UserId: split.UserID,
			Amount: split.Amount,
		}
	}

	return ExpenseResponse{
		ID:          expense.ID,
		Description: expense.Description,
		Amount:      expense.Amount,
		Category:    expense.Category,
		PaidBy:      expense.PaidBy,
		Date:        expense.Date,
		CreatedAt:   expense.CreatedAt,
		Splits:      splits,
	}
}

// ToExpenseResponses converts multiple expenses to response DTOs
func ToExpenseResponses(expenses []*domain.Expense) []ExpenseResponse {
	responses := make([]ExpenseResponse, len(expenses))
	for i, expense := range expenses {
		responses[i] = ToExpenseResponse(expense)
	}
	return responses
}
