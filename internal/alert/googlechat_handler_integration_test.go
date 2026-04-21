package alert_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/winhowes/portwatch/internal/alert"
	"github.com/winhowes/portwatch/internal/monitor"
)

func TestGoogleChatHandler_RoundTrip(t *testing.T) {
	var received map[string]interface{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	h := alert.NewGoogleChatHandler(srv.URL)

	changes := []monitor.Change{
		{Type: monitor.Opened, Port: 9090, Proto: "tcp"},
		{Type: monitor.Closed, Port: 8080, Proto: "tcp"},
	}

	if err := h.Handle(changes); err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}

	if received == nil {
		t.Fatal("server received no payload")
	}
	text, ok := received["text"].(string)
	if !ok || text == "" {
		t.Errorf("expected non-empty 'text' field in payload, got: %v", received)
	}
}

func TestGoogleChatHandler_RoundTrip_NoChanges(t *testing.T) {
	called := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	h := alert.NewGoogleChatHandler(srv.URL)
	if err := h.Handle([]monitor.Change{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call when there are no changes")
	}
}
