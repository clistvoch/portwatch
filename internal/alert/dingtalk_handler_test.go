package alert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/patrickward/portwatch/internal/monitor"
)

func buildDingTalkChange(port int, kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Port: monitor.PortInfo{Port: port, Proto: "tcp"},
		Kind: kind,
	}
}

func TestDingTalkHandler_NoChanges(t *testing.T) {
	h := NewDingTalkHandler("http://example.com", "", "text")
	if err := h.Handle(nil); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestDingTalkHandler_SendsPayload(t *testing.T) {
	var received []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := NewDingTalkHandler(ts.URL, "", "text")
	changes := []monitor.Change{
		buildDingTalkChange(8080, monitor.ChangeOpened),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(received, &payload); err != nil {
		t.Fatalf("invalid JSON payload: %v", err)
	}
	if payload["msgtype"] != "text" {
		t.Errorf("expected msgtype=text, got %v", payload["msgtype"])
	}
}

func TestDingTalkHandler_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := NewDingTalkHandler(ts.URL, "", "text")
	changes := []monitor.Change{buildDingTalkChange(443, monitor.ChangeOpened)}
	if err := h.Handle(changes); err == nil {
		t.Fatal("expected error on server error, got nil")
	}
}

func TestDingTalkHandler_MarkdownType(t *testing.T) {
	var received []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := NewDingTalkHandler(ts.URL, "", "markdown")
	changes := []monitor.Change{buildDingTalkChange(22, monitor.ChangeClosed)}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(received, &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if payload["msgtype"] != "markdown" {
		t.Errorf("expected msgtype=markdown, got %v", payload["msgtype"])
	}
	if payload["markdown"] == nil {
		t.Error("expected markdown field in payload")
	}
}
