package alert_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

func buildGrafanaChange(port int, kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Port: scanner.PortInfo{Port: port, Proto: "tcp"},
	}
}

func TestGrafanaHandler_NoChanges(t *testing.T) {
	h := alert.NewGrafanaHandler("http://localhost", "key", "dash1", 5)
	if err := h.Handle(nil); err != nil {
		t.Fatalf("expected no error on empty changes, got %v", err)
	}
}

func TestGrafanaHandler_SendsPayload(t *testing.T) {
	var captured map[string]interface{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/annotations" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("missing or wrong Authorization header")
		}
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &captured)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	h := alert.NewGrafanaHandler(srv.URL, "test-key", "portwatch", 5)
	changes := []monitor.Change{
		buildGrafanaChange(8080, monitor.Opened),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if captured["dashboardId"] != "portwatch" {
		t.Errorf("unexpected dashboardId: %v", captured["dashboardId"])
	}
	if captured["text"] == "" {
		t.Error("expected non-empty annotation text")
	}
}

func TestGrafanaHandler_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	h := alert.NewGrafanaHandler(srv.URL, "key", "dash1", 5)
	changes := []monitor.Change{buildGrafanaChange(443, monitor.Closed)}
	if err := h.Handle(changes); err == nil {
		t.Fatal("expected error on server 500")
	}
}
