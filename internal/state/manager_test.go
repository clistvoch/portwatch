package state_test

import (
	"testing"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

func TestManager_IsFirstRun_True(t *testing.T) {
	store := state.NewStore(tempPath(t))
	mgr := state.NewManager(store)

	if !mgr.IsFirstRun() {
		t.Error("expected IsFirstRun to be true when no state file exists")
	}
}

func TestManager_IsFirstRun_False(t *testing.T) {
	store := state.NewStore(tempPath(t))
	mgr := state.NewManager(store)

	_, _ = mgr.Update([]scanner.PortInfo{})

	if mgr.IsFirstRun() {
		t.Error("expected IsFirstRun to be false after first Update")
	}
}

func TestManager_Update_DetectsChanges(t *testing.T) {
	store := state.NewStore(tempPath(t))
	mgr := state.NewManager(store)

	initial := []scanner.PortInfo{
		{Port: 22, Proto: "tcp", State: "open"},
	}
	_, err := mgr.Update(initial)
	if err != nil {
		t.Fatalf("first Update: %v", err)
	}

	next := []scanner.PortInfo{
		{Port: 22, Proto: "tcp", State: "open"},
		{Port: 8080, Proto: "tcp", State: "open"},
	}
	changes, err := mgr.Update(next)
	if err != nil {
		t.Fatalf("second Update: %v", err)
	}

	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != monitor.Opened || changes[0].Port.Port != 8080 {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestManager_Update_EmptyToEmpty(t *testing.T) {
	store := state.NewStore(tempPath(t))
	mgr := state.NewManager(store)

	changes, err := mgr.Update([]scanner.PortInfo{})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %v", changes)
	}
}
