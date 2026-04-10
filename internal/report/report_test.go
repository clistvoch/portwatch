package report_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/report"
)

func makeSnapshot(ports []monitor.PortInfo, changes []monitor.Change) report.Snapshot {
	return report.Snapshot{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		OpenPorts: ports,
		Changes:   changes,
	}
}

func TestPrintSnapshot_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	p := report.NewPrinter(&buf)

	ports := []monitor.PortInfo{
		{Port: 8080, Proto: "tcp", Address: "127.0.0.1"},
	}
	s := makeSnapshot(ports, nil)
	if err := p.PrintSnapshot(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "8080") {
		t.Errorf("expected port 8080 in output, got:\n%s", out)
	}
	if strings.Contains(out, "Changes detected") {
		t.Errorf("did not expect changes section when no changes")
	}
}

func TestPrintSnapshot_WithChanges(t *testing.T) {
	var buf bytes.Buffer
	p := report.NewPrinter(&buf)

	changes := []monitor.Change{
		{Type: monitor.Opened, Port: monitor.PortInfo{Port: 9090, Proto: "tcp", Address: "0.0.0.0"}},
	}
	s := makeSnapshot(nil, changes)
	p.PrintSnapshot(s)

	out := buf.String()
	if !strings.Contains(out, "Changes detected") {
		t.Errorf("expected changes section, got:\n%s", out)
	}
	if !strings.Contains(out, "9090") {
		t.Errorf("expected port 9090 in changes output")
	}
}

func TestPrintSnapshot_EmptyPorts(t *testing.T) {
	var buf bytes.Buffer
	p := report.NewPrinter(&buf)
	p.PrintSnapshot(makeSnapshot(nil, nil))

	if !strings.Contains(buf.String(), "(none)") {
		t.Errorf("expected '(none)' for empty port list")
	}
}

func TestPrintChangesOnly_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	p := report.NewPrinter(&buf)
	p.PrintChangesOnly(nil)

	if !strings.Contains(buf.String(), "No changes") {
		t.Errorf("expected 'No changes' message")
	}
}

func TestNewPrinter_DefaultsToStdout(t *testing.T) {
	// Should not panic when w is nil
	p := report.NewPrinter(nil)
	if p == nil {
		t.Fatal("expected non-nil printer")
	}
}
