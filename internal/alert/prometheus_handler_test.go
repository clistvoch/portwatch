package alert

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"portwatch/internal/monitor"
	"portwatch/internal/scanner"
)

func buildPrometheusChange(t monitor.ChangeType, port uint16) monitor.Change {
	return monitor.Change{
		Type: t,
		Port: scanner.PortInfo{Port: port, Proto: "tcp"},
	}
}

func TestPrometheusHandler_NoChanges(t *testing.T) {
	h, err := NewPrometheusHandler("127.0.0.1:0", "/metrics")
	if err != nil {
		t.Fatalf("new handler: %v", err)
	}
	defer h.Close()

	if err := h.Handle(nil); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if h.total.Load() != 0 {
		t.Errorf("expected 0 changes, got %d", h.total.Load())
	}
}

func TestPrometheusHandler_CountsChanges(t *testing.T) {
	h, err := NewPrometheusHandler("127.0.0.1:19091", "/metrics")
	if err != nil {
		t.Fatalf("new handler: %v", err)
	}
	defer h.Close()

	changes := []monitor.Change{
		buildPrometheusChange(monitor.Opened, 8080),
		buildPrometheusChange(monitor.Opened, 443),
		buildPrometheusChange(monitor.Closed, 22),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("handle: %v", err)
	}
	if h.total.Load() != 3 {
		t.Errorf("expected 3, got %d", h.total.Load())
	}
}

func TestPrometheusHandler_MetricsEndpoint(t *testing.T) {
	h, err := NewPrometheusHandler("127.0.0.1:19092", "/metrics")
	if err != nil {
		t.Fatalf("new handler: %v", err)
	}
	defer h.Close()

	_ = h.Handle([]monitor.Change{buildPrometheusChange(monitor.Opened, 9000)})

	time.Sleep(50 * time.Millisecond)

	resp, err := http.Get("http://127.0.0.1:19092/metrics")
	if err != nil {
		t.Fatalf("GET /metrics: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "portwatch_ports_opened_total") {
		t.Error("expected portwatch_ports_opened_total in metrics output")
	}
}
