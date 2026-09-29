package state

import (
	"errors"
	"testing"
)

var errTest = errors.New("test failure")

func TestMonitorInitialState(t *testing.T) {
	monitor := NewMonitor()

	status, changed, thresholdReached, recovered := monitor.Update("Google", nil)

	if status != UP {
		t.Fatalf("expected UP, got %s", status)
	}

	if !changed {
		t.Fatal("expected initial state to be changed")
	}

	if thresholdReached {
		t.Fatal("did not expect threshold alert")
	}

	if recovered {
		t.Fatal("did not expect recovery")
	}
}

func TestFailureThreshold(t *testing.T) {
	monitor := NewMonitor()

	monitor.Update("Google", errTest)
	monitor.Update("Google", errTest)

	_, _, thresholdReached, _ := monitor.Update("Google", errTest)

	if !thresholdReached {
		t.Fatal("expected failure threshold after 3 consecutive failures")
	}
}

func TestRecovery(t *testing.T) {
	monitor := NewMonitor()

	monitor.Update("Google", errTest)
	_, _, _, recovered := monitor.Update("Google", nil)

	if !recovered {
		t.Fatal("expected recovery after DOWN → UP")
	}
}

func TestFailureCounterResets(t *testing.T) {
	monitor := NewMonitor()

	monitor.Update("Google", errTest)
	monitor.Update("Google", errTest)
	monitor.Update("Google", nil)

	_, _, thresholdReached, _ := monitor.Update("Google", errTest)

	if thresholdReached {
		t.Fatal("expected failure counter to reset after recovery")
	}
}
