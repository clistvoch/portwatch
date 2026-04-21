package alert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/robfig/portwatch/internal/monitor"
)

func buildTransformChange(port int, kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{Port: port, Kind: kind}
}

func TestWebhookTransformHandler_NoChanges(t *testing.T) {
	h, err := NewWebhookTransformHandler("http://localhost", "application/json", "", true)
	if err != nil {
		t.Fatal(err)
	}
	if err := h.Handle(nil); err != nil {
		t.Errorf("expected no error on empty changes, got: %v", err)
	}
}

func TestWebhookTransformHandler_DefaultPayload(t *testing.T) {
	var received []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	h, err := NewWebhookTransformHandler(srv.URL, "application/json", "", true)
	if err != nil {
		t.Fatal(err)
	}
	changes := []monitor.Change{buildTransformChange(8080, monitor.Opened)}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(received, &payload); err != nil {
		t.Fatalf("invalid JSON payload: %v", err)
	}
	if _, ok := payload["changes"]; !ok {
		t.Error("expected 'changes' key in payload")
	}
}

func TestWebhookTransformHandler_TemplatePayload(t *testing.T) {
	var received string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		received = string(b)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	h, err := NewWebhookTransformHandler(srv.URL, "text/plain", "port={{.Port}}", false)
	if err != nil {
		t.Fatal(err)
	}
	changes := []monitor.Change{buildTransformChange(9090, monitor.Opened)}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(received, "port=9090") {
		t.Errorf("expected template output in body, got: %s", received)
	}
}

func TestWebhookTransformHandler_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	h, err := NewWebhookTransformHandler(srv.URL, "application/json", "", true)
	if err != nil {
		t.Fatal(err)
	}
	changes := []monitor.Change{buildTransformChange(443, monitor.Closed)}
	if err := h.Handle(changes); err == nil {
		t.Error("expected error on 500 response")
	}
}

func TestWebhookTransformHandler_InvalidURL(t *testing.T) {
	_, err := NewWebhookTransformHandler("", "application/json", "", true)
	if err == nil {
		t.Error("expected error for empty URL")
	}
}
