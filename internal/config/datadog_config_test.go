package config

import (
	"os"
	"testing"
)

func writeDatadogConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.toml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write config: %v", err)
	}
	_ = f.Close()
	return f.Name()
}

func TestLoad_DatadogDefaults(t *testing.T) {
	path := writeDatadogConfig(t, "[monitor]\ninterval = '10s'\n")
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	dd := cfg.Datadog
	if dd.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if dd.Site != "datadoghq.com" {
		t.Errorf("expected site=datadoghq.com, got %q", dd.Site)
	}
	if dd.Service != "portwatch" {
		t.Errorf("expected service=portwatch, got %q", dd.Service)
	}
}

func TestLoad_DatadogSection(t *testing.T) {
	path := writeDatadogConfig(t, `
[monitor]
interval = '10s'

[datadog]
enabled = true
api_key = "abc123"
site    = "datadoghq.eu"
service = "myapp"
tags    = ["env:prod", "team:ops"]
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	dd := cfg.Datadog
	if !dd.Enabled {
		t.Error("expected Enabled=true")
	}
	if dd.APIKey != "abc123" {
		t.Errorf("expected api_key=abc123, got %q", dd.APIKey)
	}
	if dd.Site != "datadoghq.eu" {
		t.Errorf("expected site=datadoghq.eu, got %q", dd.Site)
	}
	if len(dd.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(dd.Tags))
	}
}

func TestLoad_DatadogMissingAPIKey(t *testing.T) {
	path := writeDatadogConfig(t, `
[monitor]
interval = '10s'

[datadog]
enabled = true
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing api_key")
	}
}
