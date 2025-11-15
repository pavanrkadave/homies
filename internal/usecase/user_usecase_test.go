package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/pavanrkadave/homies/internal/domain"
)

type mockUserRepository struct {
	users map[string]*domain.User
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *mockUserRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	var users []*domain.User
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func (m *mockUserRepository) Update(ctx context.Context, user *domain.User) error {
	if _, ok := m.users[user.ID]; !ok {
		return errors.New("user not found")
	}
	m.users[user.ID] = user
	return nil
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[string]*domain.User),
	}
}

func Test_userUseCase_CreateUser(t *testing.T) {

	repo := newMockUserRepository()
	userUseCase := NewUserUseCase(repo)
	ctx := context.Background()

	user, err := userUseCase.CreateUser(ctx, "Pavan", "pavan@email.com")
	if err != nil {
		t.Fatalf("userUseCase.CreateUser returned error: %v", err)
	}

	if user.Name != "Pavan" {
		t.Fatalf("userUseCase.CreateUser returned user.Name: %v but expected: %v", user.Name, "Pavan")
	}

	if user.CreatedAt.IsZero() {
		t.Fatalf("userUseCase.CreateUser should have createdAt set")
	}
}

func Test_userUseCase_CreateUser_ValidationError(t *testing.T) {
	repo := newMockUserRepository()
	userUseCase := NewUserUseCase(repo)
	ctx := context.Background()

	_, err := userUseCase.CreateUser(ctx, "", "pavan@email.com")
	if err == nil {
		t.Fatal("Expected validation error for empty name, got nil")
	}

	_, err = userUseCase.CreateUser(ctx, "Pavan", "")
	if err == nil {
		t.Fatal("Expected validation error for empty email, got nil")
	}
}

func Test_userUseCase_UpdateUser(t *testing.T) {
	repo := newMockUserRepository()
	userUseCase := NewUserUseCase(repo)
	ctx := context.Background()

	// First, create a user
	user, err := userUseCase.CreateUser(ctx, "Pavan", "pavan@email.com")
	if err != nil {
		t.Fatalf("userUseCase.CreateUser returned error: %v", err)
	}

	// Update the user
	updatedUser, err := userUseCase.UpdateUser(ctx, user.ID, "Pavan Reddy", "pavan.reddy@email.com")
	if err != nil {
		t.Fatalf("userUseCase.UpdateUser returned error: %v", err)
	}

	if updatedUser.Name != "Pavan Reddy" {
		t.Fatalf("Expected name 'Pavan Reddy', got: %v", updatedUser.Name)
	}

	if updatedUser.Email != "pavan.reddy@email.com" {
		t.Fatalf("Expected email 'pavan.reddy@email.com', got: %v", updatedUser.Email)
	}

	if updatedUser.UpdatedAt.IsZero() {
		t.Fatal("UpdatedAt should be set")
	}
}

func Test_userUseCase_UpdateUser_UserNotFound(t *testing.T) {
	repo := newMockUserRepository()
	userUseCase := NewUserUseCase(repo)
	ctx := context.Background()

	// Try to update a user that doesn't exist
	_, err := userUseCase.UpdateUser(ctx, "nonexistent", "Ghost", "ghost@email.com")
	if err == nil {
		t.Fatal("Expected error when updating non-existent user, got nil")
	}
}

func Test_userUseCase_UpdateUser_EmailAlreadyExists(t *testing.T) {
	repo := newMockUserRepository()
	userUseCase := NewUserUseCase(repo)
	ctx := context.Background()

	// Create two users
	user1, _ := userUseCase.CreateUser(ctx, "User1", "user1@email.com")
	_, _ = userUseCase.CreateUser(ctx, "User2", "user2@email.com")

	// Try to update user1 with user2's email
	_, err := userUseCase.UpdateUser(ctx, user1.ID, "User1", "user2@email.com")
	if err == nil {
		t.Fatal("Expected error when updating user with existing email, got nil")
	}

	if err != domain.ErrEmailAlreadyExists {
		t.Fatalf("Expected ErrEmailAlreadyExists, got: %v", err)
	}
}
