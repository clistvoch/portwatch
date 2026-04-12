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

func buildTelegramChange(port int, state monitor.ChangeType) monitor.Change {
	return monitor.Change{
		Type: state,
		Port: scanner.PortInfo{Port: port, Proto: "tcp"},
	}
}

func TestTelegramHandler_NoChanges(t *testing.T) {
	h := NewTelegramHandler("tok", "123", "HTML")
	if err := h.Handle(nil); err != nil {
		t.Fatalf("expected no error for empty changes, got %v", err)
	}
}

func TestTelegramHandler_SendsPayload(t *testing.T) {
	var received map[string]interface{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	h := NewTelegramHandler("mytoken", "-10099", "HTML")
	h.apiBase = srv.URL

	changes := []monitor.Change{
		buildTelegramChange(8080, monitor.Opened),
		buildTelegramChange(9090, monitor.Closed),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["chat_id"] != "-10099" {
		t.Errorf("unexpected chat_id: %v", received["chat_id"])
	}
	text, _ := received["text"].(string)
	if !strings.Contains(text, "2 port change") {
		t.Errorf("expected change count in message, got: %s", text)
	}
	if received["parse_mode"] != "HTML" {
		t.Errorf("unexpected parse_mode: %v", received["parse_mode"])
	}
}

func TestTelegramHandler_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	h := NewTelegramHandler("badtoken", "123", "HTML")
	h.apiBase = srv.URL

	err := h.Handle([]monitor.Change{buildTelegramChange(22, monitor.Opened)})
	if err == nil {
		t.Fatal("expected error on non-2xx response")
	}
}
