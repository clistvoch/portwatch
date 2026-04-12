package config

import (
	"os"
	"testing"
)

func writeMSTeamsConfig(t *testing.T, content string) string {
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

func TestLoad_MSTeamsDefaults(t *testing.T) {
	path := writeMSTeamsConfig(t, "[scan]\nport_range = \"1-1024\"\n")
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	d := defaultMSTeamsConfig()
	if cfg.MSTeams.Enabled != d.Enabled {
		t.Errorf("Enabled: got %v, want %v", cfg.MSTeams.Enabled, d.Enabled)
	}
	if cfg.MSTeams.Title != d.Title {
		t.Errorf("Title: got %q, want %q", cfg.MSTeams.Title, d.Title)
	}
	if cfg.MSTeams.ThemeColor != d.ThemeColor {
		t.Errorf("ThemeColor: got %q, want %q", cfg.MSTeams.ThemeColor, d.ThemeColor)
	}
}

func TestLoad_MSTeamsSection(t *testing.T) {
	path := writeMSTeamsConfig(t, `
[scan]
port_range = "1-1024"

[msteams]
enabled = true
webhook_url = "https://outlook.office.com/webhook/abc"
title = "Port Alert"
theme_color = "0076D7"
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !cfg.MSTeams.Enabled {
		t.Error("expected Enabled=true")
	}
	if cfg.MSTeams.WebhookURL != "https://outlook.office.com/webhook/abc" {
		t.Errorf("WebhookURL: got %q", cfg.MSTeams.WebhookURL)
	}
	if cfg.MSTeams.ThemeColor != "0076D7" {
		t.Errorf("ThemeColor: got %q", cfg.MSTeams.ThemeColor)
	}
}

func TestLoad_MSTeamsMissingWebhook(t *testing.T) {
	path := writeMSTeamsConfig(t, `
[scan]
port_range = "1-1024"

[msteams]
enabled = true
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing webhook_url")
	}
}
