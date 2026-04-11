package state

import (
	"encoding/json"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Snapshot represents a saved port scan result with metadata.
type Snapshot struct {
	Timestamp time.Time           `json:"timestamp"`
	Ports     []scanner.PortInfo  `json:"ports"`
}

// Store handles persistence of port snapshots to disk.
type Store struct {
	path string
}

// NewStore creates a Store that reads/writes to the given file path.
func NewStore(path string) *Store {
	return &Store{path: path}
}

// Save writes the snapshot to disk as JSON, overwriting any existing file.
func (s *Store) Save(snap Snapshot) error {
	f, err := os.Create(s.path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(snap)
}

// Load reads the most recent snapshot from disk.
// Returns an empty Snapshot and no error if the file does not exist.
func (s *Store) Load() (Snapshot, error) {
	f, err := os.Open(s.path)
	if os.IsNotExist(err) {
		return Snapshot{}, nil
	}
	if err != nil {
		return Snapshot{}, err
	}
	defer f.Close()

	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return Snapshot{}, err
	}
	return snap, nil
}

// Exists reports whether a persisted snapshot file is present.
func (s *Store) Exists() bool {
	_, err := os.Stat(s.path)
	return err == nil
}
