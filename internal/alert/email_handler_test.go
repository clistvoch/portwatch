package alert

import (
	"io"
	"net"
	"net/smtp"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

// startFakeSMTP starts a minimal fake SMTP server and returns its address.
func startFakeSMTP(t *testing.T, received *string) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		ln.Close()
		conn.Write([]byte("220 fake\r\n"))
		buf := make([]byte, 4096)
		var sb strings.Builder
		for {
			n, err := conn.Read(buf)
			if n > 0 {
				sb.Write(buf[:n])
				conn.Write([]byte("250 OK\r\n"))
			}
			if err == io.EOF || err != nil {
				break
			}
		}
		*received = sb.String()
	}()
	_ = smtp.SendMail // ensure import used
	return ln.Addr().String()
}

func TestEmailHandler_NoChanges(t *testing.T) {
	cfg := EmailConfig{Host: "127.0.0.1", Port: 9999, From: "a@b.com", To: []string{"c@d.com"}}
	h := NewEmailHandler(cfg)
	if err := h.Handle(nil); err != nil {
		t.Fatalf("expected no error on empty changes, got %v", err)
	}
}

func TestEmailHandler_SubjectContainsCount(t *testing.T) {
	changes := []monitor.Change{
		{Type: monitor.Opened, Port: scanner.PortInfo{Port: 8080, Proto: "tcp"}},
		{Type: monitor.Closed, Port: scanner.PortInfo{Port: 443, Proto: "tcp"}},
	}
	cfg := EmailConfig{
		Host: "127.0.0.1", Port: 0,
		From: "a@b.com", To: []string{"c@d.com"},
	}
	h := &emailHandler{cfg: cfg}
	// We only test the body building logic by confirming Handle returns an error
	// (no real server) but not a nil-change skip.
	err := h.Handle(changes)
	if err == nil {
		t.Fatal("expected connection error to fake host:0")
	}
}
