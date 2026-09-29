package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Ganesh-12-spec/pulsemon/internal/state"
)

func TestStatusHandler(t *testing.T) {
	monitor := state.NewMonitor()
	monitor.Update("Google", nil)

	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	recorder := httptest.NewRecorder()

	handler := StatusHandler(monitor)
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	body := recorder.Body.String()

	if !strings.Contains(body, "Google") {
		t.Fatalf("expected response to contain Google, got %s", body)
	}

	if !strings.Contains(body, "UP") {
		t.Fatalf("expected response to contain UP, got %s", body)
	}
}
