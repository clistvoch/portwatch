package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeEmailConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "portwatch.toml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return p
}

func TestLoad_EmailDefaults(t *testing.T) {
	p := writeEmailConfig(t, `
port_range = "1-1024"
interval_seconds = 30
state_file = "/tmp/pw.state"
`)
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Email.Enabled {
		t.Error("expected email disabled by default")
	}
	if cfg.Email.Port != 0 {
		t.Errorf("expected zero port, got %d", cfg.Email.Port)
	}
}

func TestLoad_EmailSection(t *testing.T) {
	p := writeEmailConfig(t, `
port_range = "1-1024"
interval_seconds = 30
state_file = "/tmp/pw.state"

[email]
enabled = true
host = "smtp.example.com"
port = 587
username = "user"
password = "secret"
from = "alerts@example.com"
to = ["ops@example.com", "dev@example.com"]
`)
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Email.Enabled {
		t.Error("expected email enabled")
	}
	if cfg.Email.Host != "smtp.example.com" {
		t.Errorf("unexpected host: %s", cfg.Email.Host)
	}
	if len(cfg.Email.To) != 2 {
		t.Errorf("expected 2 recipients, got %d", len(cfg.Email.To))
	}
}

func TestLoad_EmailMissingHost(t *testing.T) {
	p := writeEmailConfig(t, `
port_range = "1-1024"
interval_seconds = 30
state_file = "/tmp/pw.state"

[email]
enabled = true
port = 587
`)
	// No host — still loads; validation of email fields is caller's responsibility.
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Email.Host != "" {
		t.Errorf("expected empty host, got %q", cfg.Email.Host)
	}
}
