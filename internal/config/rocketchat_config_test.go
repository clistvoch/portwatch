package config_test

import (
	"os"
	"testing"

	"github.com/user/portwatch/internal/config"
)

func writeRocketChatConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.toml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_RocketChatDefaults(t *testing.T) {
	path := writeRocketChatConfig(t, "[general]\nport_range = \"1-1024\"\n")
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.RocketChat.Enabled {
		t.Error("expected rocketchat to be disabled by default")
	}
	if cfg.RocketChat.Channel != "#general" {
		t.Errorf("expected default channel #general, got %q", cfg.RocketChat.Channel)
	}
	if cfg.RocketChat.Username != "portwatch" {
		t.Errorf("expected default username portwatch, got %q", cfg.RocketChat.Username)
	}
}

func TestLoad_RocketChatSection(t *testing.T) {
	path := writeRocketChatConfig(t, `
[general]
port_range = "1-1024"

[rocketchat]
enabled = true
webhook_url = "https://chat.example.com/hooks/abc123"
channel = "#alerts"
username = "bot"
icon_emoji = ":bell:"
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.RocketChat.Enabled {
		t.Error("expected rocketchat to be enabled")
	}
	if cfg.RocketChat.WebhookURL != "https://chat.example.com/hooks/abc123" {
		t.Errorf("unexpected webhook_url: %q", cfg.RocketChat.WebhookURL)
	}
	if cfg.RocketChat.Channel != "#alerts" {
		t.Errorf("unexpected channel: %q", cfg.RocketChat.Channel)
	}
}

func TestLoad_RocketChatMissingWebhook(t *testing.T) {
	path := writeRocketChatConfig(t, `
[general]
port_range = "1-1024"

[rocketchat]
enabled = true
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for missing webhook_url")
	}
}
