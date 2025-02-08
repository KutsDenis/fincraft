package transport

import (
	"encoding/json"
	"github.com/KutsDenis/logzap"
	"go.uber.org/zap"
	"net/http"
)

// WriteSuccess отправляет успешный ответ.
func WriteSuccess(w http.ResponseWriter, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logzap.Error("failed to write response", zap.Error(err))
	}
}
