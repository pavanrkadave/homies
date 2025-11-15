package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/pavanrkadave/homies/pkg/errors"
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
		errors.ResponseWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	dbStatus := "healthy"
	if err := h.db.Ping(); err != nil {
		dbStatus = "unhealthy"
	}

	response := HealthResponse{
		Status:   "ok",
		Database: dbStatus,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
