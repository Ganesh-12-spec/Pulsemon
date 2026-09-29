package webhook

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSend(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("expected JSON content type, got %s", r.Header.Get("Content-Type"))
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}

		if !strings.Contains(string(body), "Google") {
			t.Fatalf("expected body to contain Google, got %s", body)
		}

		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	err := Send(server.URL, Payload{
		Target: "Google",
		Status: "DOWN",
		Event:  "failure_threshold",
	})

	if err != nil {
		t.Fatalf("expected webhook request to succeed, got %v", err)
	}
}
