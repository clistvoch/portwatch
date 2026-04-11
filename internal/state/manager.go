package state

import (
	"time"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

// Manager coordinates loading, diffing, and saving snapshots.
type Manager struct {
	store *Store
}

// NewManager creates a Manager backed by the given Store.
func NewManager(store *Store) *Manager {
	return &Manager{store: store}
}

// Update loads the previous snapshot, diffs it against current ports,
// saves the new snapshot, and returns any detected changes.
func (m *Manager) Update(current []scanner.PortInfo) ([]monitor.Change, error) {
	prev, err := m.store.Load()
	if err != nil {
		return nil, err
	}

	changes := monitor.Diff(prev.Ports, current)

	next := Snapshot{
		Timestamp: time.Now(),
		Ports:     current,
	}
	if err := m.store.Save(next); err != nil {
		return nil, err
	}

	return changes, nil
}

// IsFirstRun reports whether no prior state exists on disk.
func (m *Manager) IsFirstRun() bool {
	return !m.store.Exists()
}
