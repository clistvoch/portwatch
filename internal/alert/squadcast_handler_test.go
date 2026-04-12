package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

func buildSquadcastChange(port int, kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Port: scanner.PortInfo{Port: port, Proto: "tcp"},
		Kind: kind,
	}
}

func TestSquadcastHandler_NoChanges(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	h := NewSquadcastHandler(ts.URL, "production", 5)
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty changes")
	}
}

func TestSquadcastHandler_SendsPayload(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := NewSquadcastHandler(ts.URL, "staging", 5)
	changes := []monitor.Change{
		buildSquadcastChange(8080, monitor.Opened),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["status"] != "trigger" {
		t.Errorf("expected status=trigger, got %v", received["status"])
	}
	tags, ok := received["tags"].(map[string]interface{})
	if !ok {
		t.Fatal("expected tags map")
	}
	if tags["environment"] != "staging" {
		t.Errorf("expected environment=staging, got %v", tags["environment"])
	}
}

func TestSquadcastHandler_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := NewSquadcastHandler(ts.URL, "production", 5)
	changes := []monitor.Change{buildSquadcastChange(9090, monitor.Closed)}
	if err := h.Handle(changes); err == nil {
		t.Error("expected error on 500 response")
	}
}
