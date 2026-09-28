package state

type Status string

const (
	UP   Status = "UP"
	DOWN Status = "DOWN"
)

type Monitor struct {
	states map[string]Status
}

func NewMonitor() *Monitor {
	return &Monitor{
		states: make(map[string]Status),
	}
}

func (m *Monitor) Update(target string, err error) (Status, bool) {
	newStatus := UP

	if err != nil {
		newStatus = DOWN
	}

	oldStatus, exists := m.states[target]

	m.states[target] = newStatus

	if !exists {
		return newStatus, true
	}

	return newStatus, oldStatus != newStatus
}
