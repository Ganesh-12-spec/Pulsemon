package history

import (
	"testing"
	"time"
)

func TestUptime(t *testing.T) {
	h := History{}

	h.Add(Record{
		Target: "Google",
		Status: "UP",
	})

	h.Add(Record{
		Target: "Google",
		Status: "UP",
	})

	h.Add(Record{
		Target: "Google",
		Status: "DOWN",
	})

	uptime := h.Uptime("Google")

	if uptime != 66.66666666666666 {
		t.Fatalf("expected approximately 66.67%% uptime, got %f", uptime)
	}
}

func TestAverageLatency(t *testing.T) {
	h := History{}

	h.Add(Record{
		Target:  "Google",
		Latency: 100 * time.Millisecond,
	})

	h.Add(Record{
		Target:  "Google",
		Latency: 300 * time.Millisecond,
	})

	average := h.AverageLatency("Google")

	if average != 200*time.Millisecond {
		t.Fatalf("expected 200ms, got %v", average)
	}
}
