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
