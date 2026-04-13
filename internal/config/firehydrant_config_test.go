package config_test

import (
	"os"
	"testing"

	"github.com/example/portwatch/internal/config"
)

func writeFirehydrantConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.toml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	_ = f.Close()
	return f.Name()
}

func TestLoad_FirehydrantDefaults(t *testing.T) {
	path := writeFirehydrantConfig(t, "[general]\ninterval_seconds = 30\n")
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Firehydrant.Enabled {
		t.Error("expected firehydrant disabled by default")
	}
	if cfg.Firehydrant.BaseURL != "https://api.firehydrant.io/v1" {
		t.Errorf("unexpected base_url: %s", cfg.Firehydrant.BaseURL)
	}
	if cfg.Firehydrant.Timeout != 10 {
		t.Errorf("unexpected timeout: %d", cfg.Firehydrant.Timeout)
	}
}

func TestLoad_FirehydrantSection(t *testing.T) {
	path := writeFirehydrantConfig(t, `
[general]
interval_seconds = 30

[firehydrant]
enabled = true
api_key = "secret-key"
service_id = "svc-123"
base_url = "https://api.firehydrant.io/v1"
timeout_seconds = 15
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Firehydrant.Enabled {
		t.Error("expected firehydrant enabled")
	}
	if cfg.Firehydrant.APIKey != "secret-key" {
		t.Errorf("unexpected api_key: %s", cfg.Firehydrant.APIKey)
	}
	if cfg.Firehydrant.ServiceID != "svc-123" {
		t.Errorf("unexpected service_id: %s", cfg.Firehydrant.ServiceID)
	}
}

func TestLoad_FirehydrantMissingAPIKey(t *testing.T) {
	path := writeFirehydrantConfig(t, `
[general]
interval_seconds = 30

[firehydrant]
enabled = true
service_id = "svc-123"
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing api_key")
	}
}

func TestLoad_FirehydrantMissingServiceID(t *testing.T) {
	path := writeFirehydrantConfig(t, `
[general]
interval_seconds = 30

[firehydrant]
enabled = true
api_key = "secret-key"
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing service_id")
	}
}
