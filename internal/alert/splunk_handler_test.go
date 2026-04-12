package alert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/monitor"
)

func buildSplunkChange(port int, kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{Port: port, Kind: kind}
}

func TestSplunkHandler_NoChanges(t *testing.T) {
	h := NewSplunkHandler("http://localhost", "tok", "main", "portwatch", 5)
	if err := h.Handle(nil); err != nil {
		t.Errorf("expected nil error for empty changes, got %v", err)
	}
}

func TestSplunkHandler_SendsPayload(t *testing.T) {
	var received []map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Splunk test-token" {
			t.Errorf("unexpected auth header: %s", r.Header.Get("Authorization"))
		}
		body, _ := io.ReadAll(r.Body)
		var ev map[string]interface{}
		if err := json.Unmarshal(body, &ev); err != nil {
			t.Errorf("invalid JSON payload: %v", err)
		}
		received = append(received, ev)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := NewSplunkHandler(ts.URL, "test-token", "portwatch", "portwatch", 5)
	changes := []monitor.Change{
		buildSplunkChange(8080, monitor.Opened),
		buildSplunkChange(9090, monitor.Closed),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(received) != 2 {
		t.Errorf("expected 2 requests, got %d", len(received))
	}
	event, ok := received[0]["event"].(map[string]interface{})
	if !ok {
		t.Fatal("expected event field in payload")
	}
	if int(event["port"].(float64)) != 8080 {
		t.Errorf("expected port 8080, got %v", event["port"])
	}
}

func TestSplunkHandler_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := NewSplunkHandler(ts.URL, "tok", "main", "portwatch", 5)
	err := h.Handle([]monitor.Change{buildSplunkChange(443, monitor.Opened)})
	if err == nil {
		t.Error("expected error on server 500 response")
	}
}

func TestSplunkHandler_InvalidURL(t *testing.T) {
	h := NewSplunkHandler("://bad-url", "tok", "main", "portwatch", 5)
	err := h.Handle([]monitor.Change{buildSplunkChange(80, monitor.Opened)})
	if err == nil {
		t.Error("expected error for invalid URL")
	}
}
