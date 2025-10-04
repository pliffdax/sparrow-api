package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/pliffdax/sparrow-api/internal/util"
)

type Health struct {
	Status  int    `json:"status"`
	TimeUTC string `json:"time"`
	Version string `json:"version"`
}

func HealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := Health{
			Status:  http.StatusOK,
			TimeUTC: time.Now().UTC().Format(time.RFC3339),
			Version: util.Getenv("APP_VERSION", "0.1.0"),
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
