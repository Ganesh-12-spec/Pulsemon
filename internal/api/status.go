package api

import (
	"encoding/json"
	"net/http"

	"github.com/Ganesh-12-spec/pulsemon/internal/state"
)

type StatusResponse struct {
	Targets map[string]state.Status `json:"targets"`
}

func StatusHandler(monitor *state.Monitor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := StatusResponse{
			Targets: monitor.Snapshot(),
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode status", http.StatusInternalServerError)
			return
		}
	}
}
