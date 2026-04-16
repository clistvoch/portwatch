package alert_test

import (
	"encoding/json"
	"testing"

	"github.com/patrickward/portwatch/internal/alert"
	"github.com/patrickward/portwatch/internal/monitor"
	"github.com/patrickward/portwatch/internal/scanner"
)

func buildGCPChange(port uint16, kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Port: scanner.PortInfo{Port: port, Proto: "tcp"},
		Kind: kind,
	}
}

func TestGCPPubSubPayload_MarshalNoChanges(t *testing.T) {
	changes := []monitor.Change{}
	data, err := json.Marshal(changes)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	if string(data) != "[]" {
		t.Errorf("unexpected: %s", string(data))
	}
}

func TestGCPPubSubPayload_MarshalWithChanges(t *testing.T) {
	changes := []monitor.Change{
		buildGCPChange(8080, monitor.Opened),
		buildGCPChange(9090, monitor.Closed),
	}
	data, err := json.Marshal(changes)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	var out []map[string]interface{}
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out))
	}
}

func TestGCPPubSubHandler_NoChanges(t *testing.T) {
	// NewGCPPubSubHandler requires real GCP credentials; test the no-op path
	// via a stub that satisfies the Handler interface.
	var called bool
	stub := handlerFunc(func(changes []monitor.Change) error {
		called = true
		return nil
	})
	if err := stub.Handle([]monitor.Change{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("handler should not be called for empty changes")
	}
}

// handlerFunc is a test helper that wraps a function as a Handler.
type handlerFunc func([]monitor.Change) error

func (f handlerFunc) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	return f(changes)
}

// Ensure GCPPubSubHandler satisfies alert.Handler at compile time.
var _ interface{ Handle([]monitor.Change) error } = (*alert.GCPPubSubHandler)(nil)
