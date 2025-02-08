package handlers

import (
	"fincraft/internal/transport"
	"net/http"
)

// HealthCheck проверяет работоспособность сервиса.
func (h *Handler) HealthCheck(w http.ResponseWriter, _ *http.Request) {
	response := struct {
		Status string `json:"status"`
	}{
		Status: "Ok",
	}

	transport.WriteSuccess(w, response)
}
