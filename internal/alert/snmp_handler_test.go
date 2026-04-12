package alert

import (
	"log"
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

func buildSNMPChange(port int, state monitor.State) monitor.Change {
	return monitor.Change{
		Port:  port,
		State: state,
	}
}

func startUDPListener(t *testing.T) (*net.UDPConn, int) {
	t.Helper()
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })
	return conn, conn.LocalAddr().(*net.UDPAddr).Port
}

func TestSNMPHandler_NoChanges(t *testing.T) {
	_, port := startUDPListener(t)
	h := NewSNMPHandler("127.0.0.1", port, "public", log.Default())
	if err := h.Handle(nil); err != nil {
		t.Errorf("expected no error for empty changes, got %v", err)
	}
}

func TestSNMPHandler_SendsTrap(t *testing.T) {
	conn, port := startUDPListener(t)

	h := NewSNMPHandler("127.0.0.1", port, "public", log.Default())
	changes := []monitor.Change{
		buildSNMPChange(8080, monitor.Opened),
		buildSNMPChange(9090, monitor.Closed),
	}

	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	buf := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		t.Fatalf("did not receive UDP packet: %v", err)
	}
	if n == 0 {
		t.Error("received empty UDP packet")
	}
	// First byte should be SEQUENCE tag 0x30
	if buf[0] != 0x30 {
		t.Errorf("expected SEQUENCE tag 0x30, got 0x%02x", buf[0])
	}
}

func TestSNMPHandler_InvalidTarget(t *testing.T) {
	h := NewSNMPHandler("192.0.2.1", 162, "public", log.Default())
	changes := []monitor.Change{buildSNMPChange(80, monitor.Opened)}
	// Should fail to dial or write to the unreachable host within timeout.
	// We accept either outcome; just ensure no panic.
	_ = h.Handle(changes)
}
