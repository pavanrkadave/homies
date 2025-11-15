-- Create splits table
CREATE TABLE IF NOT EXISTS splits (
                                      id SERIAL PRIMARY KEY,
                                      expense_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL CHECK (amount > 0),
    FOREIGN KEY (expense_id) REFERENCES expenses(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(expense_id, user_id)
    );

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_splits_expense_id ON splits(expense_id);
CREATE INDEX IF NOT EXISTS idx_splits_user_id ON splits(user_id);