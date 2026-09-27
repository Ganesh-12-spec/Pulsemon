package history

import "time"

type HealthCheckResult struct {
	Target string
	URl    string
	Time   time.Time
}

type Record struct {
	Target string
	Status string
	Time   time.Time
}

type History struct {
	Records []Record
}

func (h *History) Add(record Record) {
	h.Records = append(h.Records, record)
}
