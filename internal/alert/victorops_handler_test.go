package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/monitor"
)

func buildVictorOpsChange(port int, kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Port: monitor.PortInfo{Port: port, Proto: "tcp"},
		Kind: kind,
	}
}

func TestVictorOpsHandler_NoChanges(t *testing.T) {
	h := NewVictorOpsHandler("https://example.victorops.com", "team", "CRITICAL")
	if err := h.Handle(nil); err != nil {
		t.Fatalf("expected no error for empty changes, got %v", err)
	}
}

func TestVictorOpsHandler_SendsPayload(t *testing.T) {
	var received map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	h := NewVictorOpsHandler(server.URL, "ops-team", "WARNING")
	changes := []monitor.Change{
		buildVictorOpsChange(8080, monitor.Opened),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["message_type"] != "WARNING" {
		t.Errorf("expected message_type WARNING, got %v", received["message_type"])
	}
	if received["entity_id"] != "portwatch.port-change" {
		t.Errorf("unexpected entity_id: %v", received["entity_id"])
	}
}

func TestVictorOpsHandler_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	h := NewVictorOpsHandler(server.URL, "ops-team", "CRITICAL")
	changes := []monitor.Change{buildVictorOpsChange(443, monitor.Closed)}
	if err := h.Handle(changes); err == nil {
		t.Fatal("expected error on server 500 response")
	}
}

func TestVictorOpsHandler_InvalidURL(t *testing.T) {
	h := NewVictorOpsHandler("http://127.0.0.1:0", "key", "INFO")
	changes := []monitor.Change{buildVictorOpsChange(22, monitor.Opened)}
	if err := h.Handle(changes); err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
