package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeExecConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "portwatch.toml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeExecConfig: %v", err)
	}
	return p
}

func TestLoad_ExecDefaults(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Exec.Enabled {
		t.Error("expected exec disabled by default")
	}
	if cfg.Exec.Timeout != 5*time.Second {
		t.Errorf("expected 5s timeout, got %s", cfg.Exec.Timeout)
	}
	if cfg.Exec.Shell != "" {
		t.Errorf("expected empty shell, got %q", cfg.Exec.Shell)
	}
}

func TestLoad_ExecSection(t *testing.T) {
	p := writeExecConfig(t, `
[exec]
enabled = true
path = "/usr/local/bin/notify.sh"
args = ["--verbose"]
timeout = "10s"
shell = "/bin/bash"
`)
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !cfg.Exec.Enabled {
		t.Error("expected exec enabled")
	}
	if cfg.Exec.Path != "/usr/local/bin/notify.sh" {
		t.Errorf("unexpected path: %s", cfg.Exec.Path)
	}
	if len(cfg.Exec.Args) != 1 || cfg.Exec.Args[0] != "--verbose" {
		t.Errorf("unexpected args: %v", cfg.Exec.Args)
	}
	if cfg.Exec.Timeout != 10*time.Second {
		t.Errorf("expected 10s, got %s", cfg.Exec.Timeout)
	}
	if cfg.Exec.Shell != "/bin/bash" {
		t.Errorf("unexpected shell: %s", cfg.Exec.Shell)
	}
}

func TestLoad_ExecMissingPath(t *testing.T) {
	p := writeExecConfig(t, `
[exec]
enabled = true
`)
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if err := Validate(*cfg); err == nil {
		t.Error("expected validation error for missing path")
	}
}

func TestLoad_ExecInvalidTimeout(t *testing.T) {
	p := writeExecConfig(t, `
[exec]
enabled = true
path = "/bin/alert.sh"
timeout = "120s"
`)
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if err := Validate(*cfg); err == nil {
		t.Error("expected validation error for timeout > 60s")
	}
}
