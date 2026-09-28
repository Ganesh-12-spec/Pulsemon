package state

type Status string

const (
	UP   Status = "UP"
	DOWN Status = "DOWN"
)

type Monitor struct {
	states        map[string]Status
	failureCounts map[string]int
}

func NewMonitor() *Monitor {
	return &Monitor{
		states:        make(map[string]Status),
		failureCounts: make(map[string]int),
	}
}

func (m *Monitor) Update(target string, err error) (Status, bool, bool) {
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
		return newStatus, true, thresholdReached
	}

	changed := oldStatus != newStatus
	thresholdReached := m.failureCounts[target] == 3

	return newStatus, changed, thresholdReached
}
