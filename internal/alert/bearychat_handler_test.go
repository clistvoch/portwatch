package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/monitor"
)

func buildBearyChatChange(t monitor.ChangeType, port int) monitor.Change {
	return monitor.Change{Type: t, Port: monitor.PortInfo{Port: port, Proto: "tcp"}}
}

func TestBearyChatHandler_NoChanges(t *testing.T) {
	h := NewBearyChatHandler("http://example.com", "#test", 5)
	if err := h.Handle(nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestBearyChatHandler_SendsPayload(t *testing.T) {
	var received bearyChatPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := NewBearyChatHandler(ts.URL, "#alerts", 5)
	changes := []monitor.Change{
		buildBearyChatChange(monitor.Opened, 8080),
		buildBearyChatChange(monitor.Closed, 9090),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(received.Attachments) != 2 {
		t.Errorf("expected 2 attachments, got %d", len(received.Attachments))
	}
	if received.Channel != "#alerts" {
		t.Errorf("expected channel #alerts, got %s", received.Channel)
	}
}

func TestBearyChatHandler_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := NewBearyChatHandler(ts.URL, "#alerts", 5)
	changes := []monitor.Change{buildBearyChatChange(monitor.Opened, 443)}
	if err := h.Handle(changes); err == nil {
		t.Fatal("expected error on 500 response")
	}
}
