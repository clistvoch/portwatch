package config

import (
	"os"
	"testing"
)

func writeWebhookConfig(t *testing.T, content string) string {
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

func TestLoad_WebhookDefaults(t *testing.T) {
	path := writeWebhookConfig(t, "[scan]\nstart = 1\nend = 1024\n")
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Webhook.Enabled {
		t.Error("expected webhook disabled by default")
	}
	if cfg.Webhook.Timeout != 10 {
		t.Errorf("expected default timeout 10, got %d", cfg.Webhook.Timeout)
	}
}

func TestLoad_WebhookSection(t *testing.T) {
	path := writeWebhookConfig(t, `
[scan]
start = 1
end = 1024

[webhook]
enabled = true
url = "https://example.com/hook"
secret = "s3cr3t"
timeout_seconds = 5
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Webhook.Enabled {
		t.Error("expected webhook enabled")
	}
	if cfg.Webhook.URL != "https://example.com/hook" {
		t.Errorf("unexpected URL: %s", cfg.Webhook.URL)
	}
	if cfg.Webhook.Secret != "s3cr3t" {
		t.Errorf("unexpected secret: %s", cfg.Webhook.Secret)
	}
	if cfg.Webhook.Timeout != 5 {
		t.Errorf("expected timeout 5, got %d", cfg.Webhook.Timeout)
	}
}

func TestLoad_WebhookMissingURL(t *testing.T) {
	path := writeWebhookConfig(t, `
[scan]
start = 1
end = 1024

[webhook]
enabled = true
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing webhook URL")
	}
}
