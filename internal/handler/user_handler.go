package handler

import (
	"encoding/json"
	"net/http"

	"github.com/pavanrkadave/homies/internal/usecase"
	"github.com/pavanrkadave/homies/pkg/response"
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
		response.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.userUC.CreateUser(r.Context(), req.Name, req.Email)
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, ToUserResponse(user))
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	users, err := h.userUC.GetAllUsers(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusOK, ToUserResponses(users))
}
