package alert

import (
	"net"
	"strings"
	"testing"

	"github.com/netwatch/portwatch/internal/monitor"
)

func buildXMPPChange(port int, kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Port: monitor.PortInfo{Port: port, Proto: "tcp"},
		Kind: kind,
	}
}

func TestXMPPHandler_NoChanges(t *testing.T) {
	h := NewXMPPHandler("localhost", 5222, "u", "p", "to@example.com", false)
	if err := h.Handle(nil); err != nil {
		t.Fatalf("expected no error on empty changes, got %v", err)
	}
}

func TestXMPPHandler_SendsPayload(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()

	received := make(chan string, 1)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		buf := make([]byte, 4096)
		n, _ := conn.Read(buf)
		received <- string(buf[:n])
	}()

	addr := ln.Addr().(*net.TCPAddr)
	h := NewXMPPHandler("127.0.0.1", addr.Port, "u", "p", "to@example.com", false)
	changes := []monitor.Change{
		buildXMPPChange(8080, monitor.ChangeOpened),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("Handle: %v", err)
	}

	payload := <-received
	if !strings.Contains(payload, "portwatch") {
		t.Errorf("expected 'portwatch' in payload, got: %s", payload)
	}
	if !strings.Contains(payload, "8080") {
		t.Errorf("expected port 8080 in payload, got: %s", payload)
	}
	if !strings.Contains(payload, "<message") {
		t.Errorf("expected XMPP stanza in payload, got: %s", payload)
	}
}

func TestXMPPHandler_ServerError(t *testing.T) {
	h := NewXMPPHandler("127.0.0.1", 1, "u", "p", "to@example.com", false)
	changes := []monitor.Change{buildXMPPChange(9000, monitor.ChangeOpened)}
	if err := h.Handle(changes); err == nil {
		t.Error("expected error when server is unreachable")
	}
}
