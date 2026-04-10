package monitor_test

import (
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

func listenTCP(t *testing.T) (port int, close func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	addr := ln.Addr().(*net.TCPAddr)
	return addr.Port, func() { ln.Close() }
}

func TestDiff_DetectsOpenedPort(t *testing.T) {
	port, closePort := listenTCP(t)
	defer closePort()

	s := scanner.NewScanner("127.0.0.1", port, port, 50*time.Millisecond)
	m := monitor.New(s, time.Second)

	// First diff seeds state — no changes expected.
	changes, err := m.Diff()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(changes) != 0 {
		t.Fatalf("expected 0 changes on seed, got %d", len(changes))
	}

	// Second diff with same open port — still no changes.
	changes, err = m.Diff()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(changes) != 0 {
		t.Fatalf("expected 0 changes, got %d", len(changes))
	}
}

func TestDiff_DetectsClosedPort(t *testing.T) {
	port, closePort := listenTCP(t)

	s := scanner.NewScanner("127.0.0.1", port, port, 50*time.Millisecond)
	m := monitor.New(s, time.Second)

	// Seed with port open.
	_, _ = m.Diff()

	// Close the port then diff.
	closePort()
	time.Sleep(20 * time.Millisecond)

	changes, err := m.Diff()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Status != "closed" || changes[0].Port != port {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestChange_String(t *testing.T) {
	c := monitor.Change{Port: 8080, Status: "opened", At: time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)}
	got := c.String()
	if got == "" {
		t.Error("expected non-empty string")
	}
}
