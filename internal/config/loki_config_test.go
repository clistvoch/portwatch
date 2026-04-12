package config_test

import (
	"os"
	"testing"

	"github.com/user/portwatch/internal/config"
)

func writeLokiConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.toml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_LokiDefaults(t *testing.T) {
	cfg := config.DefaultConfig()
	if cfg.Loki.Enabled {
		t.Error("expected loki to be disabled by default")
	}
	if cfg.Loki.Timeout != 5 {
		t.Errorf("expected default timeout 5, got %d", cfg.Loki.Timeout)
	}
	if cfg.Loki.Labels["app"] != "portwatch" {
		t.Errorf("expected default label app=portwatch, got %v", cfg.Loki.Labels)
	}
}

func TestLoad_LokiSection(t *testing.T) {
	path := writeLokiConfig(t, `
[loki]
enabled = true
url = "http://loki:3100"
tenant_id = "tenant1"
timeout_seconds = 10
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Loki.Enabled {
		t.Error("expected loki to be enabled")
	}
	if cfg.Loki.URL != "http://loki:3100" {
		t.Errorf("unexpected url: %s", cfg.Loki.URL)
	}
	if cfg.Loki.TenantID != "tenant1" {
		t.Errorf("unexpected tenant_id: %s", cfg.Loki.TenantID)
	}
	if cfg.Loki.Timeout != 10 {
		t.Errorf("expected timeout 10, got %d", cfg.Loki.Timeout)
	}
}

func TestLoad_LokiMissingURL(t *testing.T) {
	path := writeLokiConfig(t, `
[loki]
enabled = true
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for missing loki url")
	}
}
