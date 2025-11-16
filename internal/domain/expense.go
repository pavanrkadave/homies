package domain

import (
	"errors"
	"math"
	"time"
)

type Expense struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	PaidBy      string    `json:"paid_by"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Splits      []Split   `json:"splits"`
}

type Split struct {
	ExpenseID string  `json:"expense_id"`
	UserID    string  `json:"user_id"`
	Amount    float64 `json:"amount"`
}

func (e *Expense) Validate() error {
	if e.Description == "" {
		return errors.New("expense description is required")
	}
	if e.Amount <= 0 {
		return errors.New("expense amount must be greater than zero")
	}
	if e.PaidBy == "" {
		return errors.New("paidBy is required")
	}
	if len(e.Splits) == 0 {
		return errors.New("at least one split is required")
	}

	var sum float64
	for _, split := range e.Splits {
		sum += split.Amount
	}

	if math.Abs(sum-e.Amount) > 0.01 {
		return errors.New("sum of splits must equal the expense amount")
	}

	return nil
}

func (e *Expense) Update(description, category string, amount float64, splits []Split) error {
	if description != "" {
		e.Description = description
	}
	if category != "" {
		e.Category = category
	}
	if amount > 0 {
		e.Amount = amount
	}
	if len(splits) > 0 {
		e.Splits = splits
	}
	e.UpdatedAt = time.Now()

	return e.Validate()
}
