package scanner

import (
	"net"
	"strconv"
	"testing"
)

// startTCPListener opens a TCP listener on an ephemeral port and returns the port number and a closer.
func startTCPListener(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	port, _ := strconv.Atoi(ln.Addr().(*net.TCPAddr).Port.String())
	_, portStr, _ := net.SplitHostPort(ln.Addr().String())
	port64, _ := strconv.ParseInt(portStr, 10, 64)
	return int(port64), func() { ln.Close() }
}

func TestNewScanner_Defaults(t *testing.T) {
	s := NewScanner("127.0.0.1", 1, 1024)
	if s.Host != "127.0.0.1" {
		t.Errorf("expected host 127.0.0.1, got %s", s.Host)
	}
	if s.PortRange[0] != 1 || s.PortRange[1] != 1024 {
		t.Errorf("unexpected port range: %v", s.PortRange)
	}
	if len(s.Protocols) == 0 {
		t.Error("expected at least one protocol")
	}
}

func TestScan_InvalidRange(t *testing.T) {
	s := NewScanner("127.0.0.1", 1024, 1)
	_, err := s.Scan()
	if err == nil {
		t.Error("expected error for invalid port range")
	}
}

func TestScan_DetectsOpenPort(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	defer ln.Close()

	_, portStr, _ := net.SplitHostPort(ln.Addr().String())
	port64, _ := strconv.ParseInt(portStr, 10, 64)
	port := int(port64)

	s := &Scanner{
		Host:      "127.0.0.1",
		PortRange: [2]int{port, port},
		Protocols: []string{"tcp"},
	}

	results, err := s.Scan()
	if err != nil {
		t.Fatalf("scan error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 open port, got %d", len(results))
	}
	if results[0].Port != port {
		t.Errorf("expected port %d, got %d", port, results[0].Port)
	}
}

func TestPortInfo_String(t *testing.T) {
	p := PortInfo{Port: 8080, Protocol: "tcp", Address: "127.0.0.1"}
	got := p.String()
	want := "127.0.0.1:8080 (tcp)"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
