package config

import (
	"os"
	"testing"
)

func writeMattermostConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.toml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	_ = f.Close()
	return f.Name()
}

func TestLoad_MattermostDefaults(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Mattermost.Enabled {
		t.Error("expected mattermost disabled by default")
	}
	if cfg.Mattermost.Username != "portwatch" {
		t.Errorf("expected username 'portwatch', got %q", cfg.Mattermost.Username)
	}
	if cfg.Mattermost.IconEmoji != ":shield:" {
		t.Errorf("expected icon_emoji ':shield:', got %q", cfg.Mattermost.IconEmoji)
	}
}

func TestLoad_MattermostSection(t *testing.T) {
	path := writeMattermostConfig(t, `
[mattermost]
enabled = true
webhook_url = "https://mattermost.example.com/hooks/abc123"
channel = "#alerts"
username = "bot"
icon_emoji = ":bell:"
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Mattermost.Enabled {
		t.Error("expected mattermost enabled")
	}
	if cfg.Mattermost.WebhookURL != "https://mattermost.example.com/hooks/abc123" {
		t.Errorf("unexpected webhook_url: %q", cfg.Mattermost.WebhookURL)
	}
	if cfg.Mattermost.Channel != "#alerts" {
		t.Errorf("unexpected channel: %q", cfg.Mattermost.Channel)
	}
	if cfg.Mattermost.IconEmoji != ":bell:" {
		t.Errorf("unexpected icon_emoji: %q", cfg.Mattermost.IconEmoji)
	}
}

func TestLoad_MattermostMissingWebhook(t *testing.T) {
	path := writeMattermostConfig(t, `
[mattermost]
enabled = true
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing webhook_url")
	}
}
