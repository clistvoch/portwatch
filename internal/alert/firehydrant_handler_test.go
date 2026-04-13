package alert_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/portwatch/internal/alert"
	"github.com/example/portwatch/internal/monitor"
)

func buildFHChange(port int, kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Port: monitor.PortInfo{Port: port, Proto: "tcp"},
		Kind: kind,
	}
}

func TestFirehydrantHandler_NoChanges(t *testing.T) {
	h := alert.NewFirehydrantHandler("key", "svc-1", "http://localhost", 5)
	if err := h.Handle(nil); err != nil {
		t.Fatalf("expected no error on empty changes, got: %v", err)
	}
}

func TestFirehydrantHandler_SendsPayload(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("missing or wrong Authorization header")
		}
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	h := alert.NewFirehydrantHandler("test-key", "svc-42", ts.URL, 5)
	changes := []monitor.Change{
		buildFHChange(8080, monitor.Opened),
		buildFHChange(9090, monitor.Closed),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["service_id"] != "svc-42" {
		t.Errorf("unexpected service_id: %v", received["service_id"])
	}
	if received["summary"] == "" {
		t.Error("expected non-empty summary")
	}
}

func TestFirehydrantHandler_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := alert.NewFirehydrantHandler("key", "svc-1", ts.URL, 5)
	err := h.Handle([]monitor.Change{buildFHChange(443, monitor.Opened)})
	if err == nil {
		t.Fatal("expected error on 500 response")
	}
}

func TestFirehydrantHandler_InvalidURL(t *testing.T) {
	h := alert.NewFirehydrantHandler("key", "svc-1", "http://127.0.0.1:1", 2)
	err := h.Handle([]monitor.Change{buildFHChange(80, monitor.Closed)})
	if err == nil {
		t.Fatal("expected connection error")
	}
}
