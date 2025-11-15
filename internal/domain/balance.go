package domain

type Balance struct {
	UserID string  `json:"user_id"`
	Amount float64 `json:"amount"`
}

type Settlement struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

type BalanceSummary struct {
	Balances    []Balance    `json:"balances"`
	Settlements []Settlement `json:"settlements"`
}
