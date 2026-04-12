package config

import (
	"os"
	"testing"
)

func writeNtfyConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.toml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write config: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_NtfyDefaults(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Ntfy.Enabled {
		t.Error("expected ntfy disabled by default")
	}
	if cfg.Ntfy.ServerURL != "https://ntfy.sh" {
		t.Errorf("expected default server_url https://ntfy.sh, got %s", cfg.Ntfy.ServerURL)
	}
	if cfg.Ntfy.Priority != 3 {
		t.Errorf("expected default priority 3, got %d", cfg.Ntfy.Priority)
	}
}

func TestLoad_NtfySection(t *testing.T) {
	path := writeNtfyConfig(t, `
[ntfy]
enabled = true
topic = "portwatch-alerts"
token = "tk_secret"
priority = 4
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Ntfy.Enabled {
		t.Error("expected ntfy enabled")
	}
	if cfg.Ntfy.Topic != "portwatch-alerts" {
		t.Errorf("expected topic portwatch-alerts, got %s", cfg.Ntfy.Topic)
	}
	if cfg.Ntfy.Token != "tk_secret" {
		t.Errorf("expected token tk_secret, got %s", cfg.Ntfy.Token)
	}
	if cfg.Ntfy.Priority != 4 {
		t.Errorf("expected priority 4, got %d", cfg.Ntfy.Priority)
	}
}

func TestLoad_NtfyMissingTopic(t *testing.T) {
	path := writeNtfyConfig(t, `
[ntfy]
enabled = true
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing topic")
	}
}

func TestLoad_NtfyInvalidPriority(t *testing.T) {
	path := writeNtfyConfig(t, `
[ntfy]
enabled = true
topic = "alerts"
priority = 9
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for invalid priority")
	}
}
