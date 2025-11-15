# HTTPie Test Commands

This document contains HTTPie commands for testing the Homies API endpoints.

## Prerequisites
```bash
# Install HTTPie (if not already installed)
brew install httpie

# Start the application
docker-compose up -d
```

## User Endpoints

### Create User
```bash
http POST localhost:3000/users name="John Doe" email="john@example.com"
```

### Get All Users
```bash
http GET localhost:3000/users
```

### Get User by ID
```bash
http GET localhost:3000/users id==USER_ID_HERE
```

### Update User (Phase 1 - NEW ✅)
```bash
http PUT localhost:3000/users id==USER_ID_HERE name="New Name" email="new@email.com"
```

## Using curl (alternative)

### Create User
```bash
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'
```

### Get User by ID
```bash
curl -X GET "http://localhost:3000/users?id=USER_ID_HERE"
```

### Update User
```bash
curl -X PUT "http://localhost:3000/users?id=USER_ID_HERE" \
  -H "Content-Type: application/json" \
  -d '{"name":"New Name","email":"new@email.com"}'
```

## Test Scenarios

### Scenario 1: Complete User Update Flow
```bash
# 1. Create a user
http POST localhost:3000/users name="Alice Smith" email="alice@example.com"
# Copy the returned ID

# 2. Get the user
http GET localhost:3000/users id==PASTE_ID_HERE

# 3. Update the user
http PUT localhost:3000/users id==PASTE_ID_HERE name="Alice Johnson" email="alice.j@example.com"

# 4. Verify the update
http GET localhost:3000/users id==PASTE_ID_HERE
```

### Scenario 2: Test Error Handling
```bash
# Test 404 - User not found
http GET localhost:3000/users id==nonexistent-id

# Test 409 - Email already exists
# First create two users
http POST localhost:3000/users name="User1" email="user1@test.com"
http POST localhost:3000/users name="User2" email="user2@test.com"

# Try to update User1 with User2's email (should fail with 409)
http PUT localhost:3000/users id==USER1_ID name="User1" email="user2@test.com"
```

## Expense Endpoints (Existing)

### Create Expense
```bash
http POST localhost:3000/expenses \
  description="Dinner" \
  amount:=100.00 \
  category="food" \
  paid_by="USER_ID" \
  splits:='[{"user_id":"USER_ID_1","amount":50.00},{"user_id":"USER_ID_2","amount":50.00}]'
```

### Get All Expenses
```bash
http GET localhost:3000/expenses
```

### Get Expense by ID
```bash
http GET localhost:3000/expenses id==EXPENSE_ID
```

### Get Expenses by User
```bash
http GET localhost:3000/expenses/user user_id==USER_ID
```

### Delete Expense
```bash
http DELETE localhost:3000/expenses id==EXPENSE_ID
```

### Get Balances
```bash
http GET localhost:3000/balances
```

## Health Check
```bash
http GET localhost:3000/health
```

## Tips

- **Pretty Print JSON**: HTTPie automatically pretty-prints JSON
- **Verbose Output**: Add `-v` flag for verbose output with headers
- **Save Response**: Pipe to file `http GET ... > response.json`
- **HTTP Status**: HTTPie shows status codes in color (green=success, red=error)

## Quick Test All New Endpoints
```bash
# Run this script to test all Phase 1 endpoints
USER_RESPONSE=$(http --body POST localhost:3000/users name="Test User" email="test@example.com")
USER_ID=$(echo $USER_RESPONSE | jq -r '.id')

echo "Created user with ID: $USER_ID"

# Get user by ID
http GET localhost:3000/users id==$USER_ID

# Update user
http PUT localhost:3000/users id==$USER_ID name="Updated Name" email="updated@example.com"

# Verify update
http GET localhost:3000/users id==$USER_ID

echo "✅ All tests completed!"
```

