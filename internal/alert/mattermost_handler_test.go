package alert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/monitor"
)

func buildMattermostChange(port int, kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Port: monitor.PortInfo{Port: port, Proto: "tcp"},
		Kind: kind,
	}
}

func TestMattermostHandler_NoChanges(t *testing.T) {
	h := NewMattermostHandler("http://invalid", "#general", "portwatch", ":shield:")
	if err := h.Handle(nil); err != nil {
		t.Fatalf("expected no error for empty changes, got %v", err)
	}
}

func TestMattermostHandler_SendsPayload(t *testing.T) {
	var received []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := NewMattermostHandler(ts.URL, "#alerts", "bot", ":bell:")
	changes := []monitor.Change{
		buildMattermostChange(8080, monitor.Opened),
		buildMattermostChange(9090, monitor.Closed),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(received, &payload); err != nil {
		t.Fatalf("invalid JSON payload: %v", err)
	}
	text, _ := payload["text"].(string)
	if !strings.Contains(text, "2 port change") {
		t.Errorf("expected change count in text, got: %q", text)
	}
	if payload["channel"] != "#alerts" {
		t.Errorf("unexpected channel: %v", payload["channel"])
	}
	if payload["username"] != "bot" {
		t.Errorf("unexpected username: %v", payload["username"])
	}
}

func TestMattermostHandler_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := NewMattermostHandler(ts.URL, "", "portwatch", ":shield:")
	err := h.Handle([]monitor.Change{buildMattermostChange(443, monitor.Opened)})
	if err == nil {
		t.Fatal("expected error on 500 response")
	}
}
