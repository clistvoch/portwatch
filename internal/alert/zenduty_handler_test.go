package alert_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/monitor"
)

func buildZendutyChange(port int, kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{Port: port, Kind: kind}
}

func TestZendutyHandler_NoChanges(t *testing.T) {
	h := alert.NewZendutyHandler("key", "svc", "critical", "portwatch alert", "http://example.com")
	if err := h.Handle([]monitor.Change{}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestZendutyHandler_SendsPayload(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	h := alert.NewZendutyHandler("key", "svc", "critical", "portwatch alert", ts.URL)
	changes := []monitor.Change{
		buildZendutyChange(8080, monitor.Opened),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["summary"] == nil {
		t.Error("expected summary field in payload")
	}
}

func TestZendutyHandler_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := alert.NewZendutyHandler("key", "svc", "critical", "portwatch alert", ts.URL)
	changes := []monitor.Change{buildZendutyChange(443, monitor.Closed)}
	if err := h.Handle(changes); err == nil {
		t.Fatal("expected error on 500 response")
	}
}
