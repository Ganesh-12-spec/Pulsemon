package history

import "time"

type Record struct {
	Target  string
	Status  string
	Time    time.Time
	Latency time.Duration
}

type History struct {
	Records []Record
}

func (h *History) Add(record Record) {
	h.Records = append(h.Records, record)
}

// Uptime calculates the percentage of successful health checks.
func (h *History) Uptime(target string) float64 {
	if len(h.Records) == 0 {
		return 0
	}

	total := 0
	up := 0

	for _, record := range h.Records {
		if record.Target != target {
			continue
		}

		total++

		if record.Status == "UP" {
			up++
		}
	}

	if total == 0 {
		return 0
	}

	return float64(up) / float64(total) * 100
}

// AverageLatency calculates the average response time for a target.
func (h *History) AverageLatency(target string) time.Duration {
	var totalLatency time.Duration
	count := 0

	for _, record := range h.Records {
		if record.Target != target {
			continue
		}

		totalLatency += record.Latency
		count++
	}

	if count == 0 {
		return 0
	}

	return totalLatency / time.Duration(count)
}
