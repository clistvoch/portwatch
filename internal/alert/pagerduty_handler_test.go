package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/monitor"
)

func buildPDChange(port int, kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Port: monitor.PortInfo{Port: port, Proto: "tcp"},
		Kind: kind,
	}
}

func TestPagerDutyHandler_NoChanges(t *testing.T) {
	h := NewPagerDutyHandler("test-key")
	if err := h.Handle(nil); err != nil {
		t.Fatalf("expected no error for empty changes, got %v", err)
	}
}

func TestPagerDutyHandler_SendsPayload(t *testing.T) {
	var received pdPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	h := NewPagerDutyHandler("my-routing-key")
	h.apiURL = ts.URL

	changes := []monitor.Change{
		buildPDChange(8080, monitor.Opened),
		buildPDChange(9090, monitor.Closed),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.RoutingKey != "my-routing-key" {
		t.Errorf("routing key = %q, want %q", received.RoutingKey, "my-routing-key")
	}
	if received.EventAction != "trigger" {
		t.Errorf("event_action = %q, want trigger", received.EventAction)
	}
	if received.Payload.Source != "portwatch" {
		t.Errorf("source = %q, want portwatch", received.Payload.Source)
	}
	if received.Payload.Severity != "warning" {
		t.Errorf("severity = %q, want warning", received.Payload.Severity)
	}
}

func TestPagerDutyHandler_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := NewPagerDutyHandler("key")
	h.apiURL = ts.URL

	changes := []monitor.Change{buildPDChange(443, monitor.Opened)}
	if err := h.Handle(changes); err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
}

func TestPagerDutyHandler_EmptySlice(t *testing.T) {
	// Ensure that an explicitly empty (non-nil) slice is also treated as a no-op
	// and does not attempt a network request.
	h := NewPagerDutyHandler("test-key")
	if err := h.Handle([]monitor.Change{}); err != nil {
		t.Fatalf("expected no error for empty changes slice, got %v", err)
	}
}
