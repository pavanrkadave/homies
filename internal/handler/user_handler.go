package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/pavanrkadave/homies/internal/usecase"
	"github.com/pavanrkadave/homies/pkg/errors"
)

type UserHandler struct {
	userUC usecase.UserUseCase
}

func NewUserHandler(userUC usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUC: userUC,
	}
}

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		errors.ResponseWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errors.ResponseWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.userUC.CreateUser(r.Context(), req.Name, req.Email)
	if err != nil {
		errors.ResponseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	userResponse := UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(userResponse)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errors.ResponseWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	users, err := h.userUC.GetAllUsers(r.Context())
	if err != nil {
		errors.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	usersResponse := make([]UserResponse, 0)
	for _, user := range users {
		usersResponse = append(usersResponse, UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(usersResponse)
}
