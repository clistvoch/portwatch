package monitor

import (
	"fmt"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Change represents a detected port state change.
type Change struct {
	Port   int
	Status string // "opened" or "closed"
	At     time.Time
}

func (c Change) String() string {
	return fmt.Sprintf("[%s] port %d %s", c.At.Format(time.RFC3339), c.Port, c.Status)
}

// Monitor watches a port range and reports changes between scans.
type Monitor struct {
	scanner  *scanner.Scanner
	previous map[int]bool
	Interval time.Duration
}

// New creates a Monitor using the provided Scanner.
func New(s *scanner.Scanner, interval time.Duration) *Monitor {
	return &Monitor{
		scanner:  s,
		previous: make(map[int]bool),
		Interval: interval,
	}
}

// Diff performs a single scan and returns any changes since the last scan.
func (m *Monitor) Diff() ([]Change, error) {
	ports, err := m.scanner.Scan()
	if err != nil {
		return nil, fmt.Errorf("scan failed: %w", err)
	}

	current := make(map[int]bool, len(ports))
	for _, p := range ports {
		current[p.Port] = true
	}

	var changes []Change
	now := time.Now()

	for port := range current {
		if !m.previous[port] {
			changes = append(changes, Change{Port: port, Status: "opened", At: now})
		}
	}
	for port := range m.previous {
		if !current[port] {
			changes = append(changes, Change{Port: port, Status: "closed", At: now})
		}
	}

	m.previous = current
	return changes, nil
}

// OpenPorts returns the set of ports that are currently considered open
// based on the most recent scan state.
func (m *Monitor) OpenPorts() []int {
	ports := make([]int, 0, len(m.previous))
	for port := range m.previous {
		ports = append(ports, port)
	}
	return ports
}

// Run starts the monitoring loop, sending changes to the returned channel.
// Close the done channel to stop.
func (m *Monitor) Run(done <-chan struct{}) <-chan Change {
	ch := make(chan Change, 16)
	go func() {
		defer close(ch)
		// Seed initial state without emitting changes.
		_, _ = m.Diff()
		ticker := time.NewTicker(m.Interval)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				changes, err := m.Diff()
				if err != nil {
					continue
				}
				for _, c := range changes {
					ch <- c
				}
			}
		}
	}()
	return ch
}
