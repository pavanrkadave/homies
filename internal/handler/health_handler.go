package handler

import (
	"database/sql"
	"net/http"

	"github.com/pavanrkadave/homies/pkg/response"
)

type HealthHandler struct {
	db *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	dbStatus := "healthy"
	if err := h.db.Ping(); err != nil {
		dbStatus = "unhealthy"
	}

	healthResponse := HealthResponse{
		Status:   "ok",
		Database: dbStatus,
	}

	response.RespondWithJSON(w, http.StatusOK, healthResponse)
}
