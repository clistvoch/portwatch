package alert_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/natemollica-nm/portwatch/internal/alert"
	"github.com/natemollica-nm/portwatch/internal/monitor"
)

func buildGoogleChatChange(port int, kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Port: monitor.PortInfo{Port: port, Proto: "tcp"},
		Kind: kind,
	}
}

func TestGoogleChatHandler_NoChanges(t *testing.T) {
	h := alert.NewGoogleChatHandler("https://example.com", "")
	if err := h.Handle(nil); err != nil {
		t.Fatalf("expected no error on empty changes, got %v", err)
	}
}

func TestGoogleChatHandler_SendsPayload(t *testing.T) {
	var received []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := alert.NewGoogleChatHandler(ts.URL, "portwatch")
	changes := []monitor.Change{
		buildGoogleChatChange(8080, monitor.Opened),
		buildGoogleChatChange(9090, monitor.Closed),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(received, &payload); err != nil {
		t.Fatalf("invalid JSON payload: %v", err)
	}
	text, _ := payload["text"].(string)
	if !strings.Contains(text, "2 port change") {
		t.Errorf("expected change count in text, got: %q", text)
	}
}

func TestGoogleChatHandler_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := alert.NewGoogleChatHandler(ts.URL, "")
	err := h.Handle([]monitor.Change{buildGoogleChatChange(443, monitor.Opened)})
	if err == nil {
		t.Fatal("expected error on 500 response")
	}
}
