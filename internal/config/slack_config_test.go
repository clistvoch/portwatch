package config_test

import (
	"os"
	"testing"

	"github.com/user/portwatch/internal/config"
)

func writeSlackConfig(t *testing.T, content string) string {
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

func TestLoad_SlackDefaults(t *testing.T) {
	path := writeSlackConfig(t, "[core]\nport_range = \"1-1024\"\n")
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Slack.Enabled {
		t.Error("expected slack disabled by default")
	}
	if cfg.Slack.Channel != "#alerts" {
		t.Errorf("expected default channel '#alerts', got %q", cfg.Slack.Channel)
	}
	if cfg.Slack.Username != "portwatch" {
		t.Errorf("expected default username 'portwatch', got %q", cfg.Slack.Username)
	}
}

func TestLoad_SlackSection(t *testing.T) {
	path := writeSlackConfig(t, `[core]
port_range = "1-1024"

[slack]
enabled = true
webhook_url = "https://hooks.slack.com/services/TEST"
channel = "#security"
username = "bot"
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Slack.Enabled {
		t.Error("expected slack enabled")
	}
	if cfg.Slack.WebhookURL != "https://hooks.slack.com/services/TEST" {
		t.Errorf("unexpected webhook_url: %q", cfg.Slack.WebhookURL)
	}
	if cfg.Slack.Channel != "#security" {
		t.Errorf("unexpected channel: %q", cfg.Slack.Channel)
	}
}

func TestLoad_SlackMissingWebhook(t *testing.T) {
	path := writeSlackConfig(t, `[core]
port_range = "1-1024"

[slack]
enabled = true
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing webhook_url")
	}
}
