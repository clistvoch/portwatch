package config_test

import (
	"os"
	"testing"

	"github.com/user/portwatch/internal/config"
)

func writeGrafanaConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "grafana-config-*.toml")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_GrafanaDefaults(t *testing.T) {
	path := writeGrafanaConfig(t, "[scan]\nports = \"80-81\"\n")
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Grafana.Enabled {
		t.Error("expected grafana disabled by default")
	}
	if cfg.Grafana.Timeout != 5 {
		t.Errorf("expected default timeout 5, got %d", cfg.Grafana.Timeout)
	}
}

func TestLoad_GrafanaSection(t *testing.T) {
	path := writeGrafanaConfig(t, `
[scan]
ports = "80-81"

[grafana]
enabled = true
url = "http://grafana.local:3000"
api_key = "glsa_abc123"
dashboard_id = "portwatch"
timeout_seconds = 10
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Grafana.Enabled {
		t.Error("expected grafana enabled")
	}
	if cfg.Grafana.URL != "http://grafana.local:3000" {
		t.Errorf("unexpected url: %s", cfg.Grafana.URL)
	}
	if cfg.Grafana.APIKey != "glsa_abc123" {
		t.Errorf("unexpected api_key: %s", cfg.Grafana.APIKey)
	}
	if cfg.Grafana.Timeout != 10 {
		t.Errorf("expected timeout 10, got %d", cfg.Grafana.Timeout)
	}
}

func TestLoad_GrafanaMissingURL(t *testing.T) {
	path := writeGrafanaConfig(t, `
[scan]
ports = "80-81"

[grafana]
enabled = true
api_key = "glsa_abc123"
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for missing grafana url")
	}
}
