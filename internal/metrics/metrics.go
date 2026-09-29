package metrics

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type Metrics struct {
	mu               sync.RWMutex
	targetStatus     map[string]int
	uptimePercent    map[string]float64
	averageLatencyMs map[string]float64
}

func New() *Metrics {
	return &Metrics{
		targetStatus:     make(map[string]int),
		uptimePercent:    make(map[string]float64),
		averageLatencyMs: make(map[string]float64),
	}
}

func (m *Metrics) Update(
	target string,
	status string,
	uptime float64,
	averageLatencyMs float64,
) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if status == "UP" {
		m.targetStatus[target] = 1
	} else {
		m.targetStatus[target] = 0
	}

	m.uptimePercent[target] = uptime
	m.averageLatencyMs[target] = averageLatencyMs
}

func (m *Metrics) Handler(w http.ResponseWriter, r *http.Request) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	w.Header().Set("Content-Type", "text/plain; version=0.0.4")

	var output strings.Builder

	for target, status := range m.targetStatus {
		fmt.Fprintf(
			&output,
			"pulsemon_target_up{target=\"%s\"} %d\n",
			target,
			status,
		)
	}

	for target, uptime := range m.uptimePercent {
		fmt.Fprintf(
			&output,
			"pulsemon_uptime_percent{target=\"%s\"} %.2f\n",
			target,
			uptime,
		)
	}

	for target, latency := range m.averageLatencyMs {
		fmt.Fprintf(
			&output,
			"pulsemon_average_latency_ms{target=\"%s\"} %.2f\n",
			target,
			latency,
		)
	}

	_, _ = w.Write([]byte(output.String()))
}
