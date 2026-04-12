package alert_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

func buildInfluxDBChange(port int, kind monitor.ChangeType) monitor.Change {
	return monitor.Change{
		Type: kind,
		Port: scanner.PortInfo{Port: port, Protocol: "tcp"},
	}
}

func TestInfluxDBHandler_NoChanges(t *testing.T) {
	h := alert.NewInfluxDBHandler("http://localhost:8086", "tok", "org", "bkt", "m", 5)
	if err := h.Handle(nil); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestInfluxDBHandler_SendsPayload(t *testing.T) {
	var body string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		body = string(b)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	h := alert.NewInfluxDBHandler(ts.URL, "mytoken", "myorg", "mybucket", "portwatch_changes", 5)
	changes := []monitor.Change{
		buildInfluxDBChange(8080, monitor.Opened),
		buildInfluxDBChange(9090, monitor.Closed),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(body, "port=8080") {
		t.Errorf("expected port=8080 in payload, got:\n%s", body)
	}
	if !strings.Contains(body, "port=9090") {
		t.Errorf("expected port=9090 in payload, got:\n%s", body)
	}
	if !strings.Contains(body, `kind=\"opened\"`) && !strings.Contains(body, "opened") {
		t.Errorf("expected 'opened' in payload, got:\n%s", body)
	}
}

func TestInfluxDBHandler_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := alert.NewInfluxDBHandler(ts.URL, "tok", "org", "bkt", "m", 5)
	err := h.Handle([]monitor.Change{buildInfluxDBChange(443, monitor.Opened)})
	if err == nil {
		t.Fatal("expected error on 500 response")
	}
}

func TestInfluxDBHandler_InvalidURL(t *testing.T) {
	h := alert.NewInfluxDBHandler("://bad-url", "tok", "org", "bkt", "m", 5)
	err := h.Handle([]monitor.Change{buildInfluxDBChange(80, monitor.Opened)})
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
