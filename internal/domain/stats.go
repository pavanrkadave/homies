package domain

// UserStats represents spending statistics for a user
type UserStats struct {
	UserID       string             `json:"user_id"`
	TotalPaid    float64            `json:"total_paid"`
	TotalOwed    float64            `json:"total_owed"`
	NetBalance   float64            `json:"net_balance"`
	ExpenseCount int                `json:"expense_count"`
	ByCategory   map[string]float64 `json:"by_category"`
}

// MonthlySummary represents expense summary for a specific month
type MonthlySummary struct {
	Year          int                `json:"year"`
	Month         int                `json:"month"`
	TotalExpenses float64            `json:"total_expenses"`
	ExpenseCount  int                `json:"expense_count"`
	ByCategory    map[string]float64 `json:"by_category"`
	TopCategory   string             `json:"top_category"`
	AveragePerDay float64            `json:"average_per_day"`
}
