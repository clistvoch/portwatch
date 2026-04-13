package alert

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"
)

func freePort(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("freePort: %v", err)
	}
	addr := l.Addr().String()
	l.Close()
	return addr
}

func TestHealthCheckHandler_ReturnsOK(t *testing.T) {
	addr := freePort(t)
	h, err := NewHealthCheckHandler(addr, "/healthz")
	if err != nil {
		t.Fatalf("NewHealthCheckHandler: %v", err)
	}
	defer h.Close()
	time.Sleep(50 * time.Millisecond)

	resp, err := http.Get(fmt.Sprintf("http://%s/healthz", addr))
	if err != nil {
		t.Fatalf("GET /healthz: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "ok\n" {
		t.Errorf("unexpected body: %q", body)
	}
}

func TestHealthCheckHandler_ReturnsUnhealthy(t *testing.T) {
	addr := freePort(t)
	h, err := NewHealthCheckHandler(addr, "/healthz")
	if err != nil {
		t.Fatalf("NewHealthCheckHandler: %v", err)
	}
	defer h.Close()
	h.SetHealthy(false)
	time.Sleep(50 * time.Millisecond)

	resp, err := http.Get(fmt.Sprintf("http://%s/healthz", addr))
	if err != nil {
		t.Fatalf("GET /healthz: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", resp.StatusCode)
	}
}

func TestHealthCheckHandler_Close(t *testing.T) {
	addr := freePort(t)
	h, err := NewHealthCheckHandler(addr, "/healthz")
	if err != nil {
		t.Fatalf("NewHealthCheckHandler: %v", err)
	}
	time.Sleep(30 * time.Millisecond)
	if err := h.Close(); err != nil {
		t.Errorf("Close returned error: %v", err)
	}
}
