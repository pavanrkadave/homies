package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/pavanrkadave/homies/internal/domain"
	"github.com/pavanrkadave/homies/internal/repository"
)

var _ repository.ExpenseRepository = (*ExpensePostgresRepository)(nil)

type ExpensePostgresRepository struct {
	db *sql.DB
}

func NewExpensePostgresRepository(db *sql.DB) *ExpensePostgresRepository {
	return &ExpensePostgresRepository{db: db}
}

func (r *ExpensePostgresRepository) Create(ctx context.Context, expense *domain.Expense) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			log.Printf("failed to rollback transaction: %w", err)
		}
	}(tx)

	// Insert expense
	expenseQuery := `
		INSERT INTO expenses (id, description, amount, category, paid_by, date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = tx.ExecContext(ctx, expenseQuery,
		expense.ID,
		expense.Description,
		expense.Amount,
		expense.Category,
		expense.PaidBy,
		expense.Date,
		expense.CreatedAt,
		expense.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create expense: %w", err)
	}

	// Insert splits
	splitQuery := `
		INSERT INTO splits (expense_id, user_id, amount)
		VALUES ($1, $2, $3)
	`
	for _, split := range expense.Splits {
		_, err = tx.ExecContext(ctx, splitQuery,
			expense.ID,
			split.UserID,
			split.Amount,
		)
		if err != nil {
			return fmt.Errorf("failed to create split: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *ExpensePostgresRepository) GetByID(ctx context.Context, id string) (*domain.Expense, error) {
	expenseQuery := `
		SELECT id, description, amount, category, paid_by, date, created_at, updated_at
		FROM expenses
		WHERE id = $1
	`

	expense := &domain.Expense{}
	err := r.db.QueryRowContext(ctx, expenseQuery, id).Scan(
		&expense.ID,
		&expense.Description,
		&expense.Amount,
		&expense.Category,
		&expense.PaidBy,
		&expense.Date,
		&expense.CreatedAt,
		&expense.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("expense not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get expense: %w", err)
	}

	// Get splits for this expense
	splitsQuery := `
		SELECT expense_id, user_id, amount
		FROM splits
		WHERE expense_id = $1
	`

	rows, err := r.db.QueryContext(ctx, splitsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get splits: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}(rows)

	var splits []domain.Split
	for rows.Next() {
		var split domain.Split
		if err := rows.Scan(&split.ExpenseID, &split.UserID, &split.Amount); err != nil {
			return nil, fmt.Errorf("failed to scan split: %w", err)
		}
		splits = append(splits, split)
	}

	expense.Splits = splits
	return expense, nil
}

func (r *ExpensePostgresRepository) GetAll(ctx context.Context) ([]*domain.Expense, error) {
	allExpenseQuery := `SELECT id, description, amount, category, paid_by, date, created_at, updated_at FROM expenses`

	var expenses []*domain.Expense
	expenseRows, err := r.db.QueryContext(ctx, allExpenseQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(expenseRows)

	for expenseRows.Next() {
		expense := &domain.Expense{}
		if err := expenseRows.Scan(&expense.ID,
			&expense.Description,
			&expense.Amount,
			&expense.Category,
			&expense.PaidBy,
			&expense.Date,
			&expense.CreatedAt,
			&expense.UpdatedAt); err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}

	if err := expenseRows.Err(); err != nil {
		return nil, err
	}

	for _, expense := range expenses {
		splitsQuery := `SELECT expense_id, user_id, amount FROM splits WHERE expense_id = $1`
		rows, err := r.db.QueryContext(ctx, splitsQuery, expense.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get splits: %w", err)
		}
		var splits []domain.Split
		for rows.Next() {
			var split domain.Split
			if err := rows.Scan(&split.ExpenseID, &split.UserID, &split.Amount); err != nil {
				return nil, fmt.Errorf("failed to scan split: %w", err)
			}
			splits = append(splits, split)
		}
		expense.Splits = splits
	}
	return expenses, nil
}

func (r *ExpensePostgresRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.Expense, error) {
	// Get all unique expense IDs where user is involved
	expenseQuery := `
        SELECT DISTINCT e.id, e.description, e.amount, e.category, e.paid_by, e.date, e.created_at, e.updated_at
        FROM expenses e
        LEFT JOIN splits s ON e.id = s.expense_id
        WHERE e.paid_by = $1 OR s.user_id = $1
    `

	rows, err := r.db.QueryContext(ctx, expenseQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get expenses: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("failed to close rows: ", err)
		}
	}(rows)

	var expenses []*domain.Expense
	for rows.Next() {
		expense := &domain.Expense{}
		if err := rows.Scan(
			&expense.ID,
			&expense.Description,
			&expense.Amount,
			&expense.Category,
			&expense.PaidBy,
			&expense.Date,
			&expense.CreatedAt,
			&expense.UpdatedAt,
		); err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}

	// Get splits for each expense
	for _, expense := range expenses {
		splitsQuery := `SELECT expense_id, user_id, amount FROM splits WHERE expense_id = $1`
		splitRows, err := r.db.QueryContext(ctx, splitsQuery, expense.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get splits: %w", err)
		}

		var splits []domain.Split
		for splitRows.Next() {
			var split domain.Split
			if err := splitRows.Scan(&split.ExpenseID, &split.UserID, &split.Amount); err != nil {
				err := splitRows.Close()
				if err != nil {
					return nil, err
				}
				return nil, fmt.Errorf("failed to scan split: %w", err)
			}
			splits = append(splits, split)
		}
		err = splitRows.Close()
		if err != nil {
			return nil, err
		}
		expense.Splits = splits
	}

	return expenses, nil
}

func (r *ExpensePostgresRepository) Delete(ctx context.Context, id string) error {
	deleteExpenseQuery := `DELETE FROM expenses WHERE id = $1`
	_, err := r.db.ExecContext(ctx, deleteExpenseQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete expense: %w", err)
	}
	return nil
}
