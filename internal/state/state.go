package state

import "sync"

type Status string

const (
	UP   Status = "UP"
	DOWN Status = "DOWN"
)

type Monitor struct {
	mu            sync.RWMutex
	states        map[string]Status
	failureCounts map[string]int
}

func NewMonitor() *Monitor {
	return &Monitor{
		states:        make(map[string]Status),
		failureCounts: make(map[string]int),
	}
}

func (m *Monitor) Update(target string, err error) (Status, bool, bool, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	newStatus := UP

	if err != nil {
		newStatus = DOWN
		m.failureCounts[target]++
	} else {
		m.failureCounts[target] = 0
	}

	oldStatus, exists := m.states[target]

	m.states[target] = newStatus

	if !exists {
		thresholdReached := m.failureCounts[target] >= 3
		return newStatus, true, thresholdReached, false
	}

	changed := oldStatus != newStatus
	thresholdReached := m.failureCounts[target] == 3
	recovered := oldStatus == DOWN && newStatus == UP

	return newStatus, changed, thresholdReached, recovered
}

func (m *Monitor) Snapshot() map[string]Status {
	m.mu.RLock()
	defer m.mu.RUnlock()

	snapshot := make(map[string]Status)

	for target, status := range m.states {
		snapshot[target] = status
	}

	return snapshot
}
