package alert

import (
	"log/syslog"
	"testing"

	"portwatch/internal/monitor"
	"portwatch/internal/scanner"
)

func buildSyslogChange(kind monitor.ChangeKind, port uint16) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Port: scanner.PortInfo{Port: port, Proto: "tcp"},
	}
}

func TestSyslogHandler_NoChanges(t *testing.T) {
	h, err := NewSyslogHandler(syslog.LOG_INFO|syslog.LOG_DAEMON, "portwatch-test")
	if err != nil {
		t.Skipf("syslog unavailable: %v", err)
	}
	defer h.Close()

	if err := h.Handle(nil); err != nil {
		t.Fatalf("expected no error on empty changes, got %v", err)
	}
}

func TestSyslogHandler_SendsChanges(t *testing.T) {
	h, err := NewSyslogHandler(syslog.LOG_WARNING|syslog.LOG_DAEMON, "portwatch-test")
	if err != nil {
		t.Skipf("syslog unavailable: %v", err)
	}
	defer h.Close()

	changes := []monitor.Change{
		buildSyslogChange(monitor.Opened, 8080),
		buildSyslogChange(monitor.Closed, 9090),
	}

	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error writing to syslog: %v", err)
	}
}

func TestSyslogHandler_DefaultTag(t *testing.T) {
	h, err := NewSyslogHandler(syslog.LOG_INFO|syslog.LOG_USER, "")
	if err != nil {
		t.Skipf("syslog unavailable: %v", err)
	}
	defer h.Close()

	if h.tag != "portwatch" {
		t.Errorf("expected default tag 'portwatch', got %q", h.tag)
	}
}

func TestSyslogHandler_Close(t *testing.T) {
	h, err := NewSyslogHandler(syslog.LOG_INFO|syslog.LOG_USER, "portwatch-test")
	if err != nil {
		t.Skipf("syslog unavailable: %v", err)
	}

	if err := h.Close(); err != nil {
		t.Fatalf("unexpected error on Close: %v", err)
	}
}
