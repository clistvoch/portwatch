package state_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "state.json")
}

func TestStore_SaveAndLoad(t *testing.T) {
	path := tempPath(t)
	store := state.NewStore(path)

	snap := state.Snapshot{
		Timestamp: time.Now().Truncate(time.Second),
		Ports: []scanner.PortInfo{
			{Port: 80, Proto: "tcp", State: "open"},
			{Port: 443, Proto: "tcp", State: "open"},
		},
	}

	if err := store.Save(snap); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(loaded.Ports) != len(snap.Ports) {
		t.Errorf("got %d ports, want %d", len(loaded.Ports), len(snap.Ports))
	}
	if !loaded.Timestamp.Equal(snap.Timestamp) {
		t.Errorf("timestamp mismatch: got %v, want %v", loaded.Timestamp, snap.Timestamp)
	}
}

func TestStore_Load_MissingFile(t *testing.T) {
	store := state.NewStore(filepath.Join(t.TempDir(), "nonexistent.json"))
	snap, err := store.Load()
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(snap.Ports) != 0 {
		t.Errorf("expected empty ports, got %v", snap.Ports)
	}
}

func TestStore_Exists(t *testing.T) {
	path := tempPath(t)
	store := state.NewStore(path)

	if store.Exists() {
		t.Error("Exists should be false before Save")
	}

	_ = store.Save(state.Snapshot{Timestamp: time.Now()})

	if !store.Exists() {
		t.Error("Exists should be true after Save")
	}
}

func TestStore_Save_InvalidPath(t *testing.T) {
	store := state.NewStore("/nonexistent_dir/state.json")
	err := store.Save(state.Snapshot{})
	if err == nil {
		t.Error("expected error for invalid path")
		os.Remove("/nonexistent_dir/state.json")
	}
}
