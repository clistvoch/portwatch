package alert

import (
	"bufio"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/wlynxg/portwatch/internal/monitor"
)

func buildIRCChange(port int, kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Port: monitor.PortInfo{Port: port, Proto: "tcp"},
		Kind: kind,
	}
}

func TestIRCHandler_NoChanges(t *testing.T) {
	h := NewIRCHandler("localhost", 6667, "bot", "#test", "", false)
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIRCHandler_SendsPayload(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	lines := make(chan string, 20)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
	}()

	addr := ln.Addr().(*net.TCPAddr)
	h := &IRCHandler{
		server:  "127.0.0.1",
		port:    addr.Port,
		nick:    "bot",
		channel: "#test",
		dial:    net.Dial,
	}

	changes := []monitor.Change{
		buildIRCChange(8080, monitor.Opened),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}

	timeout := time.After(2 * time.Second)
	var received []string
loop:
	for {
		select {
		case l := <-lines:
			received = append(received, l)
			if strings.HasPrefix(l, "QUIT") {
				break loop
			}
		case <-timeout:
			t.Fatal("timed out waiting for IRC messages")
		}
	}

	var found bool
	for _, l := range received {
		if strings.Contains(l, "PRIVMSG") && strings.Contains(l, "portwatch") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected PRIVMSG with portwatch, got: %v", received)
	}
}

func TestIRCHandler_InvalidTarget(t *testing.T) {
	h := NewIRCHandler("127.0.0.1", 1, "bot", "#test", "", false)
	changes := []monitor.Change{buildIRCChange(22, monitor.Opened)}
	if err := h.Handle(changes); err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
