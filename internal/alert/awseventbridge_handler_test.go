package alert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/monitor"
)

func buildEBChange(kind monitor.ChangeKind, port int) monitor.Change {
	return monitor.Change{Kind: kind, Port: monitor.PortInfo{Port: port, Proto: "tcp"}}
}

func TestAWSEventBridgeHandler_NoChanges(t *testing.T) {
	h := NewAWSEventBridgeHandler("http://localhost", "default", "portwatch", "PortChange")
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAWSEventBridgeHandler_SendsPayload(t *testing.T) {
	var received []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	h := NewAWSEventBridgeHandler(srv.URL, "default", "portwatch", "PortChange")
	changes := []monitor.Change{
		buildEBChange(monitor.ChangeOpened, 9200),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(received, &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	entries, ok := payload["Entries"].([]interface{})
	if !ok || len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %v", payload["Entries"])
	}
}

func TestAWSEventBridgeHandler_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	h := NewAWSEventBridgeHandler(srv.URL, "default", "portwatch", "PortChange")
	changes := []monitor.Change{buildEBChange(monitor.ChangeOpened, 443)}
	if err := h.Handle(changes); err == nil {
		t.Error("expected error on 500 response")
	}
}

func TestAWSEventBridgeHandler_InvalidURL(t *testing.T) {
	h := NewAWSEventBridgeHandler("http://127.0.0.1:1", "default", "portwatch", "PortChange")
	changes := []monitor.Change{buildEBChange(monitor.ChangeClosed, 80)}
	if err := h.Handle(changes); err == nil {
		t.Error("expected connection error")
	}
}
