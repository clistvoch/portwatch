package config_test

import (
	"os"
	"testing"

	"github.com/user/portwatch/internal/config"
)

func writeZendutyConfig(t *testing.T, content string) string {
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

func TestLoad_ZendutyDefaults(t *testing.T) {
	path := writeZendutyConfig(t, "[scanner]\nport_range = \"1-1024\"\n")
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Zenduty.Enabled {
		t.Error("expected zenduty disabled by default")
	}
	if cfg.Zenduty.AlertType != "critical" {
		t.Errorf("expected default alert_type=critical, got %q", cfg.Zenduty.AlertType)
	}
}

func TestLoad_ZendutySection(t *testing.T) {
	path := writeZendutyConfig(t, `
[scanner]
port_range = "1-1024"

[zenduty]
enabled = true
api_key = "testkey"
service_id = "svc-123"
integration_id = "int-456"
alert_type = "warning"
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Zenduty.Enabled {
		t.Error("expected zenduty enabled")
	}
	if cfg.Zenduty.APIKey != "testkey" {
		t.Errorf("expected api_key=testkey, got %q", cfg.Zenduty.APIKey)
	}
	if cfg.Zenduty.AlertType != "warning" {
		t.Errorf("expected alert_type=warning, got %q", cfg.Zenduty.AlertType)
	}
}

func TestLoad_ZendutyMissingAPIKey(t *testing.T) {
	path := writeZendutyConfig(t, `
[scanner]
port_range = "1-1024"

[zenduty]
enabled = true
service_id = "svc-123"
integration_id = "int-456"
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing api_key")
	}
}

func TestLoad_ZendutyInvalidAlertType(t *testing.T) {
	path := writeZendutyConfig(t, `
[scanner]
port_range = "1-1024"

[zenduty]
enabled = true
api_key = "key"
service_id = "svc"
integration_id = "int"
alert_type = "urgent"
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for invalid alert_type")
	}
}
