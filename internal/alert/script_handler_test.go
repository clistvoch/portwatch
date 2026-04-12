package alert

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/wander/portwatch/internal/monitor"
	"github.com/wander/portwatch/internal/scanner"
)

func buildScriptChange(port int, kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Port: scanner.PortInfo{Port: port, Proto: "tcp"},
	}
}

func TestScriptHandler_NoChanges(t *testing.T) {
	h := NewScriptHandler("/nonexistent", 5, nil)
	if err := h.Handle(nil); err != nil {
		t.Fatalf("expected nil error for empty changes, got %v", err)
	}
}

func TestScriptHandler_RunsScript(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell script test on Windows")
	}

	dir := t.TempDir()
	outFile := filepath.Join(dir, "out.txt")
	script := filepath.Join(dir, "handler.sh")

	content := "#!/bin/sh\ncat > " + outFile + "\n"
	if err := os.WriteFile(script, []byte(content), 0o755); err != nil {
		t.Fatal(err)
	}

	h := NewScriptHandler(script, 5, nil)
	changes := []monitor.Change{
		buildScriptChange(8080, monitor.Opened),
		buildScriptChange(9090, monitor.Closed),
	}

	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("script did not write output: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected non-empty JSON output from script")
	}
}

func TestScriptHandler_InvalidPath(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping on Windows")
	}
	h := NewScriptHandler("/no/such/script.sh", 5, nil)
	changes := []monitor.Change{buildScriptChange(80, monitor.Opened)}
	if err := h.Handle(changes); err == nil {
		t.Fatal("expected error for missing script, got nil")
	}
}

func TestScriptHandler_Timeout(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping on Windows")
	}

	dir := t.TempDir()
	script := filepath.Join(dir, "slow.sh")
	if err := os.WriteFile(script, []byte("#!/bin/sh\nsleep 10\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	h := NewScriptHandler(script, 1, nil)
	changes := []monitor.Change{buildScriptChange(443, monitor.Opened)}
	if err := h.Handle(changes); err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}
