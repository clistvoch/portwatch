package alert_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/monitor"
)

func TestWebhookHandler_SendsPayload(t *testing.T) {
	var received alert.WebhookPayload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := alert.NewWebhookHandler(ts.URL)
	changes := []monitor.Change{
		{Port: 8080, Kind: monitor.Opened},
		{Port: 9090, Kind: monitor.Closed},
	}

	if err := h.Handle(changes); err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}

	if len(received.Changes) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(received.Changes))
	}
	if received.Changes[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", received.Changes[0].Port)
	}
	if received.Timestamp == "" {
		t.Error("expected non-empty timestamp")
	}
}

func TestWebhookHandler_NoChanges(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	h := alert.NewWebhookHandler(ts.URL)
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty changes")
	}
}

func TestWebhookHandler_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := alert.NewWebhookHandler(ts.URL)
	changes := []monitor.Change{{Port: 80, Kind: monitor.Opened}}

	if err := h.Handle(changes); err == nil {
		t.Error("expected error on 500 response, got nil")
	}
}
