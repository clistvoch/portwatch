package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/monitor"
)

func buildHipChatChange(port int, kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Port: monitor.PortInfo{Port: port, Proto: "tcp"},
		Kind: kind,
	}
}

func TestHipChatHandler_NoChanges(t *testing.T) {
	h := NewHipChatHandler("42", "tok", "yellow", "http://localhost", false)
	if err := h.Handle(nil); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestHipChatHandler_SendsPayload(t *testing.T) {
	var received map[string]interface{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			http.Error(w, "bad body", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	h := NewHipChatHandler("99", "mytoken", "red", srv.URL, true)
	changes := []monitor.Change{
		buildHipChatChange(8080, monitor.Opened),
	}

	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["color"] != "red" {
		t.Errorf("expected color red, got %v", received["color"])
	}
	if received["notify"] != true {
		t.Errorf("expected notify true, got %v", received["notify"])
	}
	if received["message"] == "" {
		t.Error("expected non-empty message")
	}
}

func TestHipChatHandler_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	h := NewHipChatHandler("1", "tok", "green", srv.URL, false)
	err := h.Handle([]monitor.Change{buildHipChatChange(443, monitor.Closed)})
	if err == nil {
		t.Fatal("expected error on 500 response")
	}
}
