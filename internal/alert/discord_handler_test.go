package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

func buildDiscordChange(port int, state monitor.ChangeType) monitor.Change {
	return monitor.Change{
		Type: state,
		Port: scanner.PortInfo{Port: port, Proto: "tcp"},
	}
}

func TestDiscordHandler_NoChanges(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	h := NewDiscordHandler(ts.URL, "portwatch", "")
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty changes")
	}
}

func TestDiscordHandler_SendsPayload(t *testing.T) {
	var received discordPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	h := NewDiscordHandler(ts.URL, "bot", "http://example.com/avatar.png")
	changes := []monitor.Change{
		buildDiscordChange(8080, monitor.Opened),
		buildDiscordChange(9090, monitor.Closed),
	}

	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.Username != "bot" {
		t.Errorf("expected username 'bot', got %q", received.Username)
	}
	if len(received.Embeds) != 1 {
		t.Fatalf("expected 1 embed, got %d", len(received.Embeds))
	}
	if received.Embeds[0].Color != 0xFF4444 {
		t.Errorf("expected color 0xFF4444, got %d", received.Embeds[0].Color)
	}
}

func TestDiscordHandler_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := NewDiscordHandler(ts.URL, "portwatch", "")
	changes := []monitor.Change{buildDiscordChange(443, monitor.Opened)}

	if err := h.Handle(changes); err == nil {
		t.Error("expected error on server 500, got nil")
	}
}
