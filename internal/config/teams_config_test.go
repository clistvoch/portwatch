package config

import (
	"os"
	"testing"
)

func writeTeamsConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.toml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_TeamsDefaults(t *testing.T) {
	path := writeTeamsConfig(t, "[scan]\nport_range = \"1-1024\"\n")
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Teams.Enabled {
		t.Error("expected Teams.Enabled to be false by default")
	}
	if cfg.Teams.Title != "PortWatch Alert" {
		t.Errorf("expected default title, got %q", cfg.Teams.Title)
	}
}

func TestLoad_TeamsSection(t *testing.T) {
	content := `
[scan]
port_range = "1-1024"

[teams]
enabled = true
webhook_url = "https://outlook.office.com/webhook/abc123"
title = "My Alert"
`
	path := writeTeamsConfig(t, content)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Teams.Enabled {
		t.Error("expected Teams.Enabled to be true")
	}
	if cfg.Teams.WebhookURL != "https://outlook.office.com/webhook/abc123" {
		t.Errorf("unexpected webhook_url: %q", cfg.Teams.WebhookURL)
	}
	if cfg.Teams.Title != "My Alert" {
		t.Errorf("unexpected title: %q", cfg.Teams.Title)
	}
}

func TestLoad_TeamsMissingWebhook(t *testing.T) {
	content := `
[scan]
port_range = "1-1024"

[teams]
enabled = true
`
	path := writeTeamsConfig(t, content)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing webhook_url, got nil")
	}
}
