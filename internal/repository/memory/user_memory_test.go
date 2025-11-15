package memory

import (
	"context"
	"testing"
	"time"

	"github.com/pavanrkadave/homies/internal/domain"
)

func TestUserMemoryRepository_CreateAndGetByID(t *testing.T) {
	repo := NewUserMemoryRepository()
	ctx := context.Background()

	createUser := &domain.User{
		ID:        "1",
		Name:      "test",
		Email:     "test@email.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Create(ctx, createUser)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	retrievedUser, err := repo.GetByID(ctx, "1")

	if err != nil {
		t.Fatalf("GetByID() failed: %v", err)
	}

	if retrievedUser.ID != createUser.ID {
		t.Errorf("Expected ID '%s', got '%s'", createUser.ID, retrievedUser.ID)
	}
	if retrievedUser.Name != createUser.Name {
		t.Errorf("Expected name '%s', got '%s'", createUser.Name, retrievedUser.Name)
	}
	if retrievedUser.Email != createUser.Email {
		t.Errorf("Expected email '%s', got '%s'", createUser.Email, retrievedUser.Email)
	}
}

func TestUserMemoryRepository_CreateAndGetByEmail(t *testing.T) {
	repo := NewUserMemoryRepository()
	ctx := context.Background()

	createUser := &domain.User{
		ID:        "1",
		Name:      "test",
		Email:     "test@email.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Create(ctx, createUser)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	retrievedUser, err := repo.GetByEmail(ctx, "test@email.com")

	if err != nil {
		t.Fatalf("GetByID() failed: %v", err)
	}

	if retrievedUser.ID != createUser.ID {
		t.Errorf("Expected ID '%s', got '%s'", createUser.ID, retrievedUser.ID)
	}
	if retrievedUser.Name != createUser.Name {
		t.Errorf("Expected name '%s', got '%s'", createUser.Name, retrievedUser.Name)
	}
	if retrievedUser.Email != createUser.Email {
		t.Errorf("Expected email '%s', got '%s'", createUser.Email, retrievedUser.Email)
	}
}

func TestUserMemoryRepository_CreateAndGetAll(t *testing.T) {
	repo := NewUserMemoryRepository()
	ctx := context.Background()

	createUsers := []*domain.User{
		{ID: "1", Name: "test1", Email: "test1@email.com", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "2", Name: "test2", Email: "test2@email.com", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	for _, user := range createUsers {
		err := repo.Create(ctx, user)
		if err != nil {
			t.Errorf("Create() failed: %v", err)
		}
	}

	retrievedUsers, err := repo.GetAll(ctx)

	if err != nil {
		t.Fatalf("GetByID() failed: %v", err)
	}

	if len(retrievedUsers) != len(createUsers) {
		t.Fatalf("Expected %d users, got %d", len(createUsers), len(retrievedUsers))
	}
}

func TestUserMemoryRepository_CreateAndGetByIDNotFound(t *testing.T) {
	repo := NewUserMemoryRepository()
	ctx := context.Background()

	retrievedUser, err := repo.GetByID(ctx, "nonexistent")

	if err == nil {
		t.Fatal("Expected error for non-existent user, got nil")
	}

	if retrievedUser != nil {
		t.Fatalf("Expected nil user, got a user")
	}
}

func TestUserMemoryRepository_CreateAndGetByEmailNotFound(t *testing.T) {
	repo := NewUserMemoryRepository()
	ctx := context.Background()

	retrievedUser, err := repo.GetByEmail(ctx, "nonexistent")

	if err == nil {
		t.Fatal("Expected error for non-existent user, got nil")
	}

	if retrievedUser != nil {
		t.Fatalf("Expected nil user, got a user")
	}
}

func TestUserMemoryRepository_CreateAndGetAllNil(t *testing.T) {
	repo := NewUserMemoryRepository()
	ctx := context.Background()

	users, err := repo.GetAll(ctx)

	if err != nil {
		t.Fatalf("GetAll() failed: %v", err)
	}

	if len(users) != 0 {
		t.Fatalf("Expected 0 users, got : %d", len(users))
	}

}
