package config_test

import (
	"os"
	"testing"

	"github.com/natemollica-nm/portwatch/internal/config"
)

func writeGoogleChatConfig(t *testing.T, body string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.toml")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(body)
	_ = f.Close()
	return f.Name()
}

func TestLoad_GoogleChatDefaults(t *testing.T) {
	path := writeGoogleChatConfig(t, "")
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.GoogleChat.Enabled {
		t.Error("expected enabled=false by default")
	}
	if cfg.GoogleChat.WebhookURL != "" {
		t.Errorf("expected empty webhook_url, got %q", cfg.GoogleChat.WebhookURL)
	}
}

func TestLoad_GoogleChatSection(t *testing.T) {
	path := writeGoogleChatConfig(t, `
[googlechat]
enabled = true
webhook_url = "https://chat.googleapis.com/v1/spaces/ABC/messages?key=xyz"
thread_key = "portwatch"
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.GoogleChat.Enabled {
		t.Error("expected enabled=true")
	}
	if cfg.GoogleChat.WebhookURL == "" {
		t.Error("expected non-empty webhook_url")
	}
	if cfg.GoogleChat.ThreadKey != "portwatch" {
		t.Errorf("expected thread_key=portwatch, got %q", cfg.GoogleChat.ThreadKey)
	}
}

func TestLoad_GoogleChatMissingWebhook(t *testing.T) {
	path := writeGoogleChatConfig(t, `
[googlechat]
enabled = true
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for missing webhook_url")
	}
}
