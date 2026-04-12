package config

import (
	"os"
	"testing"
)

func writeSquadcastConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "squadcast_*.toml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_SquadcastDefaults(t *testing.T) {
	cfg := defaultSquadcastConfig()
	if cfg.Enabled {
		t.Error("expected enabled=false by default")
	}
	if cfg.Environment != "production" {
		t.Errorf("expected environment=production, got %s", cfg.Environment)
	}
	if cfg.Timeout != 10 {
		t.Errorf("expected timeout=10, got %d", cfg.Timeout)
	}
}

func TestLoad_SquadcastSection(t *testing.T) {
	path := writeSquadcastConfig(t, `
[squadcast]
enabled = true
webhook_url = "https://api.squadcast.com/v2/incidents/api/abc123"
environment = "staging"
timeout_seconds = 15
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.Squadcast.Enabled {
		t.Error("expected enabled=true")
	}
	if cfg.Squadcast.WebhookURL != "https://api.squadcast.com/v2/incidents/api/abc123" {
		t.Errorf("unexpected webhook_url: %s", cfg.Squadcast.WebhookURL)
	}
	if cfg.Squadcast.Environment != "staging" {
		t.Errorf("expected environment=staging, got %s", cfg.Squadcast.Environment)
	}
}

func TestLoad_SquadcastMissingWebhook(t *testing.T) {
	path := writeSquadcastConfig(t, `
[squadcast]
enabled = true
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := validateSquadcast(cfg.Squadcast); err == nil {
		t.Error("expected error for missing webhook_url")
	}
}
