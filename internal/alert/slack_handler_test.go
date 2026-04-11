package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

func buildSlackChange(port uint16, state monitor.ChangeType) monitor.Change {
	return monitor.Change{
		Type: state,
		Port: scanner.PortInfo{Port: port, Proto: "tcp"},
	}
}

func TestSlackHandler_NoChanges(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	h := NewSlackHandler(ts.URL, "portwatch", ":bell:")
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty changes")
	}
}

func TestSlackHandler_SendsPayload(t *testing.T) {
	var received map[string]interface{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := NewSlackHandler(ts.URL, "portwatch", ":bell:")
	changes := []monitor.Change{
		buildSlackChange(8080, monitor.Opened),
	}

	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["username"] != "portwatch" {
		t.Errorf("expected username 'portwatch', got %v", received["username"])
	}
	text, _ := received["text"].(string)
	if !strings.Contains(text, "8080") {
		t.Errorf("expected text to contain port 8080, got: %s", text)
	}
}

func TestSlackHandler_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := NewSlackHandler(ts.URL, "portwatch", ":bell:")
	changes := []monitor.Change{buildSlackChange(443, monitor.Closed)}

	err := h.Handle(changes)
	if err == nil {
		t.Fatal("expected error on server 500, got nil")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("expected error to mention status 500, got: %v", err)
	}
}
