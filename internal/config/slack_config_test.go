package config

import (
	"testing"
)

func writeSlackConfig(t *testing.T, content string) string {
	t.Helper()
	return writeTempConfig(t, content)
}

func TestLoad_SlackDefaults(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Slack.WebhookURL != "" {
		t.Errorf("expected empty slack webhook URL by default, got %q", cfg.Slack.WebhookURL)
	}
	if cfg.Slack.Username != "portwatch" {
		t.Errorf("expected default slack username 'portwatch', got %q", cfg.Slack.Username)
	}
	if cfg.Slack.IconEmoji != ":bell:" {
		t.Errorf("expected default slack icon ':bell:', got %q", cfg.Slack.IconEmoji)
	}
}

func TestLoad_SlackSection(t *testing.T) {
	path := writeSlackConfig(t, `
[slack]
webhook_url = "https://hooks.slack.com/services/TEST"
username = "alertbot"
icon_emoji = ":rotating_light:"
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Slack.WebhookURL != "https://hooks.slack.com/services/TEST" {
		t.Errorf("unexpected webhook URL: %q", cfg.Slack.WebhookURL)
	}
	if cfg.Slack.Username != "alertbot" {
		t.Errorf("unexpected username: %q", cfg.Slack.Username)
	}
	if cfg.Slack.IconEmoji != ":rotating_light:" {
		t.Errorf("unexpected icon_emoji: %q", cfg.Slack.IconEmoji)
	}
}

func TestLoad_SlackMissingWebhook(t *testing.T) {
	path := writeSlackConfig(t, `
[slack]
username = "portwatch"
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Slack.WebhookURL != "" {
		t.Errorf("expected empty webhook URL when not set, got %q", cfg.Slack.WebhookURL)
	}
}
