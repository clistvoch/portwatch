package config

import (
	"os"
	"testing"
)

func writeSyslogConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.toml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_SyslogDefaults(t *testing.T) {
	path := writeSyslogConfig(t, "[scan]\nstart = 1\nend = 1024\n")
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Syslog.Enabled {
		t.Error("expected syslog disabled by default")
	}
	if cfg.Syslog.Tag != "portwatch" {
		t.Errorf("expected default tag 'portwatch', got %q", cfg.Syslog.Tag)
	}
	if cfg.Syslog.Priority != "info" {
		t.Errorf("expected default priority 'info', got %q", cfg.Syslog.Priority)
	}
}

func TestLoad_SyslogSection(t *testing.T) {
	path := writeSyslogConfig(t, `
[scan]
start = 1
end = 1024

[syslog]
enabled = true
network = "udp"
address = "localhost:514"
tag = "myapp"
priority = "warning"
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Syslog.Enabled {
		t.Error("expected syslog enabled")
	}
	if cfg.Syslog.Network != "udp" {
		t.Errorf("expected network 'udp', got %q", cfg.Syslog.Network)
	}
	if cfg.Syslog.Address != "localhost:514" {
		t.Errorf("expected address 'localhost:514', got %q", cfg.Syslog.Address)
	}
	if cfg.Syslog.Tag != "myapp" {
		t.Errorf("expected tag 'myapp', got %q", cfg.Syslog.Tag)
	}
}

func TestLoad_SyslogInvalidPriority(t *testing.T) {
	path := writeSyslogConfig(t, `
[scan]
start = 1
end = 1024

[syslog]
enabled = true
tag = "portwatch"
priority = "verbose"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid priority, got nil")
	}
}

func TestLoad_SyslogMissingAddress(t *testing.T) {
	path := writeSyslogConfig(t, `
[scan]
start = 1
end = 1024

[syslog]
enabled = true
network = "tcp"
tag = "portwatch"
priority = "info"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error when network set but address missing, got nil")
	}
}
