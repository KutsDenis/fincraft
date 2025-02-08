package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes регистрирует маршруты
func (h *Handler) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	r.Get("/health", h.HealthCheck)

	return r
}
