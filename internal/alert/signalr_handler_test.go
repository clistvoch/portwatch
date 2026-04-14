package alert_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/monitor"
)

func buildSignalRChange(port int, kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Port: monitor.PortInfo{Port: port, Proto: "tcp"},
		Kind: kind,
	}
}

func TestSignalRHandler_NoChanges(t *testing.T) {
	h := alert.NewSignalRHandler("https://example.service.signalr.net", "portwatch", "portChange", "", 5)
	if err := h.Handle(nil); err != nil {
		t.Fatalf("expected nil error for empty changes, got %v", err)
	}
}

func TestSignalRHandler_SendsPayload(t *testing.T) {
	var received map[string]interface{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	h := alert.NewSignalRHandler(srv.URL, "portwatch", "portChange", "secret", 5)
	changes := []monitor.Change{
		buildSignalRChange(8080, monitor.Opened),
	}

	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["target"] != "portChange" {
		t.Errorf("expected target=portChange, got %v", received["target"])
	}
	args, ok := received["arguments"].([]interface{})
	if !ok || len(args) == 0 {
		t.Errorf("expected non-empty arguments array")
	}
}

func TestSignalRHandler_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	h := alert.NewSignalRHandler(srv.URL, "portwatch", "portChange", "", 5)
	changes := []monitor.Change{buildSignalRChange(443, monitor.Closed)}

	if err := h.Handle(changes); err == nil {
		t.Fatal("expected error on server 500, got nil")
	}
}
