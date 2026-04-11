package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

func buildOGChange(port uint16, kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Port: scanner.PortInfo{Port: port, Proto: "tcp"},
	}
}

func TestOpsGenieHandler_NoChanges(t *testing.T) {
	h := NewOpsGenieHandler("key", "")
	if err := h.Handle(nil); err != nil {
		t.Fatalf("expected no error for empty changes, got %v", err)
	}
}

func TestOpsGenieHandler_SendsPayload(t *testing.T) {
	var received opsGeniePayload
	var authHeader string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	h := NewOpsGenieHandler("test-key", srv.URL)
	changes := []monitor.Change{
		buildOGChange(8080, monitor.Opened),
		buildOGChange(9090, monitor.Closed),
	}

	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if authHeader != "GenieKey test-key" {
		t.Errorf("expected GenieKey auth header, got %q", authHeader)
	}
	if received.Message == "" {
		t.Error("expected non-empty message")
	}
	if received.Priority != "P3" {
		t.Errorf("expected priority P3, got %q", received.Priority)
	}
}

func TestOpsGenieHandler_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	h := NewOpsGenieHandler("key", srv.URL)
	changes := []monitor.Change{buildOGChange(443, monitor.Opened)}

	if err := h.Handle(changes); err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestOpsGenieHandler_InvalidURL(t *testing.T) {
	h := NewOpsGenieHandler("key", "http://127.0.0.1:0")
	changes := []monitor.Change{buildOGChange(80, monitor.Opened)}
	if err := h.Handle(changes); err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
